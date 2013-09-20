/*
 * newspeak.io
 * Copyright (C) 2013 Jahn Bertsch
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License version 3
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package filter

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/fitstar/falcore"
	"io/ioutil"
	"net/http"
	"newspeak/response_messages"
	"newspeak/services"
	"regexp"
	"strconv"
	"time"
)

// TokenFilter "class"
type TokenFilter int

// access tokens
func (f TokenFilter) FilterRequest(request *falcore.Request) *http.Response {
	var response *http.Response

	if request.HttpRequest.Method == "POST" {
		request.CurrentStage.Status = byte(1)
		response = validateCredentials(request, request.HttpRequest.FormValue("username"), request.HttpRequest.FormValue("password"))
	}
	return response
}

// check if password and username only contain allowed characters
func validateCredentials(request *falcore.Request, username string, password string) *http.Response {
	var response *http.Response

	if username == "" {
		request.CurrentStage.Status = byte(2)
		response = response_messages.ErrorResponse(request, http.StatusBadRequest, "bad request", "username is missing", nil)
	} else if password == "" {
		request.CurrentStage.Status = byte(3)
		response = response_messages.ErrorResponse(request, http.StatusBadRequest, "bad request", "password is missing", nil)
	} else {
		var alnumValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")
		if !alnumValidator.MatchString(username) {
			request.CurrentStage.Status = byte(4)
			response = response_messages.ErrorResponse(request, http.StatusBadRequest, "bad request", "username is invalid", nil)
		} else {
			if !alnumValidator.MatchString(password) {
				request.CurrentStage.Status = byte(5)
				response = response_messages.ErrorResponse(request, http.StatusBadRequest, "bad request", "password is invalid", nil)
			} else {
				// password and username do not contain invalid characters
				request.CurrentStage.Status = byte(6)
				response = generateAccessToken(request, username, password)
			}
		}
	}
	return response
}

// asks for an usergrid access token with passed username and password
func generateAccessToken(request *falcore.Request, username string, password string) *http.Response {
	var response *http.Response

	uri := fmt.Sprintf("http://localhost:8080/newspeak/newspeak/token?grant_type=password&username=%s&password=%s", username, password)
	usergridResponse, usergridError := http.Get(uri)

	if usergridError != nil {
		request.CurrentStage.Status = byte(7)
		response = response_messages.InternalServerErrorResponse(request, "Could not request credentials", usergridError)
	} else if usergridResponse.StatusCode != http.StatusOK {
		request.CurrentStage.Status = byte(8)
		response = response_messages.ErrorResponse(request, http.StatusUnauthorized, "unauthorized", "Invalid credentials", nil)
		fmt.Println("invalid credentials: " + username + ", " + password)
	} else {
		usergridResponseBody, usergridError := ioutil.ReadAll(usergridResponse.Body)
		usergridResponse.Body.Close()

		if usergridError != nil {
			request.CurrentStage.Status = byte(9)
			response = response_messages.InternalServerErrorResponse(request, "Could not read response when requesting credentials", usergridError)
		} else {
			request.CurrentStage.Status = byte(10)
			response = parseUsergridResponse(request, username, usergridResponseBody)
		}
	}
	return response
}

// extract access token out of usergrid response and pass it on
func parseUsergridResponse(request *falcore.Request, username string, usergridResponseBody []byte) *http.Response {
	var response *http.Response

	var parsedJson interface{}
	unmarshalError := json.Unmarshal(usergridResponseBody, &parsedJson)
	if unmarshalError != nil {
		request.CurrentStage.Status = byte(11)
		response = response_messages.InternalServerErrorResponse(request, "Could not process response when requesting credentials", unmarshalError)
	} else {
		// extract access token and store in memcached
		var accessToken string
		parsedJsonMap := parsedJson.(map[string]interface{})
		for key, value := range parsedJsonMap {
			switch realValue := value.(type) {
			case string:
				if key == "access_token" {
					accessToken = realValue
					break
				}
			}
		}

		// store token request time and username in memcached
		time := time.Now().UnixNano()
		value := fmt.Sprintf("%s %s", strconv.FormatInt(time, 10), username)
		services.Memcached.Set(&memcache.Item{Key: accessToken, Value: []byte(value)})

		request.CurrentStage.Status = byte(12)
		response = accessTokenResponse(request, accessToken, username)
	}
	return response
}

// create json response containing the access token
func accessTokenResponse(request *falcore.Request, accessToken string, username string) *http.Response {
	type Json struct {
		AccessToken string
	}
	json := Json{accessToken}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " created acccess token '" + accessToken + "' for '" + username + "'")
	return response_messages.SuccessResponse(request, json)
}

// returns the token passed with the request
func getAccessToken(request *falcore.Request) string {
	var accessToken string
	if request.HttpRequest.FormValue("accessToken") != "" {
		accessToken = request.HttpRequest.FormValue("accessToken")
	} else if request.HttpRequest.Header.Get("Authentication") != "" {
		accessToken = request.HttpRequest.Header.Get("Authentication")
	}
	return accessToken
}

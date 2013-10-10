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
	"errors"
	"fmt"
	"github.com/fitstar/falcore"
	"io/ioutil"
	"net/http"
	"net/url"
	"newspeak/response_messages"
	"newspeak/services"
)

// DevicesFilter "class"
type DevicesFilter int

// Device managment api endpoint
func (f DevicesFilter) FilterRequest(request *falcore.Request) *http.Response {
	var response *http.Response

	if request.HttpRequest.Method == "POST" {
		response = addDeviceToUniqush(request, request.HttpRequest.FormValue("deviceToken"))
	}
	return response
}

// register device with apple/google push message server
func addDeviceToUniqush(request *falcore.Request, deviceToken string) *http.Response {
	var response *http.Response
	username := services.MemcacheGetUsername(getAccessToken(request))

	uniqushResponse, uniqushError := http.PostForm("http://localhost:9898/subscribe", url.Values{
		"service":         {"newspeak"},
		"pushservicetype": {"apns"},
		"subscriber":      {username},
		"devtoken":        {deviceToken},
	})

	if uniqushError != nil {
		request.CurrentStage.Status = byte(1)
		response = response_messages.InternalServerErrorResponse(request, "Could not read response when registering device", uniqushError)
	} else {
		uniqushResponseBody, uniqushError := ioutil.ReadAll(uniqushResponse.Body)
		uniqushResponse.Body.Close()

		if uniqushResponse.StatusCode != http.StatusOK {
			request.CurrentStage.Status = byte(2)
			response = response_messages.ErrorResponse(request, http.StatusServiceUnavailable, "service unavailable", "Error while connecting to push message server", errors.New(string(uniqushResponseBody)))
		} else if uniqushError != nil {
			request.CurrentStage.Status = byte(3)
			response = response_messages.InternalServerErrorResponse(request, "Could not process response when registering device", uniqushError)
		} else {
			var body = make(map[string]string)
			body["Message"] = "device added"
			body["Username"] = username
			body["DeviceToken"] = deviceToken
			body["Response"] = string(uniqushResponseBody)
			request.CurrentStage.Status = byte(4)
			fmt.Println("registered device:", deviceToken, "for:", username)
			response = response_messages.SuccessResponse(request, body)
		}
	}
	return response
}

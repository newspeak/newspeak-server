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
	"time"
)

// MessagesFilter "class"
type MessagesFilter int

// Device managment api endpoint
func (f MessagesFilter) FilterRequest(request *falcore.Request) *http.Response {
	var response *http.Response

	if request.HttpRequest.Method == "POST" {
		request.CurrentStage.Status = byte(1)
		// @TODO check if parameters are valid
		response = sendMessageToUniqush(request)
	}
	return response
}

// register device with apple/google push message server
// @TODO check if message is empty
func sendMessageToUniqush(request *falcore.Request) *http.Response {
	var response *http.Response
	recipient := request.HttpRequest.FormValue("recipient")
	message := request.HttpRequest.FormValue("message")

	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " delivering message '" + message + "' to '" + recipient + "'")

	uniqushResponse, uniqushError := http.PostForm("http://localhost:9898/push", url.Values{
		"service":    {"newspeak"},
		"subscriber": {recipient},
		"msg":        {message},
	})

	if uniqushError != nil {
		request.CurrentStage.Status = byte(2)
		response = response_messages.InternalServerErrorResponse(request, "Could not read response when registering device", uniqushError)
	} else {
		uniqushResponseBody, uniqushError := ioutil.ReadAll(uniqushResponse.Body)
		uniqushResponse.Body.Close()

		if uniqushError != nil {
			request.CurrentStage.Status = byte(3)
			response = response_messages.InternalServerErrorResponse(request, "Could not process response when registering device", uniqushError)
		} else if uniqushResponse.StatusCode != http.StatusOK {
			request.CurrentStage.Status = byte(4)
			response = response_messages.ErrorResponse(request, http.StatusServiceUnavailable, "service unavailable", "Error while connecting to push message server", errors.New(string(uniqushResponseBody)))
		} else {
			var body = make(map[string]string)
			body["Message"] = "message sent successfully"
			body["Recipient"] = recipient
			body["MessageSent"] = message
			body["Response"] = string(uniqushResponseBody)
			request.CurrentStage.Status = byte(5)
			fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" delivered message '"+message+"' to '", recipient+"' response: "+string(uniqushResponseBody))
			response = response_messages.SuccessResponse(request, body)
		}
	}
	return response
}

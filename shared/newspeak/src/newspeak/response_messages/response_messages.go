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

package response_messages

import (
	"fmt"
	"github.com/fitstar/falcore"
	"log"
	"net/http"
)

// return response message with encoded json data and http status code 200 (OK)
func SuccessResponse(request *falcore.Request, jsonData interface{}) *http.Response {
	response, jsonError := falcore.JSONResponse(request.HttpRequest, http.StatusOK, nil, jsonData)
	if jsonError != nil {
		response = falcore.StringResponse(request.HttpRequest, http.StatusInternalServerError, nil, fmt.Sprintf("JSON error: %s", jsonError))
	}
	return response
}

// return response message with encoded json error message and a custom error code
// if logMessage is nil, do not write to log
func ErrorResponse(request *falcore.Request, status int, error string, errorDescription string, logMessage error) *http.Response {
	if logMessage != nil {
		log.Println(fmt.Sprintf("%s: %s, %s", error, errorDescription, logMessage))
	}

	type Json struct {
		Error            string
		ErrorDescription string
	}
	json := Json{error, errorDescription}

	response, jsonError := falcore.JSONResponse(request.HttpRequest, status, nil, json)
	if jsonError != nil {
		response = falcore.StringResponse(request.HttpRequest, http.StatusInternalServerError, nil, fmt.Sprintf("JSON error: %s", jsonError))
		log.Println(fmt.Sprintf("%s: %s, %s, json error:", error, errorDescription, logMessage, jsonError))
	}
	return response
}

// return internal server error status code with error message. also log the error.
func InternalServerErrorResponse(request *falcore.Request, errorDescription string, logMessage error) *http.Response {
	return ErrorResponse(request, http.StatusInternalServerError, "internal server error", errorDescription, logMessage)
}

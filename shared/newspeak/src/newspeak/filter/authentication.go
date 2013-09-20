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
	"github.com/fitstar/falcore"
	"net/http"
	"newspeak/response_messages"
	"newspeak/services"
)

// AuthenticationFilter "class"
type AuthenticationFilter int

// authenticate request
func (f AuthenticationFilter) FilterRequest(request *falcore.Request) *http.Response {
	var response *http.Response

	if len(request.HttpRequest.URL.Path) == 7 && request.HttpRequest.URL.Path[0:7] == "/tokens" {
		// current request is a token request, do not do authentication here
		request.CurrentStage.Status = byte(1)
		response = nil
	} else {
		// check if access token is set in header "Authentication" or parameter "accessToken"
		if getAccessToken(request) == "" {
			// do not allow empty access token
			request.CurrentStage.Status = byte(2)
			response = response_messages.ErrorResponse(request, http.StatusUnauthorized, "unauthorized", "Missing access token", nil)
		} else {
			// check if access token is not expired yet
			if services.MemcacheAccessTokenIsExpired(getAccessToken(request)) {
				request.CurrentStage.Status = byte(3)
				response = response_messages.ErrorResponse(request, http.StatusUnauthorized, "unauthorized", "Invalid access token", nil)
			} else {
				// authentication ok!
				request.CurrentStage.Status = byte(4)
				response = nil
			}
		}
	}
	return response
}

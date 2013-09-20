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
)

// NotFoundFilter "class"
type NotFoundFilter int

// catch all filter: if no filter matched, return 404 not found as default response
func (f NotFoundFilter) FilterRequest(request *falcore.Request) *http.Response {
	return falcore.StringResponse(request.HttpRequest, http.StatusNotFound, nil, "Not Found")
}

// SwaggerHeadersFilter "class"
type SwaggerHeadersFilter int

// downstream filter, headers required for swagger-ui to work
func (filter *SwaggerHeadersFilter) FilterResponse(request *falcore.Request, response *http.Response) {
	response.Header.Set("Access-Control-Allow-Origin", "*")
	response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
}

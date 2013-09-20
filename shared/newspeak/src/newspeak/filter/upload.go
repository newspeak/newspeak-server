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
	"fmt"
	"github.com/fitstar/falcore"
	"io"
	"math"
	"net/http"
	"os"
)

// upload limit
const maxUploadBytes = 1024 * 1024 * 10 // 10 mb

// UploadFilter "class"
type UploadFilter int

// Handle upload
func (f UploadFilter) FilterRequest(request *falcore.Request) *http.Response {
	multipartReader, error := request.HttpRequest.MultipartReader()

	if error == http.ErrNotMultipart {
		return falcore.StringResponse(request.HttpRequest, http.StatusBadRequest, nil, "Bad Request // Not Multipart")
	} else if error != nil {
		return falcore.StringResponse(request.HttpRequest, http.StatusInternalServerError, nil, "Upload Error: "+error.Error()+"\n")
	}

	length := request.HttpRequest.ContentLength
	fmt.Println("content length:", length)

	part, error := multipartReader.NextPart()
	if error != io.EOF {
		var read int64
		var percent, previousPercent float64
		var filename string

		// find non existing filename
		for i := 1; ; i++ {
			filename = fmt.Sprintf("uploadedFile%v.mov", i)
			_, err := os.Stat(filename)
			if os.IsNotExist(err) {
				break
			}
		}

		fmt.Println(filename)
		destination, error := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
		if error != nil {
			return falcore.StringResponse(request.HttpRequest, http.StatusInternalServerError, nil, "Internal Server Error // Could not create file")
		}

		buffer := make([]byte, 1024) // 100 kB
		for {
			byteCount, error := part.Read(buffer)
			read += int64(byteCount)
			percent = math.Floor(float64(read)/float64(length)*100) + 1

			if error == io.EOF {
				break
			} else if read > maxUploadBytes {
				fmt.Println("file too big")
				return falcore.StringResponse(request.HttpRequest, http.StatusRequestEntityTooLarge, nil, "Request Entity Too Large")
			}

			if percent != previousPercent {
				previousPercent = percent
				fmt.Printf("progress %v%%, read %fmb, %v byte of %v\n", percent, (float64(read) / (1024 * 1024)), read, length)
			}
			destination.Write(buffer)
		}
	}

	return falcore.StringResponse(request.HttpRequest, http.StatusOK, nil, "upload finished\n")
}

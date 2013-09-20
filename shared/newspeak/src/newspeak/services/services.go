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

package services

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/peterbourgon/g2s"
	"log"
	"strconv"
	"strings"
	"time"
)

// statsd, package variable
var Statsd g2s.Statter

// statsd, package variable
var Memcached *memcache.Client

// token is valid for this many minutes
const tokenLifetime = 30

// get values stored in memcache for passed access token
// returns "" on chache miss
func MemcacheGetUsername(accessToken string) string {
	item, error := Memcached.Get(accessToken)
	if error == nil {
		// cache hit
		items := strings.Split(string(item.Value), " ") // example for item.Value: "1369837012584360116 user"
		return items[1]
	}
	return ""
}

// returns true if passed access token is still valid
func MemcacheAccessTokenIsExpired(accessToken string) bool {
	item, error := Memcached.Get(accessToken)
	if error == nil {
		// cache hit
		itemString := string(item.Value) // example: "1369837012584360116 user"
		items := strings.Split(itemString, " ")
		nanoSeconds, error := strconv.ParseInt(items[0], 10, 64)
		if error != nil {
			log.Println("error while parsing user cache result '" + itemString + "' with error: " + error.Error())
		} else {
			issuedTime := time.Unix(0, nanoSeconds)
			issuedMinutesAgo := int(time.Since(issuedTime).Minutes())

			if issuedMinutesAgo > tokenLifetime {
				return true // expired
			} else {
				return false // not expired
			}
		}
	}
	return true
}

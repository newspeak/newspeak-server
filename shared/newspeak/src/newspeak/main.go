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

package main

import (
	"flag"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/fitstar/falcore"
	"github.com/fitstar/falcore/router"
	"github.com/peterbourgon/g2s"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"newspeak/filter"
	"newspeak/services"
	"time"
)

// command line options
var (
	port = flag.Int("port", 80, "the port to listen on")
)

// entry point
func main() {
	// parse command line options
	flag.Parse()
	PrintLog(fmt.Sprintf("starting newspeak api server at http://localhost:%d\n", *port))

	server := setupServer()

	services.Statsd.Counter(1.0, "api.server-start", 1) // count how often server was started

	// start the server, this is normally blocking forever unless you send lifecycle commands
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Could not start server:", err)
	}
}

// init server and its filter pipeline
func setupServer() *falcore.Server {
	pipeline := falcore.NewPipeline()

	var authentication filter.AuthenticationFilter
	var token filter.TokenFilter
	var devices filter.DevicesFilter
	var messages filter.MessagesFilter
	// var upload filter.UploadFilter
	var notFound filter.NotFoundFilter
	swaggerHeaders := new(filter.SwaggerHeadersFilter)

	// route based on request path
	router := router.NewPathRouter()
	router.AddMatch("^/tokens", token)
	router.AddMatch("^/devices", devices)
	router.AddMatch("^/messages", messages)
	//router.AddMatch("^/upload", upload)

	pipeline.Upstream.PushBack(authentication) // authentication is first in pipeline since it's mandatory
	pipeline.Upstream.PushBack(router)
	pipeline.Upstream.PushBack(notFound) // catch all filter is last, default response

	pipeline.Downstream.PushBack(swaggerHeaders)

	server := falcore.NewServer(*port, pipeline)

	// add request done callback stage
	//server.CompletionCallback = completionCallback

	setupServices()

	return server
}

// setup connections to other services running on this host
func setupServices() {
	var error error
	services.Statsd, error = g2s.Dial("udp", "localhost:8125")
	if error != nil {
		log.Fatal("could not set up statsd client.")
	}
	services.Memcached = memcache.New("localhost:11211")

	// setup push service
	usingSandbox := "true" // true or false
	uniqushResponse, uniqushError := http.PostForm("http://localhost:9898/addpsp", url.Values{
		"pushservicetype": {"apns"},
		"service":         {"newspeak"},
		"cert":            {"/etc/newspeak/apns-certs/cert.pem"},
		"key":             {"/etc/newspeak/apns-certs/priv-noenc.pem"},
		"sandbox":         {usingSandbox},
	})
	if uniqushError != nil {
		log.Fatal("could not add push service provider for apple push notifications: " + string(uniqushError.Error()))
	} else {
		uniqushResponseBodyBytes, uniqushError := ioutil.ReadAll(uniqushResponse.Body)
		uniqushResponseBody := string(uniqushResponseBodyBytes)
		uniqushResponse.Body.Close()
		if uniqushError != nil {
			log.Fatal("could not read response when adding push service provider for apple push notifications: " + string(uniqushError.Error()))
		} else if uniqushResponseBody[0:30] != "[AddPushServiceProvider][Info]" {
			log.Fatal("invalid response when adding push service provider for apple push notifications: " + uniqushResponseBody)
		} else {
			fmt.Println("added push service provider for apple push notifications. usingSandbox:" + usingSandbox + ", uniqush response:" + uniqushResponseBody)
		}
	}
}

// print detailed stats about the request to the log and push stats data to graphite.
// runs as separate goroutine, based on falcore.request.go:Trace()
// falcore docs say this is a big hit on performance and should only be used for debugging or development.
var completionCallback = func(falcoreRequest *falcore.Request, response *http.Response) {
	go func() {
		Statsd, statsdError := g2s.Dial("udp", "localhost:8125")
		requestTimeDiff := falcore.TimeDiff(falcoreRequest.StartTime, falcoreRequest.EndTime)
		httpRequest := falcoreRequest.HttpRequest

		// stats for the whole request
		falcore.Trace("%s [%s] %s%s S=%v Sig=%s Tot=%.4fs", falcoreRequest.ID, httpRequest.Method, httpRequest.Host, httpRequest.URL, response.StatusCode, falcoreRequest.Signature(), requestTimeDiff)
		if statsdError == nil {
			Statsd.Timing(1.0, "api.request-time", falcoreRequest.EndTime.Sub(falcoreRequest.StartTime))
		}

		// stats for each pipeline stage
		stages := falcoreRequest.PipelineStageStats
		for stage := stages.Front(); stage != nil; stage = stage.Next() {
			pipelineStageStats, _ := stage.Value.(*falcore.PipelineStageStat)

			stageTimeDiff := falcore.TimeDiff(pipelineStageStats.StartTime, pipelineStageStats.EndTime)
			falcore.Trace("%s [%s]%-30s S=%2d Tot=%.4fs %%=%.2f", falcoreRequest.ID, pipelineStageStats.Type, pipelineStageStats.Name, pipelineStageStats.Status, stageTimeDiff, stageTimeDiff/(requestTimeDiff*100.0))

			if statsdError == nil {
				stageName := fmt.Sprintf("api.pipeline.%s", pipelineStageStats.Name)
				if pipelineStageStats.Name[0:1] == "*" {
					// remove * character, if necessary
					stageName = fmt.Sprintf("api.pipeline.%s", pipelineStageStats.Name[1:len(pipelineStageStats.Name)])
				}
				Statsd.Timing(1.0, stageName, pipelineStageStats.EndTime.Sub(pipelineStageStats.StartTime))
			}
		}
		falcore.Trace("%s %-30s     S= 0 Tot=%.4fs %%=%.2f", falcoreRequest.ID, "Overhead", float32(falcoreRequest.Overhead)/float32(time.Second), float32(falcoreRequest.Overhead)/float32(time.Second)/requestTimeDiff*100.0)
		fmt.Print("\n")
	}() // note the parentheses - must call the function
}

// prints a string to stdout with current time prepended
func PrintLog(message string) {
	fmt.Printf("%s: %s", time.Now().Format("2006-01-02 15:04:05"), message)
}

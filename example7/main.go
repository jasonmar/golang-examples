// Copyright 2017 Jason Mar. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	// inputLimit sets an upper bound on request JSON bytes to decode.
	inputLimit = 65535
)

var (
	dispatch Dispatcher
	queue    chan *Task
)

func postHandler(w http.ResponseWriter, r *http.Request, h httprouter.Params) {
	t := time.Now() // take timestamp for latency calculation

	// Decode JSON
	body := io.LimitReader(r.Body, inputLimit)
	decoder := json.NewDecoder(body)
	d := InputData{}
	if err := decoder.Decode(&d); err != nil {
		log.Println("Handler: Received invalid JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Iterate through decoded InputData
	for i := range d.Data {
		p := &d.Data[i] // get a pointer to element
		log.Printf("Handler: Received Parcel %s", p.ID)
		queue <- &Task{Time: t, Parcel: p} // create a Task and send its address to the queue channel
	}

	w.WriteHeader(http.StatusOK)
}

func init() {
	// Read queue size from environment variable
	if bufSize, err := strconv.Atoi(os.Getenv("MAX_QUEUE")); err == nil {
		queue = make(chan *Task, bufSize)
	} else {
		log.Fatalln("MAX_QUEUE not set")
	}

	// Read worker count from environment variable
	if nWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS")); err == nil {
		dispatch = NewDispatcher(int(nWorkers))
		go dispatch.Dispatch(queue)
	} else {
		log.Fatalln("MAX_WORKERS not set")
	}
}

func main() {
	route := httprouter.New()
	route.POST("/", postHandler)
	http.ListenAndServe(":8080", route)
}

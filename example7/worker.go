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
	"log"
)

// Worker represents the worker that executes the job.
type Worker struct {
	id   int
	task chan *Task
	quit chan bool
}

// NewWorker creates a Worker.
func NewWorker(id int) *Worker {
	return &Worker{
		id:   id,
		task: make(chan *Task),
		quit: make(chan bool),
	}
}

// Work begins task loop.
func (w *Worker) Work(c chan<- chan<- *Task) {
	for {
		c <- w.task // Send channel
		select {    // Block until Task or quit
		case t := <-w.task:
			log.Printf("Worker %d: Received Task with Parcel %s", w.id, t.Parcel.ID)
			if err := t.Complete(); err != nil {
				log.Printf("Worker %d: Error completing Task (%s)", w.id, err.Error())
			}
		case <-w.quit:
			return
		}
	}
}

// Stop signals worker to stop listening for work requests.
func (w *Worker) Stop() {
	w.quit <- true
}

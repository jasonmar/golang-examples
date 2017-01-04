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

// Dispatcher maintains a buffered channel of channels which receive work.
type Dispatcher struct {
	size int
}

// NewDispatcher creates a Dispatcher with a specified number of workers.
func NewDispatcher(size int) Dispatcher {
	return Dispatcher{size: size}
}

// Dispatch creates workers and begins distributing tasks.
func (d *Dispatcher) Dispatch(queue <-chan *Task) {
	c := make(chan chan<- *Task, d.size)
	for i := 0; i < d.size; i++ {
		w := NewWorker(i)
		go w.Work(c)
	}
	for {
		select {
		case task := <-queue:
			go func(t *Task) {
				w := <-c // Block until a worker is ready to receive Task
				w <- t   // Send Task
			}(task)
		}
	}
}

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
	"math/rand"
	"time"
)

// InputData is used to decode JSON into a slice of Parcels.
type InputData struct {
	Data []Parcel
}

// Parcel contains fields from JSON
type Parcel struct {
	ID   string
	Time time.Time
	Data string
}

// Task contains data for a single task.
type Task struct {
	Time   time.Time
	Parcel *Parcel
}

// Complete performs a prespecified task.
func (t *Task) Complete() (err error) {
	// TODO: implement upload to to S3 or publish to SNS, Kafka, etc
	amt := time.Duration(10 + rand.Intn(50))
	time.Sleep(time.Millisecond * amt)
	latency := time.Since(t.Time)
	log.Printf("Task: Completed for Parcel %s (%s)", t.Parcel.ID, latency)
	return
}

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
	"html/template"
	"log"
	"net/http"
)

// Dog has a name and an age
type Dog struct {
	Name string
	Age  int
}

// Cat has a name and a number of lives
type Cat struct {
	Name  string
	Lives int
}

var data = struct {
	Dogs []Dog
	Cats []Cat
}{
	Dogs: []Dog{
		Dog{"Snoopy", 8},
		Dog{"Dogbert", 3},
	},
	Cats: []Cat{
		Cat{"Garfield", 9},
		Cat{"Selina", 2},
	},
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	err = t.Execute(w, &data)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":8080", nil)
}

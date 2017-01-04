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
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var td, tc *template.Template

func init() {
	var err0, err1 error
	td, err0 = template.ParseFiles("dogs.gohtml")
	if err0 != nil {
		panic(err0)
	}
	tc, err1 = template.ParseFiles("cats.gohtml")
	if err1 != nil {
		panic(err1)
	}
}

// Dog has a name and an age
type Dog struct {
	Name string
	Age  int
}

// Cat has a name and a number of lives
type Cat struct {
	Name  string
	Age   int
	Lives int
}

var dogs = map[string]Dog{
	"Snoopy":  Dog{"Snoopy", 8},
	"Dogbert": Dog{"Dogbert", 3},
}

var cats = map[string]Cat{
	"Garfield": Cat{"Garfield", 6, 9},
	"Selina":   Cat{"Selina", 27, 2},
}

func dogHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if dog, found := dogs[p.ByName("name")]; found {
		err := td.Execute(w, &dog)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		http.NotFound(w, r)
	}
}

func catHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	l1, err := strconv.ParseInt(p.ByName("lives"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
	}

	if c0, found := cats[p.ByName("name")]; found {
		cat := Cat{c0.Name, c0.Age, int(l1)}
		err1 := tc.Execute(w, &cat)
		if err1 != nil {
			log.Fatalln(err)
		}
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	route := httprouter.New()
	route.GET("/dogs/:name", dogHandler)
	route.GET("/cats/:name/:lives", catHandler)
	http.ListenAndServe(":8080", route)
}

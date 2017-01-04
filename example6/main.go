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
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/julienschmidt/httprouter"
)

// LoginForm consists of email, password and csrf token
type LoginForm struct {
	Email    string // Names must be capital to use with json.Decoder!!
	Password string
	CSRF     string
}

var t1, t2 *template.Template

var passwords = map[string]string{
	"Snoopy":  "Charlie",
	"Dogbert": "Dilbert",
}

func init() {
	var err2 error
	t2, err2 = template.ParseFiles("success.gohtml")
	if err2 != nil {
		log.Fatalln(err2)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var f LoginForm
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&f)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("%s", f)

	if f.CSRF != "afe8492c00c784295f83330ce7dccaba9bb188b01566e87fceb6794d0e7d9e9d" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if len(f.Email) == 0 || len(f.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pw, found := passwords[f.Email]
	if !found { // Email is not registered
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if pw != f.Password { // Password is incorrect
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "email", Value: f.Email, Expires: expiration}
	http.SetCookie(w, &cookie)
	t2.Execute(w, &f.Email)
}

func getHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fp := path.Join("static", "login.html")
	http.ServeFile(w, r, fp)
}

func main() {
	route := httprouter.New()
	route.POST("/login", postHandler)
	route.GET("/login", getHandler)
	http.ListenAndServe(":8080", route)
}

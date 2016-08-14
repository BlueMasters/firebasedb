// Copyright 2016 Jacques Supcik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// References:
// https://www.firebase.com/docs/rest/api
// https://firebase.google.com/docs/reference/rest/database
// https://firebase.google.com/docs/database/rest/structure-data
// https://firebase.google.com/docs/database/rest/save-data
// https://firebase.google.com/docs/database/rest/retrieve-data

package firebasedb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	pathLib "path"
)

// WithHttpClient sets a custom HTTP client for the REST requests. If set to nil (default),
// then http.DefaultClient is used.
func (r Reference) WithHttpClient(c *http.Client) Reference {
	result := r
	result.client = c
	return result
}

// addAuth returns a new reference with authentication information (if available).
func (r Reference) addAuth() Reference {
	if r.auth != nil {
		return r.withParam(r.auth.ParamName(), r.auth.String())
	} else {
		return r
	}
}

// jsonUrl is an internal function to build the URL for the REST API
// See https://firebase.google.com/docs/reference/rest/database/ "API Usage".
func (r Reference) jsonUrl() string {
	u := r.url
	u.Path = pathLib.Clean(u.Path)
	if u.Path == "." {
		u.Path = "/.json"
	} else {
		u.Path += ".json"
	}
	return u.String()
}

// jsonReader returns a reader (io.Reader) on the JSON representation
// of the value passed as parameter.
func jsonReader(value interface{}) (io.Reader, error) {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r Reference) writeDebug(req *http.Request, response *http.Response) {
	fmt.Fprintln(r.debug, "----- BEGIN DEBUG -----")
	fmt.Fprintf(r.debug, "%v %v\n", req.Method, req.URL)
	dbg := response.Header.Get("X-Firebase-Auth-Debug");
	if (dbg != "") {
	fmt.Fprintf(r.debug, "X-Firebase-Auth-Debug: %v\n", dbg)
	}
	fmt.Fprintln(r.debug, "----- END DEBUG -----")
}

// Value reads from the database and store the content in value. It gives an error
// if it the request fails or if it can't decode the returned payload.
func (r Reference) Value(value interface{}) (err error) {
	req, err := http.NewRequest("GET", r.addAuth().jsonUrl(), nil)
	if err != nil {
		return errors.New(fmt.Sprintf("error while building the request: %v", err))
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("error while executing the request: %v", err))
	}
	defer response.Body.Close()
	if r.debug != nil {
		r.writeDebug(req, response)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("error, response is : %v", response.Status))
	}
	d := json.NewDecoder(response.Body)
	err = d.Decode(value)
	if err != nil {
		return errors.New(fmt.Sprintf("error decoding the result: %v", err))
	}
	return nil
}

// Set write data to the database location given by the Reference r.
// This will overwrite any data at this location and all child locations.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#set
// for more details.
func (r Reference) Set(value interface{}) (err error) {
	b, err := jsonReader(value)
	if err != nil {
		return errors.New(fmt.Sprintf("error reading body: %v", err))
	}
	req, err := http.NewRequest("PUT", r.addAuth().jsonUrl(), b)
	if err != nil {
		return errors.New(fmt.Sprintf("error while building the request: %v", err))
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("error while executing the request: %v", err))
	}
	defer response.Body.Close()
	if r.debug != nil {
		r.writeDebug(req, response)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("error, response is : %v", response.Status))
	}
	return nil
}

// SetWithResult does the same as the Set function and, additionally, stores the
// resulting node in result.
func (r Reference) SetWithResult(value interface{}, result interface{}) (err error) {
	b, err := jsonReader(value)
	if err != nil {
		return errors.New(fmt.Sprintf("error reading body: %v", err))
	}
	req, err := http.NewRequest("PUT", r.addAuth().jsonUrl(), b)
	if err != nil {
		return errors.New(fmt.Sprintf("error while building the request: %v", err))
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("error while executing the request: %v", err))
	}
	defer response.Body.Close()
	if r.debug != nil {
		r.writeDebug(req, response)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("error, response is : %v", response.Status))
	}
	d := json.NewDecoder(response.Body)
	err = d.Decode(result)
	if err != nil {
		return errors.New(fmt.Sprintf("error decoding the result: %v", err))
	}
	return nil
}

// Update writes multiple values to the database at once. The "value" argument contains multiple
// property/value pairs that will be written to the database together. Each child property can
// either be a simple property (for example, "name"), or a relative path (for example, "name/first")
// from the current location to the data to update.
//
// As opposed to the set() method, update() can be use to selectively update only the referenced properties
// at the current location (instead of replacing all the child properties at the current location).
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#update
// for more details.
func (r Reference) Update(value interface{}) (err error) {
	b, err := jsonReader(value)
	if err != nil {
		return errors.New(fmt.Sprintf("error reading body: %v", err))
	}
	req, err := http.NewRequest("PATCH", r.addAuth().jsonUrl(), b)
	if err != nil {
		return errors.New(fmt.Sprintf("error while building the request: %v", err))
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("error while executing the request: %v", err))
	}
	defer response.Body.Close()
	if r.debug != nil {
		r.writeDebug(req, response)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("error, response is : %v", response.Status))
	}
	return nil
}

// UpdateWithResult does the same as the Update function and, additionally, stores the
// updated node in result.
func (r Reference) UpdateWithResult(value interface{}, result interface{}) (err error) {
	b, err := jsonReader(value)
	if err != nil {
		return errors.New(fmt.Sprintf("error reading body: %v", err))
	}
	req, err := http.NewRequest("PATCH", r.addAuth().jsonUrl(), b)
	if err != nil {
		return errors.New(fmt.Sprintf("error while building the request: %v", err))
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("error while executing the request: %v", err))
	}
	defer response.Body.Close()
	if r.debug != nil {
		r.writeDebug(req, response)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("error, response is : %v", response.Status))
	}
	d := json.NewDecoder(response.Body)
	err = d.Decode(result)
	if err != nil {
		return errors.New(fmt.Sprintf("error decoding the result: %v", err))
	}
	return nil
}

// Push generates a new child location using a unique key and returns this key
// in the parameter "name".
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#push
// for more details.
func (r Reference) Push(value interface{}) (name string, err error) {
	b, err := jsonReader(value)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error reading body: %v", err))
	}
	req, err := http.NewRequest("POST", r.addAuth().jsonUrl(), b)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error while building the request: %v", err))
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("error while executing the request: %v", err))

	}
	defer response.Body.Close()
	if r.debug != nil {
		r.writeDebug(req, response)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", errors.New(fmt.Sprintf("error, response is : %v", response.Status))
	}
	result := map[string]string{}
	d := json.NewDecoder(response.Body)
	err = d.Decode(&result)
	if err != nil {
		return "", err
	}
	if name, ok := result["name"]; ok {
		return name, nil
	} else {
		return "", nil
	}
}

// Remove deletes the data at the database location given by the reference r.
// Any data at child locations will also be deleted.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#remove
// for more details.
func (r Reference) Remove() (err error) {
	req, err := http.NewRequest("DELETE", r.addAuth().jsonUrl(), nil)
	if err != nil {
		return errors.New(fmt.Sprintf("error while building the request: %v", err))
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("error while executing the request: %v", err))
	}
	defer response.Body.Close()
	if r.debug != nil {
		r.writeDebug(req, response)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("error, response is : %v", response.Status))
	}
	return nil
}

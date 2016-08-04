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

package firebasedb

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"
)

func (r Reference) WithHttpClient(c *http.Client) Reference {
	result := r
	result.client = c
	return result
}

// jsonUrl is an internal function to build the URL for the REST API
// See https://firebase.google.com/docs/reference/rest/database/ "API Usage"
func (r Reference) jsonUrl() string {
	u := r.url
	u.Path = path.Clean(u.Path)
	if u.Path == "." {
		u.Path = "/.json"
	} else {
		u.Path += ".json"
	}
	return u.String()
}

func jsonReader(value interface{}) (io.Reader, error) {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Values reads from the database and store the content in value. It gives an error
// if it the request fails or if it can't decode the returned payload.
func (r Reference) Value(value interface{}) (err error) {
	req, err := http.NewRequest("GET", r.jsonUrl(), nil)
	if err != nil {
		return err
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}
	d := json.NewDecoder(response.Body)
	return d.Decode(value)
}

func (r Reference) Set(value interface{}, result interface{}) (err error) {
	b, err := jsonReader(value)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", r.jsonUrl(), b)
	if err != nil {
		return err
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}
	if result != nil {
		d := json.NewDecoder(response.Body)
		return d.Decode(result)
	} else {
		return nil
	}
}

func (r Reference) Patch(value interface{}, result interface{}) (err error) {
	b, err := jsonReader(value)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PATCH", r.jsonUrl(), b)
	if err != nil {
		return err
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}
	if result != nil {
		d := json.NewDecoder(response.Body)
		return d.Decode(result)
	} else {
		return nil
	}
}

func (r Reference) Push(value interface{}) (name string, err error) {
	b, err := jsonReader(value)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", r.jsonUrl(), b)
	if err != nil {
		return "", err
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New(response.Status)
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

func (r Reference) Delete() (err error) {
	req, err := http.NewRequest("DELETE", r.jsonUrl(), nil)
	if err != nil {
		return err
	}
	response, err := r.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}
	return nil
}

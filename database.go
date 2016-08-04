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

// Package firebasedb implements a REST client for the Firebase Realtime Database
// (https://firebase.google.com/docs/database/). The API is as close as possible
// to the official JavaScript API.
//
// Similar / related project:
//   https://github.com/zabawaba99/firego
//   https://github.com/cosn/firebase
//
// Reference / documentation:
//   https://firebase.google.com/docs/reference/rest/database
//   https://firebase.google.com/docs/database/rest/structure-data
//   https://firebase.google.com/docs/database/rest/retrieve-data
//   https://firebase.google.com/docs/database/rest/save-data
//   https://firebase.google.com/docs/reference/js/firebase.database.Database
//   https://firebase.google.com/docs/reference/js/firebase.database.Reference
//   https://firebase.google.com/docs/reference/js/firebase.database.Query
//   https://www.firebase.com/docs/rest/api
//
// This packages uses the "Advanced Go Concurrency Patterns" presented by Sameer Ajmani:
//   https://blog.golang.org/advanced-go-concurrency-patterns
package firebasedb

import (
	"net/url"
)

// NewFirebaseDB opens a new Firebase Database connection using the URL u and the
// authentication auth. Currently, only the database secret can be used as auth.
func NewFirebaseDB(u, auth string) (Reference, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return Reference{}, err
	} else {
		ref := Reference{
			url: *parsedUrl,
		}
		if auth != "" {
			ref = ref.withParam("auth", auth)
		}
		return ref, nil
	}
}

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
	"errors"
	"fmt"
	"io"
	"net/http"
	urlLib "net/url"
	pathLib "path"
	"strconv"
	"strings"
)

// Reference represents a specific location in the database and can be used
// for reading or writing data to that database location.
type Reference struct {
	url           urlLib.URL
	Error         error
	client        *http.Client
	auth          Authenticator
	debug         io.Writer
	passKeepAlive bool
	retry         bool
}

// NewReference creates a new Firebase DB reference at url passed as parameter.
func NewReference(url string) Reference {
	parsedUrl, err := urlLib.Parse(url)
	if err != nil {
		return Reference{
			Error: err,
		}
	}
	return Reference{
		url:   *parsedUrl,
		Error: nil,
	}
}

// PassKeepAlive sets the passKeepAlive flag of the Reference. When the references is
// used in the Subscribe() method, the passKeepAlive flag controls the automatic handling
// if keep-alive messages.
func (r Reference) PassKeepAlive(value bool) Reference {
	result := r
	result.passKeepAlive = value
	return result
}

// Retry sets the retry flag for the Reference. When a references has the retry flag set,
// then the library will retry the requests in case of failures.
func (r Reference) Retry(value bool) Reference {
	result := r
	result.retry = value
	return result
}

// httpClient returns the HTTP client from the reference or
// http.DefaultClient if no client has been configured.
func (r Reference) httpClient() *http.Client {
	if r.client == nil {
		return http.DefaultClient
	} else {
		return r.client
	}
}

// withParam is a local function to add an error to a reference.
func (r Reference) withError(err error) Reference {
	result := r
	result.Error = err
	return result
}

// withParam is a local function to add query parameter to the URL of the reference.
func (r Reference) withParam(key, value string) Reference {
	result := r
	q := r.url.Query()
	q.Set(key, value)
	result.url.RawQuery = q.Encode()
	return result
}

// withParam is a local function to add quoted query parameter to the URL. Integer are
// returned as numbers and string are surrounded by double quotes.
func (r Reference) withQuotedParam(key string, value interface{}) Reference {
	var qvalue string = ""
	var err error = nil
	switch i := value.(type) {
	case uint:
		qvalue = strconv.FormatUint(uint64(i), 10)
	case int:
		qvalue = strconv.FormatInt(int64(i), 10)
	case float64:
		qvalue = strconv.FormatFloat(i, 'f', -1, 64)
	case string:
		qvalue = fmt.Sprintf(`"%s"`, strings.Trim(i, `"`))
	default:
		err = errors.New("Invalid Type")
	}
	if err == nil {
		return r.withParam(key, qvalue)
	} else {
		return r.withError(err)
	}
}

// Ref returns a reference to the root or the specified path.
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#ref
// or https://firebase.google.com/docs/reference/js/firebase.database.Database#ref
// for more details.
func (r Reference) Ref(path string) Reference {
	result := r
	result.url.Path = pathLib.Clean(pathLib.Join("/", path))
	return result
}

// RefFromUrl returns a reference to the root or the path specified in url.
// err is set if the host of the url is not the same as the current database.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Database#refFromURL
// for more details.
func (r Reference) RefFromUrl(url urlLib.URL) Reference {
	if r.url.Host != url.Host {
		return r.withError(errors.New("The URL has not the same host as the current database"))
	} else {
		return r.Ref(url.Path)
	}
}

// Rules returns a reference to the rules settings of the database.
func (r Reference) Rules() Reference {
	return r.Ref(".settings/rules")
}

func (r Reference) Debug(w io.Writer) Reference {
	result := r
	result.debug = w
	return result
}

// Auth authenticates the request to allow access to data protected by Firebase Realtime Database Rules.
// The argument is an object that implements the Authenticator interface. The String() method can either
// returns a Firebase app's secret or an authentication token.
//
// Note that when the reference is used in a streaming submission, a "auth_revoked" event will trigger
// a re-authentication, and reopen the http connection. *This will result in an additional "put" event*.
//
// See https://firebase.google.com/docs/reference/rest/database/#section-param-auth
// and https://firebase.google.com/docs/reference/rest/database/user-auth
// for more details.
func (r Reference) Auth(auth Authenticator) Reference {
	result := r
	result.auth = auth
	return result
}

// Shallow is an advanced feature, designed to help you work with large datasets without
// needing to download everything. Set this to true to limit the depth of the data returned
// at a location. If the data at the location is a JSON primitive (string, number or boolean),
// its value will simply be returned. If the data snapshot at the location is a JSON object,
// the values for each key will be truncated to true.
//
// See https://firebase.google.com/docs/reference/rest/database/#section-param-shallow
// for more details.
func (r Reference) Shallow() Reference {
	return r.withParam("shallow", "true")
}

// Pretty is used to view the data in a human-readable format. This is usually only used
// for debugging purposes.
//
// See https://firebase.google.com/docs/reference/rest/database/#section-param-print
// for more details.
func (r Reference) Pretty() Reference {
	return r.withParam("print", "pretty")
}

// Silent is used to suppress the output from the server when writing data. The resulting
// response will be empty and indicated by a 204 No Content HTTP status code.
//
// See https://firebase.google.com/docs/reference/rest/database/#section-param-print
// for more details.
func (r Reference) Silent() Reference {
	return r.withParam("print", "silent")
}

// Export returns a reference that include priority information in the response.
//
// See https://firebase.google.com/docs/reference/rest/database/#section-param-format
// for more details.
func (r Reference) Export() Reference {
	return r.withParam("format", "export")
}

// Key returns the last part of the current path.
// For example, "ada" is the key for https://sample-app.firebaseio.com/users/ada.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#key
// for more detail.
func (r Reference) Key() string {
	p := pathLib.Base(pathLib.Clean(r.url.Path))
	if p == "." || p == "/" {
		return ""
	} else {
		return p
	}
}

// Parent returns the parent location of a reference.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#parent
// for more details.
func (r Reference) Parent() Reference {
	result := r
	result.url.Path = pathLib.Clean(pathLib.Join(result.url.Path, ".."))
	return result
}

// Root returns the root location of a reference.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#root
// for more details.
func (r Reference) Root() Reference {
	result := r
	result.url.Path = "/"
	return result
}

// Childs returns a reference for the location at the specified relative path.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#child
// for more details.
func (r Reference) Child(path string) Reference {
	result := r
	result.url.Path = pathLib.Clean(pathLib.Join(result.url.Path, path))
	return result
}

// String returns the absolute URL for this location.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#toString
// for more details.
func (r Reference) String() string {
	return r.url.String()
}

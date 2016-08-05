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
	"errors"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"net/http"
)

// Reference represents a specific location in the database and can be used
// for reading or writing data to that database location.
type Reference struct {
	url url.URL
	err error
	client *http.Client
	skipKeepAlive bool
	retry bool
}

// SkipKeepAlive sets the skipKeepAlive flag for the Reference. When the references is
// used in the Subscribe() method, the skipKeepAlive flag controls the automatic handling
// if keep-alive messages.
func (r Reference) SkipKeepAlive(x bool) Reference {
	result := r
	r.skipKeepAlive = x
	return result
}

// Retry sets the retry flag for the Reference. When a references has the retry flag set,
// then the library will retry the requests in case of failures.
func (r Reference) Retry(x bool) Reference {
	result := r
	r.retry = x
	return result
}


// Error returns the error from a reference. Note that an error is set as soon as
// something wrong occurs with reference operations and is never reset.
func (r Reference) Error() error {
	return r.err
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
	result.err = err
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
func (r Reference) Ref(p string) Reference {
	result := r
	result.url.Path = path.Clean(path.Join("/", p))
	return result
}

// RefFromUrl returns a reference to the root or the path specified in url.
// err is set if the host of the url is not the same as the current database.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Database#refFromURL
// for more details.
func (r Reference) RefFromUrl(u url.URL) Reference {
	if r.url.Host != u.Host {
		return r.withError(errors.New("The URL has not the same host as the current database"))
	} else {
		return r.Ref(u.Path)
	}
}

// Rules returns a reference to the rules settings of the database.
func (r Reference) Rules() Reference {
	return r.Ref(".settings/rules")
}

// Auth authenticates the request to allow access to data protected by Firebase Realtime Database Rules.
// The argument can either be your Firebase app's secret or an authentication token
//
// See https://firebase.google.com/docs/reference/rest/database/#section-param-auth
// and https://firebase.google.com/docs/reference/rest/database/user-auth
// for more details.
func (r Reference) Auth(auth string) Reference {
	return r.withParam("auth", auth)
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
	p := path.Base(path.Clean(r.url.Path))
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
	result.url.Path = path.Clean(path.Join(result.url.Path, ".."))
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
func (r Reference) Child(p string) Reference {
	result := r
	result.url.Path = path.Clean(path.Join(result.url.Path, p))
	return result
}

// String returns the absolute URL for this location.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#toString
// for more details.
func (r Reference) String() string {
	return r.url.String()
}
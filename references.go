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
	"net/url"
	"path"
)

type Reference struct { // TODO: check if we can replace by a simple URL
	url url.URL
}

// withParam is a local function to add query parameter to the URL.
func (r Reference) withParam(key, value string) Reference {
	result := r
	q := r.url.Query()
	q.Set(key, value)
	result.url.RawQuery = q.Encode()
	return result
}

// Ref returns a reference to the root or the specified path.
func (r Reference) Ref(p string) Reference {
	result := r
	result.url.Path = path.Clean(path.Join("/", p))
	return result
}

// RefFromUrl returns a reference to the root or the path specified in url.
// err is set if the host of the url is not the same as the current database.
func (r Reference) RefFromUrl(u url.URL) (Reference, error) {
	if r.url.Host != u.Host {
		return r, errors.New("The URL has not the same host as the current database")
	}
	return r.Ref(u.Path), nil
}

func (r Reference) Shallow() Reference {
	return r.withParam("shallow", "true")
}

func (r Reference) Pretty() Reference {
	return r.withParam("print", "pretty")
}

func (r Reference) Silent() Reference {
	return r.withParam("print", "silent")
}

func (r Reference) Export() Reference {
	return r.withParam("format", "export")
}

func (r Reference) Key() string {
	p := path.Base(path.Clean(r.url.Path))
	if p == "." || p == "/" {
		return ""
	} else {
		return p
	}
}

func (r Reference) Parent() Reference {
	result := r
	result.url.Path = path.Clean(path.Join(result.url.Path, ".."))
	return result
}

func (r Reference) Root() Reference {
	result := r
	result.url.Path = "/"
	return result
}

func (r Reference) Child(p string) Reference {
	result := r
	result.url.Path = path.Clean(path.Join(result.url.Path, p))
	return result
}

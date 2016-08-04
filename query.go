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
	"strconv"
)

// OrderByChild generates a new query  ordered by the specified child key.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#orderByChild
// or https://firebase.google.com/docs/reference/js/firebase.database.Query#orderByChild
// for more details
func (r Reference) OrderByChild(childKey string) Reference {
	return r.withQuotedParam("orderBy", childKey)
}

// OrderedByKey generates a new query ordered by key.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#orderByKey
// https://firebase.google.com/docs/reference/js/firebase.database.Query#orderByKey
// for more details
func (r Reference) OrderByKey() Reference {
	return r.OrderByChild("$key")
}

// OrderByValue generates a new query ordered by child values.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#orderByValue
// or https://firebase.google.com/docs/reference/js/firebase.database.Query#orderByValue
// for more details
func (r Reference) OrderByValue() Reference {
	return r.OrderByChild("$value")
}

// LimitToFirst generates a new query limited to the first specific number of children.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#limitToFirst
// or https://firebase.google.com/docs/reference/js/firebase.database.Query#limitToFirst
// for more details
func (r Reference) LimitToFirst(n uint64) Reference {
	return r.withParam("limitToFirst", strconv.FormatUint(n, 10))
}

// LimitToLast generates a new query limited to the last specific number of children.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#limitToLast
// or https://firebase.google.com/docs/reference/js/firebase.database.Query#limitToLast
// for more details
func (r Reference) LimitToLast(n uint64) Reference {
	return r.withParam("limitToLast", strconv.FormatUint(n, 10))
}

// StartAt creates a query with the specified starting point.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#startAt
// or https://firebase.google.com/docs/reference/js/firebase.database.Query#startAt
// for more details.
func (r Reference) StartAt(n interface{}) Reference {
	return r.withQuotedParam("startAt", n)
}

// EndAt creates a query with the specified ending point.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#endAt
// or https://firebase.google.com/docs/reference/js/firebase.database.Query#endAt
// for more details
func (r Reference) EndAt(n interface{}) Reference {
	return r.withQuotedParam("endAt", n)
}

// EqualTo creates a query which includes children which match the specified value.
//
// See https://firebase.google.com/docs/reference/js/firebase.database.Reference#equalTo
// or https://firebase.google.com/docs/reference/js/firebase.database.Query#equalTo
// for more details
func (r Reference) EqualTo(n interface{}) Reference {
	return r.withQuotedParam("equalTo", n)
}

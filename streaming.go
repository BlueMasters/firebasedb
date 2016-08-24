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

// This code implements the techniques explained in "Advanced Go
// Concurrency Patterns" from Sameer Ajmani.
// https://blog.golang.org/advanced-go-concurrency-patterns

package firebasedb

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Event is the type used to represent streaming events. The type of the event
// can be read directly from the type. The data is extracted using the Value method
//
// See https://firebase.google.com/docs/reference/rest/database/#section-streaming
// for more details.
type Event struct {
	Type string // can be put, patch, keep-alive, cancel, or auth_revoked
	Err  error
	data string
}

// Value unmarshals data from an event. It returns the data in v and the path
// of the event in the path return attribute.
func (e Event) Value(v interface{}) (path string, err error) {
	var p struct {
		Path string      `json:"path"`
		Data interface{} `json:"data"`
	}
	p.Data = &v
	err = json.Unmarshal([]byte(e.data), &p)
	if err != nil {
		path = ""
	} else {
		path = p.Path
	}
	return path, err
}

// Subscription is the interface for event subscriptions. Subscriptions
// are returned by the Subscribe method.
type Subscription struct {
	reader        io.ReadCloser // from the HTTP request's body
	reference     *Reference    // copy of the reference
	events        chan Event    // sends events to the user
	closing       chan bool     // for Close
	LastKeepAlive time.Time
}

func (r Reference) openStream() (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", r.addAuth().jsonUrl(), nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error while building the request: %v", err))
	}
	req.Header.Add("Accept", "text/event-stream")
	response, err := r.do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error while executing the request: %v", err))
	}
	if response.StatusCode != 200 {
		response.Body.Close()
		return nil, errors.New(fmt.Sprintf("error, response is : %v", response.Status))
	}
	return response.Body, nil
}

// Subscribe returns a subscription on the reference. The returned subscription
// is used to access the streamed events.
func (r Reference) Subscribe() (*Subscription, error) {
	reader, err := r.openStream()
	if err != nil {
		return nil, err
	}
	s := &Subscription{
		reader:    reader,
		reference: &r,
		events:    make(chan Event), // for Events
		closing:   make(chan bool),  // for Close
	}
	go s.loop()
	return s, nil
}

// Events returns the event channel from the subscription.
func (s *Subscription) Events() <-chan Event {
	return s.events
}

// Close closes the subscription and finishes the request.
func (s *Subscription) Close() error {
	return s.reader.Close()
}

// main loop
func (s *Subscription) loop() {

	var fetchEvent = make(chan Event)
	defer close(fetchEvent)
	var pending []Event

	go func() { // read the payload and feed the fetchEvent channel
		payload := make([]string, 2)
		lineCount := 0
		r := bufio.NewReader(s.reader)
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				break
			}
			line = strings.Trim(line, " \r\n")
			if len(line) == 0 {
				// empty line
				if lineCount == len(payload) {
					if !strings.HasPrefix(payload[0], "event:") {
						fetchEvent <- Event{
							Err: errors.New("First line does not start with event:"),
						}
					} else if !strings.HasPrefix(payload[1], "data:") {
						fetchEvent <- Event{
							Err: errors.New("Second line does not start with data:"),
						}
					} else {
						eventType := strings.Trim(strings.TrimPrefix(payload[0], "event:"), " \r\n")
						eventData := strings.Trim(strings.TrimPrefix(payload[1], "data:"), " \r\n")
						switch eventType {
						case "keep-alive":
							s.LastKeepAlive = time.Now()
							if s.reference.passKeepAlive {
								fetchEvent <- Event{Type: eventType, data: eventData, Err: nil}
							}
						case "auth_revoked":
							var err error = nil
							if s.reference.auth != nil {
								if err = s.reference.auth.Renew(); err == nil {
									s.reader.Close()
									s.reader, err = s.reference.openStream()
									if err == nil {
										r = bufio.NewReader(s.reader)
										break // everything is OK, no need to send the event further.
									}
								}
							}
							// send the event with the proper error code.
							fetchEvent <- Event{Type: eventType, data: eventData, Err: err}
						default: // send "normal" event
							fetchEvent <- Event{Type: eventType, data: eventData, Err: nil}
						}
					}
				} else {
					fetchEvent <- Event{Err: errors.New("Badly formated body")}
				}
				lineCount = 0
			} else { // line is not empty
				if lineCount < len(payload) {
					payload[lineCount] = line
					lineCount++
				}
			}
		}
		s.closing <- true
	}()

	for {
		var first Event
		var events chan Event
		if len(pending) > 0 {
			first = pending[0]
			events = s.events // enable send case
		}

		select {
		case event := <-fetchEvent:
			// Currently, I am not controlling the size of the pending queue.
			// But the structure of this program enables those check if required.
			pending = append(pending, event)
		case <-s.closing:
			close(s.events)
			break
		case events <- first:
			pending = pending[1:]
		}
	}
}

# Firebase Realtime Database GO client

Package firebasedb implements a REST client for the
[Firebase Realtime Database](https://firebase.google.com/docs/database/).
The API is as close as possible to the official JavaScript API.

[![GoDoc](https://godoc.org/github.com/BlueMasters/firebasedb?status.svg)](https://godoc.org/github.com/spf13/hugo)
[![Travis](https://img.shields.io/travis/BlueMasters/firebasedb.svg)](https://travis-ci.org/BlueMasters/firebasedb)
[![Made in Switzerland](https://img.shields.io/badge/Made%20with%20♥%20in-Fribourg%20%2F%20Switzerland-blue.svg)](http://fribourg.ch/fr/)

## Credits
* Steven Berlanga for [another implementation](https://github.com/zabawaba99/firego) of
  Firebase in go. I also used some tricks from his travis config.
* Sameer Ajmani for his presentation [Advanced Go Concurrency Patterns](https://blog.golang.org/advanced-go-concurrency-patterns)
  that I used to implement the Streams.
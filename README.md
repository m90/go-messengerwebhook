# go-messengerwebhook
[![Build Status](https://travis-ci.org/m90/go-messengerwebhook.svg?branch=master)](https://travis-ci.org/m90/go-messengerwebhook)
[![godoc](https://godoc.org/github.com/m90/go-messengerwebhook?status.svg)](http://godoc.org/github.com/m90/go-messengerwebhook)

> setup a webhook for facebook messenger and subscribe to its updates

## Installation

Install using `go get`:

```sh
$go get github.com/m90/go-messengerwebhook
```

## Usage

Calling `SetupWebhook(verifyToken string)` returns a `http.HandlerFunc` and a `<-chan msngrhook.Update`:

```go
handler, updates := msngrhook.SetupWebhook("my_verify_token")
go http.ListenAndServe("0.0.0.0:3000", handler)

for update := range updates {
	if update.Error != nil {
		// handle error
	} else {
		// handle update
	}
}
```

## Tests

Run the tests:

```sh
$ make
```

### License
MIT © [Frederik Ring](http://www.frederikring.com)

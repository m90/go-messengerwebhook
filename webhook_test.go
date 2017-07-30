package msngrhook

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetupWebhook_Get(t *testing.T) {
	tests := []struct {
		name               string
		requestURL         string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			"pass",
			"?hub.verify_token=abc123&hub.challenge=xyz",
			200,
			"xyz",
		},
		{
			"bad token",
			"?hub.verify_token=123abc&hub.challenge=xyz",
			401,
			"Verify token did not match\n",
		},
		{
			"bad token only",
			"?hub.challenge=xyz",
			401,
			"Verify token did not match\n",
		},
		{
			"empty token",
			"",
			401,
			"Verify token did not match\n",
		},
		{
			"missing challenge",
			"?hub.verify_token=abc123",
			400,
			"Missing hub.challenge parameter\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler, _ := SetupWebhook("abc123")
			ts := httptest.NewServer(handler)
			res, resErr := http.DefaultClient.Get(fmt.Sprintf("%v%v", ts.URL, test.requestURL))
			if resErr != nil {
				t.Error(resErr)
			}
			if res.StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code of %v, got %v", test.expectedStatusCode, res.StatusCode)
			}
			bytes, _ := ioutil.ReadAll(res.Body)
			if string(bytes) != test.expectedBody {
				t.Errorf("Expected body of %v, got %v", test.expectedBody, string(bytes))
			}
		})
	}
}

func TestSetupWebhook_Post(t *testing.T) {
	tests := []struct {
		name                 string
		payload              string
		expectedStatusCode   int
		expectError          bool
		expectedFirstMessage string
	}{
		{
			"bad payload",
			"testdata/invalid.json",
			500,
			true,
			"",
		},
		{
			"default",
			"testdata/payload.json",
			200,
			false,
			"Gophers can you hear me?",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler, updates := SetupWebhook("abc123")
			done := make(chan bool)
			var err error
			var message string
			go func() {
				for u := range updates {
					if u.Error != nil {
						err = u.Error
					}
					if u.Message != nil {
						message = u.Message.Text
					}
					break
				}
				done <- true
			}()
			ts := httptest.NewServer(handler)
			content, _ := ioutil.ReadFile(test.payload)
			res, _ := http.DefaultClient.Post(ts.URL, "application/json; charset=utf-8", strings.NewReader(string(content)))

			if res.StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code of %v, got %v", test.expectedStatusCode, res.StatusCode)
			}
			<-done
			if test.expectError {
				if err == nil {
					t.Error("Expected error to be not nil")
				}
			} else {
				if message != test.expectedFirstMessage {
					t.Errorf("Expected first message to be %v, got %v", test.expectedFirstMessage, message)
				}
			}
		})
	}
}

func TestSetupWebhook_Other(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		expectedStatusCode int
	}{
		{"delete", http.MethodDelete, 405},
		{"put", http.MethodPut, 405},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler, _ := SetupWebhook("abc123")
			ts := httptest.NewServer(handler)
			req, _ := http.NewRequest(test.method, ts.URL, nil)
			res, _ := http.DefaultClient.Do(req)
			if res.StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code of %v, got %v", test.expectedStatusCode, res.StatusCode)
			}
		})
	}
}

type badReader int

func (b badReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("I'm a bad reader")
}
func TestSetupWebhook_BadReader(t *testing.T) {
	t.Run("bad reader", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", badReader(0))
		w := httptest.NewRecorder()
		handler, _ := SetupWebhook("abc123")
		handler.ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Unexpected response code %v", w.Code)
		}
	})
}

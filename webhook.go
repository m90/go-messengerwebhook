package msngrhook

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	queryKeyChallenge       = "hub.challenge"
	queryKeyToken           = "hub.verify_token"
	messageMissingChallenge = "Missing hub.challenge parameter"
	messageTokenMismatch    = "Verify token did not match"
	responseBodyOK          = "OK"
)

// SetupWebhook creates a http.HandlerFunc and a channel of updates
// using the given verify token string
func SetupWebhook(verifyToken string) (http.HandlerFunc, <-chan Update) {
	updates := make(chan Update)

	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if token := r.URL.Query().Get(queryKeyToken); token != verifyToken {
				http.Error(w, messageTokenMismatch, http.StatusUnauthorized)
				return
			}
			if challenge := r.URL.Query().Get(queryKeyChallenge); challenge != "" {
				io.WriteString(w, challenge)
			} else {
				http.Error(w, messageMissingChallenge, http.StatusBadRequest)
			}
		case http.MethodPost:
			bytes, bytesErr := ioutil.ReadAll(r.Body)
			if bytesErr != nil {
				http.Error(w, bytesErr.Error(), http.StatusBadRequest)
				return
			}
			defer r.Body.Close()
			update := &UpdateRequest{}
			if unmarshalErr := json.Unmarshal(bytes, &update); unmarshalErr != nil {
				http.Error(w, unmarshalErr.Error(), http.StatusInternalServerError)
				updates <- Update{Error: unmarshalErr}
				return
			}
			for _, entry := range *update.Entry {
				for _, messaging := range *entry.Messaging {
					updates <- messaging
				}
			}
			io.WriteString(w, responseBodyOK)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
	return handler, updates
}

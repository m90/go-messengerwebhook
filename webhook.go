package msngrhook

import (
	"encoding/json"
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
				w.Write([]byte(challenge))
			} else {
				http.Error(w, messageMissingChallenge, http.StatusBadRequest)
			}
		case http.MethodPost:
			update := UpdateRequest{}
			if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				updates <- Update{Error: err}
				return
			}

			for _, entry := range *update.Entry {
				for _, messaging := range *entry.Messaging {
					updates <- messaging
				}
			}
			w.Write([]byte(responseBodyOK))
		default:
			http.Error(
				w,
				http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed,
			)
		}
	}
	return handler, updates
}

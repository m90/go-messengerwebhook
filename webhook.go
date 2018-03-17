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

type hook struct {
	token   string
	updates chan Update
}

// New creates a http.Handler and a channel of updates
// using the given verify token string
func New(verifyToken string) (http.Handler, <-chan Update) {
	updates := make(chan Update)
	handler := hook{token: verifyToken, updates: updates}
	return &handler, updates
}

func (h *hook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if token := r.URL.Query().Get(queryKeyToken); token != h.token {
			http.Error(w, messageTokenMismatch, http.StatusUnauthorized)
			return
		}
		if challenge := r.URL.Query().Get(queryKeyChallenge); challenge != "" {
			w.Write([]byte(challenge))
		} else {
			http.Error(w, messageMissingChallenge, http.StatusBadRequest)
		}
	case http.MethodPost:
		update := updateRequest{}
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			h.updates <- Update{Error: err}
			return
		}

		for _, entry := range *update.Entry {
			for _, messaging := range *entry.Messaging {
				h.updates <- messaging
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

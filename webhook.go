package msngrhook

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// SetupWebhook creates a http.HandlerFunc and a channel of updates
// using the given verify token string
func SetupWebhook(verifyToken string) (http.HandlerFunc, <-chan Update) {
	updates := make(chan Update, 1)

	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			token := r.URL.Query().Get("hub.verify_token")
			if token != verifyToken {
				http.Error(w, "Verify token did not match", http.StatusUnauthorized)
				return
			}
			if challenge := r.URL.Query().Get("hub.challenge"); challenge != "" {
				io.WriteString(w, challenge)
			} else {
				http.Error(w, "Missing hub.challenge parameter", http.StatusBadRequest)
			}
		case http.MethodPost:
			bytes, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			update := &UpdateRequest{}
			unmarshalErr := json.Unmarshal(bytes, &update)
			if unmarshalErr != nil {
				out := Update{Error: unmarshalErr}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				updates <- out
				return
			}
			if len(*update.Entry) == 0 {
				return
			}
			for _, entry := range *update.Entry {
				for _, messaging := range *entry.Messaging {
					updates <- messaging
				}
			}
			io.WriteString(w, "OK")
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
	return handler, updates
}

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
	updates := make(chan Update)

	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if token := r.URL.Query().Get("hub.verify_token"); token != verifyToken {
				http.Error(w, "Verify token did not match", http.StatusUnauthorized)
				return
			}
			if challenge := r.URL.Query().Get("hub.challenge"); challenge != "" {
				io.WriteString(w, challenge)
			} else {
				http.Error(w, "Missing hub.challenge parameter", http.StatusBadRequest)
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
			io.WriteString(w, "OK")
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
	return handler, updates
}

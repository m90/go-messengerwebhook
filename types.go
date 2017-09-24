package msngrhook

import (
	"fmt"
	"net/url"
)

// UpdateRequest describes the request body's top level wrapper
type UpdateRequest struct {
	Object string         `json:"object"`
	Entry  *[]UpdateEntry `json:"entry"`
}

// UpdateEntry describes an entry contained in the top level wrapper
type UpdateEntry struct {
	ID        string    `json:"id"`
	Time      int       `json:"time"`
	Messaging *[]Update `json:"messaging"`
}

// Update contains the actual data that an update will be composed of
type Update struct {
	Sender    *UpdateSender    `json:"sender"`
	Recipient *UpdateRecipient `json:"recipient"`
	Message   *UpdateMessage   `json:"message"`
	Timestamp int              `json:"timestamp"`
	Postback  *UpdatePostback  `json:"postback,omitempty"`
	Error     error            `json:"-"`
}

// IsPostback checks if the update is a postback
func (u *Update) IsPostback() bool {
	return u.Postback != nil
}

// NormalizedTextMessage returns the applicable text message of an update
// depending on if it is a postback, standard message or an image
func (u *Update) NormalizedTextMessage() string {
	if u.IsPostback() {
		return u.Postback.Payload
	}
	if u.Message == nil {
		return ""
	}
	if u.Message.Attachments != nil {
		for _, a := range *u.Message.Attachments {
			switch a.Type {
			case "location":
				if coords, ok := a.Payload["coordinates"]; ok {
					if cast, ok := coords.(map[string]interface{}); ok {
						return fmt.Sprintf("%v, %v", cast["lat"], cast["long"])
					}
				}
			case "template":
				if templateType, ok := a.Payload["template_type"]; !ok || templateType != "generic" {
					break
				}
				elements, ok := a.Payload["elements"].([]interface{})
				if !ok {
					break
				}
				for _, element := range elements {
					castElement, ok := element.(map[string]interface{})
					if !ok {
						continue
					}
					buttons, ok := castElement["buttons"].([]interface{})
					if !ok {
						continue
					}
					for _, button := range buttons {
						castButton, ok := button.(map[string]interface{})
						if !ok {
							continue
						}
						t, ok := castButton["type"].(string)
						if !ok && t != "element_share" {
							continue
						}
						for _, button := range buttons {
							castButton, ok := button.(map[string]interface{})
							if !ok {
								continue
							}
							if u, ok := castButton["url"].(string); ok {
								return u
							}
						}
					}
				}
			default:
				if value, ok := a.Payload["url"]; ok {
					if urlStr, ok := value.(string); ok {
						u, uErr := url.Parse(urlStr)
						if uErr != nil {
							continue
						}
						return fmt.Sprintf("%v://%v%v", u.Scheme, u.Host, u.Path)
					}
				}
			}
		}
	}
	return u.Message.Text

}

// UpdatePostback contains the postback payload of an update
type UpdatePostback struct {
	Payload string `json:"payload"`
}

// UpdateMessage describes the message data that an update contains
type UpdateMessage struct {
	Text        string              `json:"text,omitempty"`
	MID         string              `json:"mid,omitempty"`
	Attachments *[]UpdateAttachment `json:"attachments,omitempty"`
	QuickReply  *UpdateQuickReply   `json:"quick_reply,omitempty"`
}

// UpdateAttachment describes an update's attachment
type UpdateAttachment struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// UpdateQuickReply contains the payload of a quick reply
type UpdateQuickReply struct {
	Payload string `json:"payload"`
}

// UpdateSender contains the ID of an update's sender
type UpdateSender struct {
	ID string `json:"id"`
}

// UpdateRecipient contains the ID of an update's sender
type UpdateRecipient struct {
	ID string `json:"id"`
}

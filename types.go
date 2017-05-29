package msngrhook

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
// depending on if it is a postback or standard message
func (u *Update) NormalizedTextMessage() string {
	if u.IsPostback() {
		return u.Postback.Payload
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
	Type    string `json:"type"`
	Payload string `json:"payload"`
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

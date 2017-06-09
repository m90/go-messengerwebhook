package msngrhook

import "testing"

func TestIsPostback(t *testing.T) {
	tests := []struct {
		update         Update
		expectedResult bool
	}{
		{
			Update{
				Postback: &UpdatePostback{"postback!"},
			},
			true,
		},
		{
			Update{
				Message: &UpdateMessage{
					Text: "message!",
					MID:  "some-value",
				},
			},
			false,
		},
	}
	for _, test := range tests {
		if test.update.IsPostback() != test.expectedResult {
			t.Errorf(
				"Expected result of %v, got %v",
				test.expectedResult,
				test.update.IsPostback(),
			)
		}
	}
}

func TestNormalizedTextMessage(t *testing.T) {
	tests := []struct {
		update          Update
		expectedMessage string
	}{
		{
			Update{
				Postback: &UpdatePostback{"postback!"},
			},
			"postback!",
		},
		{
			Update{
				Message: &UpdateMessage{
					Text: "message!",
					MID:  "some-value",
				},
			},
			"message!",
		},
		{
			Update{
				Message: &UpdateMessage{
					Attachments: &[]UpdateAttachment{
						UpdateAttachment{
							Type: "image",
							Payload: map[string]interface{}{
								"url": "https://scontent.xx.fbcdn.net/v/t39.1997-6/p100x100/851582_369239386556143_1497813874_n.png?_nc_ad=z-m&oh=52af86654b8cdb071a0d23c2c0208e88&oe=59E9AA4D",
							},
						},
					},
				},
			},
			"https://scontent.xx.fbcdn.net/v/t39.1997-6/p100x100/851582_369239386556143_1497813874_n.png",
		},
	}
	for _, test := range tests {
		if test.update.NormalizedTextMessage() != test.expectedMessage {
			t.Errorf(
				"Expected result of %v, got %v",
				test.expectedMessage,
				test.update.NormalizedTextMessage(),
			)
		}
	}
}

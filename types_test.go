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

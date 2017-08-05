package msngrhook

import "testing"

func TestIsPostback(t *testing.T) {
	tests := []struct {
		name           string
		update         Update
		expectedResult bool
	}{
		{
			"true",
			Update{
				Postback: &UpdatePostback{"postback!"},
			},
			true,
		},
		{
			"false",
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
		t.Run(test.name, func(t *testing.T) {
			if test.update.IsPostback() != test.expectedResult {
				t.Errorf(
					"Expected result of %v, got %v",
					test.expectedResult,
					test.update.IsPostback(),
				)
			}
		})
	}
}

func TestNormalizedTextMessage(t *testing.T) {
	tests := []struct {
		name            string
		update          Update
		expectedMessage string
	}{
		{
			"postback",
			Update{
				Postback: &UpdatePostback{"postback!"},
			},
			"postback!",
		},
		{
			"default",
			Update{
				Message: &UpdateMessage{
					Text: "message!",
					MID:  "some-value",
				},
			},
			"message!",
		},
		{
			"invalid URL",
			Update{
				Message: &UpdateMessage{
					Text: "message!",
					MID:  "some-value",
					Attachments: &[]UpdateAttachment{
						UpdateAttachment{
							Type: "image",
							Payload: map[string]interface{}{
								"url": "%%%%%%%%%%%%%%%%%%%%%%%%%%%",
							},
						},
					},
				},
			},
			"message!",
		},
		{
			"empty",
			Update{},
			"",
		},
		{
			"image attachment",
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
		{
			"location attachment",
			Update{
				Message: &UpdateMessage{
					Attachments: &[]UpdateAttachment{
						UpdateAttachment{
							Type: "location",
							Payload: map[string]interface{}{
								"coordinates": map[string]interface{}{
									"lat":  52.520007,
									"long": 13.404954,
								},
							},
						},
					},
				},
			},
			"52.520007, 13.404954",
		},
		{
			"share button",
			Update{
				Message: &UpdateMessage{
					Attachments: &[]UpdateAttachment{
						UpdateAttachment{
							Type: "template",
							Payload: map[string]interface{}{
								"elements": []map[string]interface{}{
									{
										"buttons": []map[string]interface{}{
											{"title": "Open me", "type": "web_url", "url": "http://m.me/foo.bar"},
											{"type": "element_share"},
										},
									},
									{
										"image_url": "http://www.example.net/bar.jpg",
									},
									{
										"title": "some super interesting thing",
									},
								},
								"sharable":      true,
								"template_type": "generic",
							},
						},
					},
				},
			},
			"http://m.me/foo.bar",
		},
		{
			"unknown template type",
			Update{
				Message: &UpdateMessage{
					Attachments: &[]UpdateAttachment{
						UpdateAttachment{
							Type: "template",
							Payload: map[string]interface{}{
								"template_type": "super new",
							},
						},
					},
				},
			},
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.update.NormalizedTextMessage() != test.expectedMessage {
				t.Errorf(
					"Expected result of %v, got %v",
					test.expectedMessage,
					test.update.NormalizedTextMessage(),
				)
			}
		})
	}
}

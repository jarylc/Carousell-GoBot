package responses

// {"data":{"fieldset_data":{"biography":"Normally deal at Choa Chu Kang or downtown line.\n\nhttps://jarylchng.com","follower_following":{"items":[{"action":{"deep_link":"carousell://followers/344194","type":"go_to_deep_link"},"header":"102","id":"followers","text":"Followers"},{"action":{"deep_link":"carousell://following/344194","type":"go_to_deep_link"},"header":"5","id":"following","text":"Following"}]},"response_rate":"H","trust_badges":{"items":null}}}}

type Profile struct {
	Data struct {
		FieldsetData struct {
			Biography         string `json:"biography"`
			FollowerFollowing struct {
				Items []struct {
					Action struct {
						DeepLink string `json:"deep_link"`
						Type     string `json:"type"`
					} `json:"action"`
					Header string `json:"header"`
					ID     string `json:"id"`
					Text   string `json:"text"`
				} `json:"items"`
			} `json:"follower_following"`
			ResponseRate string `json:"response_rate"`
			TrustBadges  struct {
				Items interface{} `json:"items"`
			} `json:"trust_badges"`
		} `json:"fieldset_data"`
	} `json:"data"`
}

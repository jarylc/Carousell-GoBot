package responses

// {"data":{"token":"akjhiwpe5423097fakesscvtoken21kmxsdfb3gd"}}

type Token struct {
	Data struct {
		Token  string `json:"token"`
		Detail string `json:"detail"`
	} `json:"data"`
}

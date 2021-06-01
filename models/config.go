package models

type Config struct {
	Carousell struct {
		Cookie       string  `yaml:"cookie"`
		PingInterval int16   `yaml:"ping_interval"`
		LowBall      float64 `yaml:"low_ball"`
	} `yaml:"carousell"`
	MessageTemplates struct {
		FAQ           string `yaml:"faq"`
		Initial       string `yaml:"initial"`
		NotAvailable  string `yaml:"not_available"`
		Offered       string `yaml:"offered"`
		PossibleOffer string `yaml:"possible_offer"`
		LowerOffer    string `yaml:"lower_offer"`
		LowBalled     string `yaml:"low_balled"`
		Reminder      string `yaml:"reminder"`
	} `yaml:"message_templates"`
	Reminders     []int16 `yaml:"reminders"`
	CommandPrefix string  `yaml:"command_prefix"`
	StatePrune    int16   `yaml:"state_prune"`
	Forwarders    []struct {
		Type             string `yaml:"type"`
		Token            string `yaml:"token"`
		WebhookURL       string `yaml:"webhook_url"`
		ChatID           string `yaml:"chat_id"`
		MessageTemplates struct {
			Standard string `yaml:"standard"`
			Reminder string `yaml:"reminder"`
		} `yaml:"message_templates"`
	} `yaml:"forwarders"`
}

func DefaultConfig() *Config {
	return &Config{
		Carousell: struct {
			Cookie       string  `yaml:"cookie"`
			PingInterval int16   `yaml:"ping_interval"`
			LowBall      float64 `yaml:"low_ball"`
		}{
			Cookie:       "",
			PingInterval: 300,
			LowBall:      0.7,
		},
		MessageTemplates: struct {
			FAQ           string `yaml:"faq"`
			Initial       string `yaml:"initial"`
			NotAvailable  string `yaml:"not_available"`
			Offered       string `yaml:"offered"`
			PossibleOffer string `yaml:"possible_offer"`
			LowerOffer    string `yaml:"lower_offer"`
			LowBalled     string `yaml:"low_balled"`
			Reminder      string `yaml:"reminder"`
		}{
			FAQ:           "",
			Initial:       "Hello @{{NAME}}!\n\nThanks for your interest in my item `{{ITEM}}`!",
			NotAvailable:  "Please note that this listing might not be available anymore as it was {{REASON}}.",
			Offered:       "Thank you for your offer of ${{OFFER}}!",
			PossibleOffer: "It looks like you are making an offer of ${{OFFER}}.",
			LowerOffer:    "WARNING: Offer was lowered!",
			LowBalled:     "WARNING: Your offer is {{PERCENT}}% below listing price, it's pretty low!",
			Reminder:      "REMINDER: We are dealing this in {{HOURS}} hour(s)!",
		},
		Reminders:     []int16{1, 4, 24},
		CommandPrefix: ".",
		StatePrune:    14,
		Forwarders:    nil,
	}
}

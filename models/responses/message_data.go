package responses

type MessageData struct {
	ImageProgressiveURL         string `json:"image_progressive_url"`
	ImageProgressiveLowRange    int    `json:"image_progressive_low_range"`
	SbSyncedAt                  int64  `json:"sb_synced_at"`
	ImageProgressiveMediumRange int    `json:"image_progressive_medium_range"`
	OfferID                     string `json:"offer_id"`
	OiID                        int64  `json:"oi_id"`
	Version                     string `json:"version"`
	Source                      string `json:"source"`
	CurrencySymbol              string `json:"currency_symbol"`
	OfferAmount                 string `json:"offer_amount"`
}

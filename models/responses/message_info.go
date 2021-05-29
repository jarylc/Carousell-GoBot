package responses

import "time"

// {"data":{"offers":[{"id":1328386920,"latest_price":"660.00","latest_price_formatted":"660","latest_price_message":"test","latest_price_created":"2021-05-24T15:19:10Z","state":"O","currency_symbol":"S$","chat_only":false,"has_accepted_offer":false,"is_product_sold":false,"offer_type":"received","vertical":null,"channel_url":"F3CB6187-CB42-4CD1-95FC-1C46F8856006-carousell-1328386920","order":null,"unread_count":0,"user":{"id":531373,"username":"shirleymew","is_suspended":false,"profile":{"image_url":"https://media.karousell.com/media/photos/profiles/2015/11/09/shirleymew_1447046671.jpg","image":"photos/profiles/2015/11/09/shirleymew_1447046671.jpg","verification_type":"VE","is_facebook_verified":true,"is_email_verified":true,"is_gplus_verified":false,"is_mobile_verified":true,"is_affiliate":false,"affiliate_name":null,"is_id_verified":false}},"product":{"id":1089654526,"title":"Dell XPS 13\" 9380 Lightweight Windows Laptop","primary_photo_url":"https://media.karousell.com/media/photos/products/2021/5/24/dell_xps_13_9380_lightweight_w_1621838215_1ea5f42d_thumbnail","price":"1100.00","price_formatted":"1,100","status":"L","vertical":null,"smart_attributes":{"brand_enum":"BRAND_DELL","caroupay":"true","city":"Singapore","condition_v2":"USED","covid_meetup_msg":"Stay safe at home and opt for delivery instead.\r\n{{{X}}}","deal_options":"MEETUP,CAROUPAY","description":"Price negotiable.\r\n\r\nOfficial Dell charger and premier laptop sleeve provided (as pictured)\r\n\r\nVery good productivity laptop for work. Portable and lightweight.\r\n\r\n8th Generation Intel® Core™ i7-8565U Processor (8MB Cache, up to 4.6 GHz, 4 cores)\r\nIntel® UHD Graphics 620\r\n13.3″ 4K Ultra HD (3840×2160) InfinityEdge Touch Display\r\n16GB LPDDR3 2133MHz\r\n512GB M.2 PCIe NVMe Solid State Drive\r\n2 Thunderbolt™ 3 with power delivery and DisplayPort (4 lanes of PCI Express Gen 3)\r\n1 USB-C 3.1 with power delivery and DisplayPort\r\n1 MicroSD card reader\r\n1 Headset jack\r\nHeight: 0.3″- 0.46″ (7.8mm – 11.6mm)\r\nWidth: 11.9″ (302mm) x Depth: 7.8″ (199mm)\r\nWeight: Starting at 2.7 lbs (1.23 kg)\r\nWaves MaxxAudio® Pro Stereo Speakers\r\nWidescreen HD (720p) 2.25mm webcam with 4 array digital microphones\r\nWindows 10 64bit English","fieldset_id":"5cb186ca-b878-4b20-8746-3bf39b1fa82a","is_free":"false","last_liked":"2021-05-24T14:13:52.483669Z","meetup_count":"3","memory_ram":"16","memory_ssd":"512","model_number":"13-9380","multi_quantities":"false","region":""},"marketplace":{"id":1880252,"name":"Singapore","country":{"id":1880251,"name":"Singapore","code":"SG","city_count":1},"location":{"latitude":1.28967,"longitude":103.85007}},"collection":{"cc_id":2195,"id":1793,"name":"Laptops"}},"dispute":null,"chat_status":null,"is_bot_offer":false,"make_offer_type":"normal"}]}}

type MessageInfo struct {
	ID                   int         `json:"id"`
	LatestPrice          string      `json:"latest_price"`
	LatestPriceFormatted string      `json:"latest_price_formatted"`
	LatestPriceMessage   string      `json:"latest_price_message"`
	LatestPriceCreated   time.Time   `json:"latest_price_created"`
	State                string      `json:"state"`
	CurrencySymbol       string      `json:"currency_symbol"`
	ChatOnly             bool        `json:"chat_only"`
	HasAcceptedOffer     bool        `json:"has_accepted_offer"`
	IsProductSold        bool        `json:"is_product_sold"`
	OfferType            string      `json:"offer_type"`
	Vertical             interface{} `json:"vertical"`
	ChannelURL           string      `json:"channel_url"`
	Order                interface{} `json:"order"`
	UnreadCount          int         `json:"unread_count"`
	User                 struct {
		ID          int    `json:"id"`
		Username    string `json:"username"`
		IsSuspended bool   `json:"is_suspended"`
		Profile     struct {
			ImageURL           string      `json:"image_url"`
			Image              string      `json:"image"`
			VerificationType   string      `json:"verification_type"`
			IsFacebookVerified bool        `json:"is_facebook_verified"`
			IsEmailVerified    bool        `json:"is_email_verified"`
			IsGplusVerified    bool        `json:"is_gplus_verified"`
			IsMobileVerified   bool        `json:"is_mobile_verified"`
			IsAffiliate        bool        `json:"is_affiliate"`
			AffiliateName      interface{} `json:"affiliate_name"`
			IsIDVerified       bool        `json:"is_id_verified"`
		} `json:"profile"`
	} `json:"user"`
	Product struct {
		ID              int         `json:"id"`
		Title           string      `json:"title"`
		PrimaryPhotoURL string      `json:"primary_photo_url"`
		Price           string      `json:"price"`
		PriceFormatted  string      `json:"price_formatted"`
		Status          string      `json:"status"`
		Vertical        interface{} `json:"vertical"`
		SmartAttributes struct {
			BrandEnum       string `json:"brand_enum"`
			Caroupay        string `json:"caroupay"`
			City            string `json:"city"`
			ConditionV2     string `json:"condition_v2"`
			CovidMeetupMsg  string `json:"covid_meetup_msg"`
			DealOptions     string `json:"deal_options"`
			Description     string `json:"description"`
			FieldsetID      string `json:"fieldset_id"`
			IsFree          string `json:"is_free"`
			LastLiked       string `json:"last_liked"`
			MeetupCount     string `json:"meetup_count"`
			MemoryRAM       string `json:"memory_ram"`
			MemorySsd       string `json:"memory_ssd"`
			ModelNumber     string `json:"model_number"`
			MultiQuantities string `json:"multi_quantities"`
			Region          string `json:"region"`
		} `json:"smart_attributes"`
		Marketplace struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Country struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				Code      string `json:"code"`
				CityCount int    `json:"city_count"`
			} `json:"country"`
			Location struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"location"`
		} `json:"marketplace"`
		Collection struct {
			CcID int    `json:"cc_id"`
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"collection"`
	} `json:"product"`
	Dispute       interface{} `json:"dispute"`
	ChatStatus    interface{} `json:"chat_status"`
	IsBotOffer    bool        `json:"is_bot_offer"`
	MakeOfferType string      `json:"make_offer_type"`
}

package responses

// {"custom_type":"MESSAGE","msg_id":7031000798,"is_op_msg":false,"request_id":"1621864528896","is_guest_msg":true,"message":"Hi","message_retention_hour":-1,"last_updated_at":1621864535167,"ts":1621864535167,"scrap_id":"","channel_url":"F3CB6187-CB42-4CD1-95FC-1C46F8856006-carousell-1328386920","mentioned_users":[],"translations":{},"is_removed":false,"user":{"guest_id":"344194","require_auth_for_profile_image":false,"name":"","is_active":true,"image":"https:\/\/static.sendbird.com\/sample\/profiles\/profile_35_512px.png","is_bot":false,"is_blocked_by_me":false,"friend_discovery_key":null,"role":"","friend_name":null,"id":263804,"metadata":{}},"data":"{\"is_dispute_acknowledgement\":false,\"offer_id\":\"1328386920\",\"refresh\":false,\"source\":\"ANDROID\",\"tags\":[]}","silent":false,"is_super":false,"mention_type":"users","channel_type":"group","channel_id":657371374,"sts":1621864535167}
// {"scrap_id":"","msg_id":7031028612,"is_super":false,"message_retention_hour":null,"last_updated_at":1621864864095,"mention_type":"users","is_op_msg":false,"mentioned_users":[],"translations":{},"ts":1621864864095,"data":"{\"image_progressive_url\": \"\", \"image_progressive_low_range\": 0, \"sb_synced_at\": 1621864864390, \"image_progressive_medium_range\": 0, \"offer_id\": \"1328386920\", \"oi_id\": 1328386920186486409, \"version\": \"1.0\", \"source\": \"DJANGO-API\", \"currency_symbol\": \"S$\", \"offer_amount\": \"660.0\"}","channel_type":"group","channel_id":657371374,"is_guest_msg":true,"user":{"guest_id":"531373","require_auth_for_profile_image":false,"name":"shirleymew","is_active":true,"image":"https:\/\/media.karousell.com\/media\/photos\/profiles\/2015\/11\/09\/shirleymew_1447046671.jpg","is_bot":false,"is_blocked_by_me":false,"friend_discovery_key":null,"role":"","friend_name":null,"id":1380681,"metadata":{}},"is_removed":false,"sts":1621864864095,"message":"Make Offer","channel_url":"F3CB6187-CB42-4CD1-95FC-1C46F8856006-carousell-1328386920","custom_type":"MAKE_OFFER","silent":false}
// {"scrap_id":"","msg_id":7049135652,"is_super":false,"message_retention_hour":null,"last_updated_at":1622192143043,"mention_type":"users","is_op_msg":false,"mentioned_users":[],"translations":{},"ts":1622192143043,"data":"{\"image_progressive_url\": \"\", \"image_progressive_low_range\": 0, \"sb_synced_at\": 1622192143319, \"image_progressive_medium_range\": 0, \"offer_id\": \"1328386920\", \"oi_id\": 1328386920219214304, \"version\": \"1.0\", \"source\": \"DJANGO-API\", \"currency_symbol\": \"S$\", \"offer_amount\": \"660.0\"}","channel_type":"group","channel_id":657371374,"is_guest_msg":true,"user":{"guest_id":"344194","require_auth_for_profile_image":false,"name":"","is_active":true,"image":"https:\/\/static.sendbird.com\/sample\/profiles\/profile_35_512px.png","is_bot":false,"is_blocked_by_me":false,"friend_discovery_key":null,"role":"","friend_name":null,"id":263804,"metadata":{}},"is_removed":false,"sts":1622192143043,"message":"Decline Offer","channel_url":"F3CB6187-CB42-4CD1-95FC-1C46F8856006-carousell-1328386920","custom_type":"DECLINE_OFFER","silent":false}

type Message struct {
	ScrapID              string        `json:"scrap_id"`
	MsgID                int64         `json:"msg_id"`
	IsSuper              bool          `json:"is_super"`
	MessageRetentionHour interface{}   `json:"message_retention_hour"`
	LastUpdatedAt        int64         `json:"last_updated_at"`
	MentionType          string        `json:"mention_type"`
	IsOpMsg              bool          `json:"is_op_msg"`
	MentionedUsers       []interface{} `json:"mentioned_users"`
	Translations         struct {
	} `json:"translations"`
	TS          int64  `json:"ts"`
	Data        string `json:"data"`
	ChannelType string `json:"channel_type"`
	ChannelID   int    `json:"channel_id"`
	IsGuestMsg  bool   `json:"is_guest_msg"`
	User        struct {
		GuestID                    string      `json:"guest_id"`
		RequireAuthForProfileImage bool        `json:"require_auth_for_profile_image"`
		Name                       string      `json:"name"`
		IsActive                   bool        `json:"is_active"`
		Image                      string      `json:"image"`
		IsBot                      bool        `json:"is_bot"`
		IsBlockedByMe              bool        `json:"is_blocked_by_me"`
		FriendDiscoveryKey         interface{} `json:"friend_discovery_key"`
		Role                       string      `json:"role"`
		FriendName                 interface{} `json:"friend_name"`
		ID                         int         `json:"id"`
		Metadata                   struct {
		} `json:"metadata"`
	} `json:"user"`
	IsRemoved  bool   `json:"is_removed"`
	Sts        int64  `json:"sts"`
	Message    string `json:"message"`
	ChannelURL string `json:"channel_url"`
	CustomType string `json:"custom_type"`
	Silent     bool   `json:"silent"`
}

package requests

// {"custom_type":"MESSAGE","msg_id":7031000798,"is_op_msg":false,"request_id":"1621864528896","is_guest_msg":true,"message":"Hi","message_retention_hour":-1,"last_updated_at":1621864535167,"ts":1621864535167,"scrap_id":"","channel_url":"F3CB6187-CB42-4CD1-95FC-1C46F8856006-carousell-1328386920","mentioned_users":[],"translations":{},"is_removed":false,"user":{"guest_id":"344194","require_auth_for_profile_image":false,"name":"","is_active":true,"image":"https:\/\/static.sendbird.com\/sample\/profiles\/profile_35_512px.png","is_bot":false,"is_blocked_by_me":false,"friend_discovery_key":null,"role":"","friend_name":null,"id":263804,"metadata":{}},"data":"{\"is_dispute_acknowledgement\":false,\"offer_id\":\"1328386920\",\"refresh\":false,\"source\":\"ANDROID\",\"tags\":[]}","silent":false,"is_super":false,"mention_type":"users","channel_type":"group","channel_id":657371374,"sts":1621864535167}
// {"scrap_id":"","msg_id":7031028612,"is_super":false,"message_retention_hour":null,"last_updated_at":1621864864095,"mention_type":"users","is_op_msg":false,"mentioned_users":[],"translations":{},"ts":1621864864095,"data":"{\"image_progressive_url\": \"\", \"image_progressive_low_range\": 0, \"sb_synced_at\": 1621864864390, \"image_progressive_medium_range\": 0, \"offer_id\": \"1328386920\", \"oi_id\": 1328386920186486409, \"version\": \"1.0\", \"source\": \"DJANGO-API\", \"currency_symbol\": \"S$\", \"offer_amount\": \"660.0\"}","channel_type":"group","channel_id":657371374,"is_guest_msg":true,"user":{"guest_id":"531373","require_auth_for_profile_image":false,"name":"shirleymew","is_active":true,"image":"https:\/\/media.karousell.com\/media\/photos\/profiles\/2015\/11\/09\/shirleymew_1447046671.jpg","is_bot":false,"is_blocked_by_me":false,"friend_discovery_key":null,"role":"","friend_name":null,"id":1380681,"metadata":{}},"is_removed":false,"sts":1621864864095,"message":"Make Offer","channel_url":"F3CB6187-CB42-4CD1-95FC-1C46F8856006-carousell-1328386920","custom_type":"MAKE_OFFER","silent":false}

type Reply struct {
	ChannelURL       string   `json:"channel_url"`
	Message          string   `json:"message"`
	Data             string   `json:"data"`
	MentionType      string   `json:"mention_type"`
	MentionedUserIds []string `json:"mentioned_user_ids"`
	CustomType       string   `json:"custom_type"`
	ReqID            int64    `json:"req_id"`
}

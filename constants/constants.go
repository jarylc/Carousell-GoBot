package constants

const CHANNEL = "f3cb6187-cb42-4cd1-95fc-1c46f8856006"
const QUERY = "p=JS&pv=Mozilla%2F5.0%20(X11%3B%20Linux%20x86_64%3B%20rv%3A90.0)%20Gecko%2F20100101%20Firefox%2F90.0&sv=3.1.32&ai={{CHANNEL}}&user_id={{USERID}}&access_token={{TOKEN}}&active=1&SB-User-Agent=JS%2Fc3.1.32%2F%2F&Request-Sent-Timestamp={{TIME}}&include_extra_data%3Dpremium_feature_list%2Cfile_upload_size_limit%2Capplication_attributes%2Cemoji_hash%26use_local_cache%3D1"

const PRICE_EXPRESSION = "(^\\d{1,5}(\\.{1}\\d{1,2})?$)|((?<=(\\$)|((offer)|(quote)|(can)|(please)|(pls)|(do)|(deal)|(sell)|(at))\\s+)(\\d{1,5}(\\.{1}\\d{1,2})?)(?!((\\d)|(\\s*pm)|(\\s*am)|(\\s*at)|(\\s*in)|(\\s*evening)|(\\s*night)|(\\s*tdy)|(\\s*today)|(\\s*tonight)|(\\s*tmr)|(\\s*tomorrow)|(\\s*ltr)|(\\s*later)))|((?<!(\\d))(\\d{1,5}(\\.{1}\\d{1,2})?)(?=((\\$)|(\\?))|\\s+((offer)|(quote)|(can)|(please)|(pls)|(deal)|(fast)|(quick)|(bucks)|(dollar)|(ok)))))"
const READABLE_DATE_FORMAT = "Monday, 02 January 2006, 03:04:05PM"

const MESSAGE = "MESSAGE"
const MAKE_OFFER = "MAKE_OFFER"

const NEW_CHAT = "NEW CHAT"
const OFFICIAL = "OFFICIAL"
const LOW_BALL = "LOW BALL"
const SUPER_LOW_BALL = "SUPER LOW BALL"
const LOWERED = "LOWERED"

const CAROUSELL_CHANNEL = "{{CHANNEL}}-carousell-{{CHATID}}"
const CAROUSELL_PING = "PING{\"id\":{{TIME}},\"active\":1,\"req_id\":\"\"}\n"
const CAROUSELL_MESG = "MESG{DATA}\n"

const CAROUSELL_URL_TOKEN = "https://www.carousell.sg/ds/api/1.0/chat/token/?_path=/1.0/chat/token/" //nolint:gosec
const CAROUSELL_URL_CHAT = "ws-" + CHANNEL + ".sendbird.com/"
const CAROUSELL_URL_MESSAGES = "https://www.carousell.sg/api-service/offer/1.0/me/?count=20&type={{TYPE}}"

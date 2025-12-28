package consts

const (
	ApiV1URL = "https://api.coze.ai/v1"
	ApiV3URL = "https://api.coze.ai/v3"

	CreateConversationURL   = ApiV1URL + "/conversation/create"
	RetrieveConversationURL = ApiV3URL + "/chat/retrieve"
	MessageListURL          = ApiV3URL + "/chat/message/list"
)

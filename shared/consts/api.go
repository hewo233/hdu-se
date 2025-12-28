package consts

const (
	ApiV1URL = "https://api.coze.ai/v1"
	ApiV3URL = "https://api.coze.ai/v3"

	BotID = "7563218003241058343"

	CreateConversationURL   = ApiV1URL + "/conversation/create"
	CreateChatURL           = ApiV3URL + "/chat"
	RetrieveConversationURL = ApiV3URL + "/chat/retrieve"
	ChatMessageListURL      = ApiV3URL + "/chat/message/list"

	ConversationMessageListURL = ApiV1URL + "/conversation/message/list"
)

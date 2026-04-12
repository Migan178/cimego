package cimego

const (
	APIVersion = "v1"
	APIBaseURL = "https://ci.me/api/openapi"
	APIOpen    = APIBaseURL + "/open/" + APIVersion

	EndpointAuth  = APIBaseURL + "/auth/" + APIVersion
	EndpointToken = EndpointAuth + "/token"
	EndpointMe    = APIOpen + "/users/me"

	EndpointChannels           = APIOpen + "/channels"
	EndpointChannelFollowers   = EndpointChannels + "/followers"
	EndpointChannelSubscribers = EndpointChannels + "/subscribers"
	EndpointChannelManagers    = EndpointChannels + "/streaming-roles"
	EndpointRestrictChannels   = APIOpen + "/restrict-channels"

	EndpointLives        = APIOpen + "/lives"
	EndpointLivesSetting = EndpointLives + "/setting"

	EndpointStreams   = APIOpen + "/streams"
	EndpointStreamKey = EndpointStreams + "/key"

	EndpointChats        = APIOpen + "/chats"
	EndpointChatSettings = EndpointChats + "/settings"
	EndpointChatSend     = EndpointChats + "/send"
	EndpointChatNotice   = EndpointChats + "/notice"

	EndpointCategories     = APIOpen + "/categories"
	EndpointCategorySearch = EndpointCategories + "/search"
)

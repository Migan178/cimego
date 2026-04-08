package cimego

var (
	Version = "v1"
	BaseURL = "https://ci.me/api/openapi/" + Version

	EndpointAuthorization = BaseURL + "/token"
	EndpointMe            = BaseURL + "/users/me"

	EndpointChannels           = BaseURL + "/channels"
	EndpointChannelFollowers   = EndpointChannels + "/followers"
	EndPointChannelSubscribers = EndpointChannels + "/subscribers"
)

package cimego

const (
	APIVersion = "v1"
	APIBaseURL = "https://ci.me/api/openapi/" + APIVersion

	EndpointAuthorization = APIBaseURL + "/token"
	EndpointMe            = APIBaseURL + "/users/me"

	EndpointChannels           = APIBaseURL + "/channels"
	EndpointChannelFollowers   = EndpointChannels + "/followers"
	EndPointChannelSubscribers = EndpointChannels + "/subscribers"
)

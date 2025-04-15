package pb

const (
	FlowRequest int32 = iota + 1
	FlowResponse
)

// Type
const (
	TypeHeartbeat int32 = iota + 1
	TypeCommand
	TypeMessage
	TypeEvent
)

// NeedAck
const (
	NO int32 = iota
	YES
)

const (
	CommandTypeUserLogin      string = "USER_LOGIN"
	CommandTypeUserLogout            = "USER_LOGOUT"
	CommandTypeFriendAdd             = "FRIEND_ADD"
	CommandTypeFriendAddAgree        = "FRIEND_ADD_AGREE"
	CommandTypeFriendReject          = "FRIEND_ADD_REJECT"
)

const (
	MessageTypeText  string = "TEXT"
	MessageTypeImage string = "IMAGE"
	MessageTypeAudio string = "AUDIO"
	MessageTypeVideo string = "VIDEO"
)

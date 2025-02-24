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
	CTypeUserLogin      string = "USER_LOGIN"
	CTypeUserLogout            = "USER_LOGOUT"
	CTypeFriendAdd             = "FRIEND_ADD"
	CTypeFriendAddAgree        = "FRIEND_ADD_AGREE"
	CTypeFriendReject          = "FRIEND_ADD_REJECT"
)

const (
	CTypeText  string = "TEXT"
	CTypeImage string = "IMAGE"
	CTypeAudio string = "AUDIO"
	CTypeVideo string = "VIDEO"
)

const (
	TTypeSingle int32 = iota + 1
	TTypeGroup
)

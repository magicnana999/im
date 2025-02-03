package pb

const (
	FlowRequest int32 = iota + 1
	FlowResponse
)

// Type
const (
	BTypeHeartbeat int32 = iota + 1
	BTypeCommand
	BTypeMessage
	BTypeNotice
	BTypeTips
)

// NeedAck
const (
	NO int32 = iota
	YES
)

const (
	CTypeUserLogin      string = "USER_LOGIN"
	CTypeUserLogout     string = "USER_LOGOUT"
	CTypeFriendAdd      string = "FRIEND_ADD"
	CTypeFriendRemove   string = "FRIEND_REMOVE"
	CTypeMessageHistory string = "MESSAGE_HISTORY"
)

const (
	CTypeText  string = "TEXT"
	CTypeImage string = "IMAGE"
	CTypeAudio string = "AUDIO"
	CTypeVideo string = "VIDEO"
)

const (
	TargetTypeSingle int32 = iota + 1
	TargetTypeGroup
)

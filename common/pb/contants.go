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
	CTypeUserLogin string = "USER_LOGIN"
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

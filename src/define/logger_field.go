package define

// Zap field keys for structured logging.
const (
	OP        = "op"
	Conn      = "conn"
	MessageId = "messageId"
	Status    = "status"
	Url       = "url"
	Arg       = "arg"
	Ret       = "ret"
	Req       = "req"
	Res       = "res"
)

const (
	OpInit    = "init"
	OpStart   = "start"
	OpTicking = "ticking"
	OpWrite   = "write"
	OpResave  = "resave"
	OpSubmit  = "submit"
	OpStop    = "stop"
	OpClose   = "close"
	OpAccept  = "accept"
)

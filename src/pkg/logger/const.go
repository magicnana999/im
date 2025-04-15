package logger

// define key of zap.field
const (
	Op  = "op"
	Url = "url"
	Arg = "arg"
	Ret = "ret"
	Req = "req"
	Res = "res"
)

// define value of zap.field，and there are match to Op
const (
	OpInit    = "init"
	OpStart   = "start"
	OpAdvance = "advance"
	OpSubmit  = "submit"
	OpSlowSQL = "slow sql"
	OpClose   = "close"
	OpSend    = "send"
	OpReceive = "receive"
	OpQuery   = "query"
)

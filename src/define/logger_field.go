package define

// Zap field keys for structured logging.
const (
	OP        = "op"        // OP is the key for the operation field.
	Conn      = "conn"      // Conn is the key for the connection description field.
	MessageId = "messageId" // MessageId is the key for the message ID field.
	Status    = "status"    // Status is the key for the operation status field.
	Url       = "url"       // Url is the key for the URL field.
	Arg       = "arg"       // Arg is the key for the argument field.
	Ret       = "ret"       // Ret is the key for the return value field.
	Req       = "req"       // Req is the key for the request field.
	Res       = "res"       // Res is the key for the response field.
)

// Operation values for the OP field.
const (
	OpInit    = "init"     // OpInit indicates an initialization operation.
	OpStart   = "start"    // OpStart indicates a start operation.
	OpAdvance = "advance"  // OpAdvance indicates an advance operation.
	OpSubmit  = "submit"   // OpSubmit indicates a submit operation.
	OpStop    = "stop"     // OpStop indicates a stop operation.
	OpSlowSQL = "slow sql" // OpSlowSQL indicates a slow SQL query operation.
	OpClose   = "close"    // OpClose indicates a close operation.
	OpSend    = "send"     // OpSend indicates a send operation.
	OpReceive = "receive"  // OpReceive indicates a receive operation.
	OpQuery   = "query"    // OpQuery indicates a query operation.
)

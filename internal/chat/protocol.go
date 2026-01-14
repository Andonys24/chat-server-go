package chat

const (
	CmdEnter        = "ENTER"
	CmdMessage      = "MESSAGE"
	CmdAll          = "ALL"
	CmdUsers        = "USERS"
	CmdExit         = "EXIT"
	CmdCLeanConsole = "CLEAN CONSOLE"
)

const (
	RespOkEnter    = "OK ENTER"
	RespErrorEnter = "ERROR ENTER"
	RespMsgFrom    = "OF"
	RespInfo       = "INFO"
	RespOkClean    = "OK CLEAN"
)

const (
	InfoTypeError   = "ERROR"   // Un error general
	InfoTypeSuccess = "SUCCESS" // Una acción completada
)

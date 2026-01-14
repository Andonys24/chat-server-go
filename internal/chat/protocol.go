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
	InfoTypeEnter   = "ENTER"   // Alguien se unió
	InfoTypeExit    = "EXIT"    // Alguien se fue
	InfoTypeError   = "ERROR"   // Un error general
	InfoTypeSuccess = "SUCCESS" // Una acción completada
)

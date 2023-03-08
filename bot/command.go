package bot

type Command func(Context) error

type CommandStruct struct {
	command Command
	help    string
}

type CmdMap map[string]CommandStruct

type CommandHandler struct {
	cmds CmdMap
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{make(CmdMap)}
}

func (cmdHandler CommandHandler) Get(cmdName string) (*Command, bool) {
	cmd, found := cmdHandler.cmds[cmdName]
	return &cmd.command, found
}

func (cmdHandler CommandHandler) Register(cmdName string, cmd Command, helpMsg string) {
	cmdStruct := CommandStruct{cmd, helpMsg}
	cmdHandler.cmds[cmdName] = cmdStruct
}

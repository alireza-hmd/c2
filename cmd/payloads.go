package cmd

func Payload(s *Services) {
	PayloadMenu.Clear()

	for {
		c, args := PayloadMenu.Input()
		if c == "" {
			continue
		}
		command, exists := GetCommand(c, PayloadMenu.Commands)
		if !exists {
			ErrorResponse(PayloadMenu.Name, "invalid command\n")
			continue
		}
		PayloadCommand(command, args, s)
	}
}

func PayloadCommand(c *Command, args []string, s *Services) {
	switch c.Name {
	case "list":
	case "generate":
	case "home":
		C2(s)
	case "help":
		PayloadMenu.ShowHelp()
	case "exit":
		Exit()
	}
}

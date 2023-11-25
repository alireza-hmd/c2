package cmd

func C2(s *Services) {
	C2Menu.Clear()

	for {
		c, args := C2Menu.Input()
		if c == "" {
			continue
		}
		command, exists := GetCommand(c, C2Menu.Commands)
		if !exists {
			ErrorResponse(C2Menu.Name, "invalid command\n")
			continue
		}
		C2Command(command, args, s)
	}
}

func C2Command(c *Command, args []string, s *Services) {
	switch c.Name {
	case "listeners":
		Listener(s)
	case "clients":
		Client(s)
	case "payloads":
		Payload(s)
	case "home":
		C2(s)
	case "help":
		C2Menu.ShowHelp()
	case "exit":
		Exit()
	}
}

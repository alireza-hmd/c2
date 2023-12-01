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
		Listeners(s)
	case "clients":
		Clients(s)
	case "tasks":
		Tasks(s)
	case "home":
		C2(s)
	case "help":
		C2Menu.ShowHelp()
	case "exit":
		Exit()
	}
}

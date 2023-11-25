package cmd

func Client(s *Services) {
	ClientMenu.Clear()

	for {
		c, args := ClientMenu.Input()
		if c == "" {
			continue
		}
		command, exists := GetCommand(c, ClientMenu.Commands)
		if !exists {
			ErrorResponse(ClientMenu.Name, "invalid command\n")
			continue
		}
		ClientCommand(command, args, s)
	}
}

func ClientCommand(c *Command, args []string, s *Services) {
	switch c.Name {
	case "list":
	case "remove":
	case "home":
		C2(s)
	case "help":
		ClientMenu.ShowHelp()
	case "exit":
		Exit()
	}
}

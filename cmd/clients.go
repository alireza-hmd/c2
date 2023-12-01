package cmd

import "fmt"

func Clients(s *Services) {
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
		ClientsCommand(command, args, s)
	}
}

func ClientsCommand(c *Command, args []string, s *Services) {
	switch c.Name {
	case "list":
		cc, err := s.Clients.List()
		if err != nil {
			ErrorResponse(ClientMenu.Name, err.Error())
			return
		}
		if len(cc) == 0 {
			ErrorResponse(ClientMenu.Name, "No clients connected yet")
			return
		}
		for _, c := range cc {
			fmt.Printf("client #%d with \"%s\" token is connected to \"%s\" listener\n", c.ID, c.Token, c.Listener)
		}
	case "remove":
		if len(args) != len(c.Args) {
			ErrorResponse(ClientMenu.Name, "invalid arugment. visit the help menu.")
			return
		}
		if err := s.Clients.Delete(args[0]); err != nil {
			ErrorResponse(ClientMenu.Name, err.Error())
			return
		}
		fmt.Println("removed successfuly")
	case "home":
		C2(s)
	case "help":
		ClientMenu.ShowHelp()
	case "exit":
		Exit()
	}
}

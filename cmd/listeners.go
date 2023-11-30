package cmd

import (
	"fmt"
	"strconv"

	"github.com/alireza-hmd/c2/listener"
)

func Listener(s *Services) {
	ListenerMenu.Clear()
	for {
		c, args := ListenerMenu.Input()
		if c == "" {
			continue
		}
		command, exists := GetCommand(c, ListenerMenu.Commands)
		if !exists {
			ErrorResponse(ListenerMenu.Name, "invalid command\n")
			continue
		}
		ListenerCommand(command, args, s)
	}
}

func ListenerCommand(c *Command, args []string, s *Services) {
	switch c.Name {
	case "list":
		ll, err := s.Listener.List()
		if err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
		if len(ll) == 0 {
			ErrorResponse(ListenerMenu.Name, "No active listeners")
			return
		}
		for _, l := range ll {
			fmt.Printf("%s listener is active on port %d\n", l.Name, l.Port)
		}
	case "start":
		if len(args) != len(c.Args) {
			ErrorResponse(ListenerMenu.Name, "invalid arugment. visit the help menu.")
			return
		}
		p, _ := strconv.Atoi(args[2])
		l := listener.NewListener(args[0], args[1], p)
		_, err := s.Listener.Create(l, s.Stop)
		if err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
	case "stop":
		if len(args) != len(c.Args) {
			ErrorResponse(ListenerMenu.Name, "invalid arugment. visit the help menu.")
			return
		}
		l, err := s.Listener.Get(args[0])
		if err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
		s.Stop[l.ID] = make(chan listener.Cancel, 1)
		if err := s.Listener.Activation(l, s.Stop[l.ID], listener.DeactiveStatus); err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
	case "remove":
		if len(args) != len(c.Args) {
			ErrorResponse(ListenerMenu.Name, "invalid arugment. visit the help menu.")
			return
		}
		l, err := s.Listener.Get(args[0])
		if err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
		s.Stop[l.ID] = make(chan listener.Cancel, 1)
		if err := s.Listener.Delete(args[0], s.Stop[l.ID]); err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
	case "home":
		C2(s)
	case "help":
		ListenerMenu.ShowHelp()
	case "exit":
		Exit()
	}
}

package cmd

import (
	"fmt"
	"strconv"

	"github.com/alireza-hmd/c2/listeners"
)

func Listeners(s *Services) {
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
		ListenersCommand(command, args, s)
	}
}

func ListenersCommand(c *Command, args []string, s *Services) {
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
	case "add":
		if len(args) != len(c.Args) {
			ErrorResponse(ListenerMenu.Name, "invalid arugment. visit the help menu.")
			return
		}
		p, _ := strconv.Atoi(args[2])
		l := listeners.NewListener(args[0], args[1], p)
		_, err := s.Listener.Create(l, s.Stop)
		if err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
	case "start":
		if len(args) != len(c.Args) {
			ErrorResponse(ListenerMenu.Name, "invalid arugment. visit the help menu.")
			return
		}
		l, err := s.Listener.Get(args[0])
		if err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
		s.Stop[l.ID] = make(chan listeners.Cancel, 1)
		if err := s.Listener.Activation(l, s.Stop[l.ID], listeners.ActiveStatus); err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
		fmt.Printf("started successfuly on port %d\n", l.Port)
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
		s.Stop[l.ID] = make(chan listeners.Cancel, 1)
		if err := s.Listener.Activation(l, s.Stop[l.ID], listeners.DeactiveStatus); err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
		fmt.Println("stopped successfuly")
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
		if err := s.Clients.DeleteListenerClient(l.Name); err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
		s.Stop[l.ID] = make(chan listeners.Cancel, 1)
		if err := s.Listener.Delete(args[0], s.Stop[l.ID]); err != nil {
			ErrorResponse(ListenerMenu.Name, err.Error())
			return
		}
		fmt.Println("removed successfuly")
	case "home":
		C2(s)
	case "help":
		ListenerMenu.ShowHelp()
	case "exit":
		Exit()
	}
}

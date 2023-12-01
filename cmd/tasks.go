package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alireza-hmd/c2/tasks"
)

var ValidCommands = []string{"timeout", "pwd"}

func Tasks(s *Services) {
	TasksMenu.Clear()

	for {
		c, args := TasksMenu.Input()
		if c == "" {
			continue
		}
		command, exists := GetCommand(c, TasksMenu.Commands)
		if !exists {
			ErrorResponse(TasksMenu.Name, "invalid command\n")
			continue
		}
		TasksCommand(command, args, s)
	}
}

func TasksCommand(c *Command, args []string, s *Services) {
	switch c.Name {
	case "list":
		if len(args) != len(c.Args) {
			ErrorResponse(TasksMenu.Name, "invalid arugment. visit the help menu.")
			return
		}
		tt, err := s.Tasks.ListToDoTasks(args[0])
		if err != nil {
			ErrorResponse(TasksMenu.Name, err.Error())
			return
		}
		if len(tt) == 0 {
			ErrorResponse(TasksMenu.Name, "No tasks found for client")
			return
		}
		for _, t := range tt {
			fmt.Printf("task #%d created on \"%s\" listener with \"%s\" command\n", t.ID, t.Listener, t.Command)
		}
	case "add":
		if len(args) < 2 {
			ErrorResponse(TasksMenu.Name, "invalid arugment. visit the help menu.")
			return
		}
		c, err := s.Clients.Get(args[0])
		if err != nil {
			ErrorResponse(TasksMenu.Name, err.Error())
			return
		}
		cmd, cmdType := ParseTask(args[1:])
		if cmd == "" {
			return
		}
		id, err := s.Tasks.Create(c.Token, c.Listener, cmd, cmdType)
		if err != nil {
			ErrorResponse(TasksMenu.Name, err.Error())
			return
		}
		fmt.Printf("tasks #%d created\n", id)
	case "results":
		if len(args) != len(c.Args) {
			ErrorResponse(TasksMenu.Name, "invalid arugment. visit the help menu.")
			return
		}
		tt, err := s.Tasks.ListDoneTasks(args[0])
		if err != nil {
			ErrorResponse(TasksMenu.Name, err.Error())
			return
		}
		if len(tt) == 0 {
			ErrorResponse(TasksMenu.Name, "No tasks found for client")
			return
		}
		for _, t := range tt {
			if t.Status == tasks.Failed {
				fmt.Printf("task #%d \"%s\" failed to run on client %s with result: %s\n", t.ID, t.Command, t.Client, t.Result)
				continue
			}
			fmt.Printf("task #%d \"%s\" done on client %s with result: %s\n", t.ID, t.Command, t.Client, t.Result)
		}
	case "remove":
		if len(args) != len(c.Args) {
			ErrorResponse(TasksMenu.Name, "invalid arugment. visit the help menu.")
			return
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			ErrorResponse(TasksMenu.Name, "invalid arugment. id is a number")
			return
		}
		if err := s.Tasks.Delete(id); err != nil {
			ErrorResponse(TasksMenu.Name, err.Error())
			return
		}
		fmt.Println("removed successfuly")
	case "home":
		C2(s)
	case "help":
		TasksMenu.ShowHelp()
	case "exit":
		Exit()
	}
}

func ParseTask(cmd []string) (string, string) {
	for _, v := range ValidCommands {
		if v == cmd[0] {
			return strings.Join(cmd, " "), cmd[0]
		}
	}
	ErrorResponse(TasksMenu.Name, "undefined command")
	return "", ""
}

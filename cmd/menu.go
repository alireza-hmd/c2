package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alireza-hmd/c2/clients"
	"github.com/alireza-hmd/c2/listeners"
	"github.com/alireza-hmd/c2/tasks"
)

type Services struct {
	Tasks    tasks.UseCase
	Clients  clients.UseCase
	Listener listeners.UseCase
	Stop     map[int](chan listeners.Cancel)
}

type Command struct {
	Name        string
	Description string
	Args        []string
}
type Menu struct {
	Name     string
	Commands []*Command
}

func NewMenu(name string) *Menu {
	var commands []*Command
	return &Menu{
		Name:     name,
		Commands: commands,
	}
}

var C2Menu *Menu
var ListenerMenu *Menu
var ClientMenu *Menu
var TasksMenu *Menu

func (m *Menu) Input() (string, []string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s) ", m.Name)
	i, _ := reader.ReadString('\n')
	i = i[:len(i)-1]
	input := strings.Split(i, " ")
	if len(input) < 1 {
		return "", []string{}
	}
	if len(input) == 1 {
		return input[0], []string{}
	}
	return input[0], input[1:]
}

func (m *Menu) AddCommand(name, description string, args []string) {
	cmd := &Command{Name: name, Description: description, Args: args}
	m.Commands = append(m.Commands, cmd)
}

func (m *Menu) AddDefaultCommands() {
	home := &Command{Name: "home", Description: "Back to home.", Args: []string{}}
	help := &Command{Name: "help", Description: "Show help.", Args: []string{}}
	exit := &Command{Name: "exit", Description: "Exit.", Args: []string{}}
	var commands []*Command
	commands = append(commands, home, help, exit)
	m.Commands = append(m.Commands, commands...)
}

func (m *Menu) ShowHelp() {
	fmt.Printf("[+] Available Commands:\n\n")
	fmt.Println("\t Command                         Description                         Arguments")
	fmt.Println("\t---------                       -------------                       -----------")

	for _, c := range m.Commands {
		fmt.Printf("\t%s%s%s%s%s%s\n", c.Name, spaceGen(32-len(c.Name)), c.Description, spaceGen(32-len(c.Description)), argsToStr(c.Args), spaceGen(32-len(argsToStr(c.Args))))
	}
	fmt.Print("\n")
}

func (m *Menu) Clear() {
	fmt.Print("\033[H\033[2J")
}

func Init(s *Services) {
	C2Menu = NewMenu("c2")
	ListenerMenu = NewMenu("listeners")
	ClientMenu = NewMenu("clients")
	TasksMenu = NewMenu("tasks")

	C2Menu.AddCommand("listeners", "Manage listeners.", []string{})
	C2Menu.AddCommand("clients", "Manage active clients.", []string{})
	C2Menu.AddCommand("tasks", "Manage tasks.", []string{})
	C2Menu.AddDefaultCommands()

	ListenerMenu.AddCommand("list", "List active listeners.", []string{})
	ListenerMenu.AddCommand("start", "Start a listener.", []string{"<name>", "<connection>", "<port>"})
	ListenerMenu.AddCommand("stop", "Stop an active listener.", []string{"<name>"})
	ListenerMenu.AddCommand("remove", "Remove a listener.", []string{"<name>"})
	ListenerMenu.AddDefaultCommands()

	ClientMenu.AddCommand("list", "List active clients.", []string{})
	ClientMenu.AddCommand("remove", "Remove a client.", []string{"<name>"})
	ClientMenu.AddDefaultCommands()

	TasksMenu.AddCommand("list", "List client tasks.", []string{"<client_token>"})
	TasksMenu.AddCommand("add", "Generate a task.", []string{"<client_token>", "<command>"})
	TasksMenu.AddCommand("remove", "Remove a task.", []string{"<id>"})
	TasksMenu.AddDefaultCommands()

	C2(s)
}

func GetCommand(cmd string, cc []*Command) (*Command, bool) {
	for _, c := range cc {
		if cmd == c.Name {
			return c, true
		}
	}
	return nil, false
}

func Exit() {
	fmt.Print("goodbye")
	os.Exit(0)
}

func spaceGen(len int) string {
	var str string
	if len < 1 {
		str += " "
		return str
	}
	for i := 0; i < len; i++ {
		str += " "
	}
	return str
}

func argsToStr(args []string) string {
	var str string
	for i, a := range args {
		str += a
		if i != len(args)-1 {
			str += " "
		}
	}
	return str
}

func ErrorResponse(name, msg string) {
	fmt.Printf("%s) %s\n", name, msg)
}

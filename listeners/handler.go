package listeners

import (
	"fmt"
	"io"
	"log"

	"github.com/alireza-hmd/c2/clients"
	"github.com/alireza-hmd/c2/pkg/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitHandler(ls UseCase, cs clients.UseCase, name string, port int) {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowCredentials: true,
	}))
	api := app.Group("/")
	Handle(api, ls, cs)
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

	lUpdate := &ListenerUptade{
		Connected: Connected,
	}
	if err := h.ls.Update(data.Listener, lUpdate); err != nil {
		response.ErrorResponse(c, 400, err.Error())
		return
	}
	cli := clients.NewClient(data.Listener, data.RemoteIP, data.ClientType)
	_, err = h.cs.Create(cli)
	if err != nil {
		response.ErrorResponse(c, 400, err.Error())
		return
	}
	body := struct {
		Token string
	}{
		Token: cli.Token,
	}
	response.OkResponseWithData(c, "client registered successfully", body)

}

func Handle(r *gin.RouterGroup, ls UseCase, cs clients.UseCase) {
	h := &handler{
		ls: ls,
		cs: cs,
	}
	api := r.Group("/clients")
	{
		api.POST("/register", h.Register)
	}
}

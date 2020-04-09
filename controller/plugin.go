package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.com/xiayesuifeng/goblog/core"
	"gitlab.com/xiayesuifeng/goblog/plugin"
	"log"
	"os"
)

type Plugin struct {
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (Plugin) Gets(ctx *gin.Context) {
	ctx.JSON(200, core.SuccessDataResult("plugins", plugin.GetPlugins()))
}

func (Plugin) Download(ctx *gin.Context) {
	name := ctx.Param("name")

	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	defer conn.Close()

	exit := make(chan bool, 1)

	go func() {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				exit <- true
			}
		}
	}()

	status := make(chan string)

	go func() {
		for {
			tmp, ok := <-status
			if !ok {
				exit <- true
				break
			}
			if err := conn.WriteJSON(gin.H{
				"status": tmp,
			}); err != nil {
				log.Println("a", err)
				exit <- true
			}
		}
	}()

	go func() {
		err := plugin.DownloadPlugin(name, status)
		if err != nil {
			status <- err.Error()
		}

		close(status)
	}()

	<-exit

	os.Remove("plugins/" + name + ".so.temp")
	core.UpgraderGoBlog()
}

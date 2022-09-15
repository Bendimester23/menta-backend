package routes

import (
	"menta-backend/controller"

	"github.com/gofiber/websocket/v2"
)

var wsController = &controller.WsController{
	Connections: make([]*controller.WsConnection, 0),
}

func HandleWs(c *websocket.Conn) {
	wsController.OnConnect(c)
}

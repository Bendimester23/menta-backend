package controller

import (
	"encoding/json"
	"log"
	"menta-backend/db"

	"github.com/gofiber/websocket/v2"
	"github.com/lucsky/cuid"
)

type WsController struct {
	Connections []*WsConnection
}

type WsConnection struct {
	Raw          *websocket.Conn
	IsIdentified bool
	UID          string
	SID          string
	CurrentRoom  string
	Controller   *WsController
	Connected    bool
}

func (c *WsConnection) Init() {
	c.Send("sid", c.SID)
}

func (c *WsController) Disconnect(sid string) {
	for k, v := range c.Connections {
		if v.SID == sid {
			v.Connected = false
			c.Connections = append(c.Connections[:k], c.Connections[k+1:]...)
			return
		}
	}
}

func (c *WsConnection) Disconnect() {
	c.Controller.Disconnect(c.SID)
}

func (c *WsConnection) Send(name string, data interface{}) {
	if c.Raw.WriteJSON(map[string]interface{}{
		"event": name,
		"data":  data,
	}) != nil {
		c.Controller.Disconnect(c.SID)
	}
}

func (c *WsConnection) HandleMessage(name string, data interface{}) {
	if data == nil {
		return
	}
	if name == "setroom" {
		c.CurrentRoom = data.(string)
		return
	}
	if name == "setuid" {
		c.UID = data.(string)
		return
	}
	if name == "msgin" {
		res, err := db.DB.ChatMessage.CreateOne(
			db.ChatMessage.Author.Link(
				db.User.ID.Equals(c.UID),
			),
			db.ChatMessage.Content.Set(data.(string)),
			db.ChatMessage.Room.Link(
				db.ChatRoom.ID.Equals(c.CurrentRoom),
			),
		).With(
			db.ChatMessage.Author.Fetch(),
		).Exec(ctx)

		if err != nil {
			log.Println(err.Error())
			return
		}

		res.Author().Password = "titok"
		res.Author().Email = "lábamközött@gmail.com"
		res.Author().EmailCode = "69420"

		c.Controller.Broadcast("msg", res)
	}
}

func (w *WsController) Broadcast(name string, data interface{}) {
	for _, v := range w.Connections {
		v.Send(name, data)
	}
}

type WsMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func (w *WsController) OnConnect(conn *websocket.Conn) {
	socket := &WsConnection{
		Raw:          conn,
		SID:          cuid.New(),
		IsIdentified: false,
		UID:          ``,
		CurrentRoom:  ``,
		Controller:   w,
		Connected:    true,
	}
	socket.Init()
	w.Connections = append(w.Connections, socket)
	for socket.Connected {
		var (
			mt  int
			msg []byte
			err error
		)

		if mt, msg, err = socket.Raw.ReadMessage(); err != nil {
			socket.Disconnect()
		}

		if mt != websocket.TextMessage {
			continue
		}

		var rawData WsMessage
		if json.Unmarshal(msg, &rawData) != nil {
			socket.Disconnect()
		}
		go socket.HandleMessage(rawData.Event, rawData.Data)
	}
}

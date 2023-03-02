/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-03-02 09:41:39
 * @LastEditTime: 2023-03-02 11:12:40
 */
package routes

import (
	"fmt"
	"sync/atomic"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/websocket"
)

// InitWebSocket 以MVC方式初始化websocket
func InitWebSocket(app iris.Party) {
	websocketAPI := app.Party("/websocket")

	m := mvc.New(websocketAPI)
	m.Register(
		&prefixedLogger{prefix: "[WebSocket]"},
	)
	m.HandleWebsocket(&websocketController{Namespace: "default", Age: 42, Otherstring: "other string"})

	websocketServer := websocket.New(websocket.DefaultGorillaUpgrader, m)

	websocketAPI.Get("/", websocket.Handler(websocketServer))
}

var visits uint64

func increment() uint64 {
	return atomic.AddUint64(&visits, 1)
}

func decrement() uint64 {
	return atomic.AddUint64(&visits, ^uint64(0))
}

type websocketController struct {
	*websocket.NSConn `stateless:"true"`
	Namespace         string
	Age               int
	Otherstring       string

	Logger LoggerService
}

// or
// func (c *websocketController) Namespace() string {
// 	return "default"
// }

func (c *websocketController) OnNamespaceDisconnect(msg websocket.Message) error {
	c.Logger.Log(fmt.Sprintf("[%s] disconnected from namespace [%s]", c.Conn, msg.Namespace))
	/*
		// visits--
		newCount := decrement()
		// This will call the "OnVisit" event on all clients, except the current one,
		// (it can't because it's left but for any case use this type of design)
		c.Conn.Server().Broadcast(nil, websocket.Message{
			Namespace: msg.Namespace,
			// Event:     "OnVisit",
			Body: []byte(fmt.Sprintf("%d", newCount)),
		})
	*/

	return nil
}

func (c *websocketController) OnNamespaceConnected(msg websocket.Message) error {
	// println("Broadcast prefix is: " + c.BroadcastPrefix)

	ctx := websocket.GetContext(c.Conn)

	c.Logger.Log(fmt.Sprintf("[%s] connected to namespace [%s] with IP [%s]",
		c.Conn, msg.Namespace,
		ctx.RemoteAddr()))
	/*
		// visits++
		newCount := increment()

		// This will call the "OnVisit" event on all clients, including the current one,
		// with the 'newCount' variable.
		//
		// There are many ways that u can do it and faster, for example u can just send a new visitor
		// and client can increment itself, but here we are just "showcasing" the websocket controller.
		c.Conn.Server().Broadcast(nil, websocket.Message{
			Namespace: msg.Namespace,
			// Event:     "OnVisit",
			Body: []byte(fmt.Sprintf("%d", newCount)),
		})
	*/
	return nil
}

func (c *websocketController) OnChat(msg websocket.Message) error {
	ctx := websocket.GetContext(c.Conn)

	ctx.Application().Logger().Infof("[IP: %s] [ID: %s]  broadcast to other clients the message [%s]",
		ctx.RemoteAddr(), c, string(msg.Body))

	c.Conn.Server().Broadcast(c, msg)

	return nil
}

type LoggerService interface {
	Log(string)
}

type prefixedLogger struct {
	prefix string
}

func (s *prefixedLogger) Log(msg string) {
	fmt.Printf("%s: %s\n", s.prefix, msg)
}

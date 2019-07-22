package majsoul

import (
	"github.com/gorilla/websocket"
	"github.com/EndlessCheng/mahjong-helper/tool"
	"net/http"
	"time"
	"github.com/golang/protobuf/proto"
	"sync"
	"fmt"
	"os"
	"encoding/binary"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
)

const (
	messageTypeNotify   = 1
	messageTypeRequest  = 2
	messageTypeResponse = 3
)

type chanMessage struct {
	message proto.Message
	done    chan struct{}
}

func newChanMessage(message proto.Message) *chanMessage {
	return &chanMessage{
		message: message,
		done:    make(chan struct{}),
	}
}

func (m *chanMessage) _done() {
	m.done <- struct{}{}
}

type rpcChannel struct {
	ws     *websocket.Conn
	closed bool

	messageIndex uint16
	callMap      *sync.Map // messageIndex -> *chanMessage
}

func newRpcChannel() *rpcChannel {
	return &rpcChannel{
		callMap: &sync.Map{},
	}
}

func (*rpcChannel) wrapMessage(name string, message proto.Message) (data []byte, err error) {
	data, err = proto.Marshal(message)
	if err != nil {
		return
	}
	wrap := lq.Wrapper{
		Name: name,
		Data: data,
	}
	return proto.Marshal(&wrap)
}

func (*rpcChannel) unwrapData(data []byte, message proto.Message) error {
	wrapper := lq.Wrapper{}
	if err := proto.Unmarshal(data, &wrapper); err != nil {
		return err
	}
	return proto.Unmarshal(wrapper.Data, message)
}

func (c *rpcChannel) run() {
	for !c.closed {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if len(data) <= 3 {
			fmt.Fprintln(os.Stderr, "数据过短", data)
			continue
		}

		if data[0] == messageTypeResponse {
			reqOrder := binary.LittleEndian.Uint16(data[1:3])
			rawMsg, ok := c.callMap.Load(reqOrder)
			if !ok {
				fmt.Fprintln(os.Stderr, "未找到消息", reqOrder)
				continue
			}
			msg := rawMsg.(*chanMessage)
			if err := c.unwrapData(data[3:], msg.message); err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			msg._done()
		}
	}
}

func (c *rpcChannel) connect(endpoint string, origin string) error {
	endPoint, err := tool.GetMajsoulWebSocketURL() // wss://mj-srv-7.majsoul.com:4131/
	if err != nil {
		return err
	}
	header := http.Header{}
	header.Set("origin", origin) // 模拟来源
	ws, _, err := websocket.DefaultDialer.Dial(endPoint, header)
	if err != nil {
		return err
	}
	c.ws = ws

	go c.run()
	return nil
}

func (c *rpcChannel) close() error {
	c.closed = true
	return c.ws.Close()
}

func (c *rpcChannel) send(reqMessage proto.Message, respChanMessage *chanMessage) error {
	// TODO: 反射 小写？
	name := ".lq.Lobby." + "login"
	data, err := c.wrapMessage(name, reqMessage)
	if err != nil {
		return err
	}

	c.messageIndex = (c.messageIndex + 1) % 60007 // from code.js

	// 填写消息序号后，发送登录请求给雀魂 WebSocket 服务器
	reqOrderBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(reqOrderBytes, c.messageIndex)
	msgHead := append([]byte{messageTypeRequest}, reqOrderBytes...)
	if err := c.ws.WriteMessage(websocket.BinaryMessage, append(msgHead, data...)); err != nil {
		return err
	}

	c.callMap.Store(c.messageIndex, respChanMessage)
	return nil
}

func (c *rpcChannel) heartbeat() error {
	for !c.closed {

		time.Sleep(6 * time.Second)
	}
	return nil
}

type GameAPI struct {
}

func NewGameAPI() *GameAPI {
	return &GameAPI{}
}

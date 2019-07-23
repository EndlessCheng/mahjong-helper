package majsoul

import (
	"encoding/binary"
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"
)

const (
	messageTypeNotify   = 1
	messageTypeRequest  = 2
	messageTypeResponse = 3
)

type rpcChannel struct {
	sync.Mutex

	ws     *websocket.Conn
	closed bool

	messageIndex       uint16
	respMessageChanMap *sync.Map // messageIndex -> chan proto.Message
}

func newRpcChannel() *rpcChannel {
	return &rpcChannel{
		respMessageChanMap: &sync.Map{},
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

func (*rpcChannel) unwrapData(rawData []byte) (methodName string, data []byte, err error) {
	wrapper := lq.Wrapper{}
	if err = proto.Unmarshal(rawData, &wrapper); err != nil {
		return
	}
	return wrapper.GetName(), wrapper.GetData(), nil
}

// TODO: auto unwrapMessage by methodName

func (c *rpcChannel) unwrapMessage(rawData []byte, message proto.Message) error {
	methodName, data, err := c.unwrapData(rawData)
	if err != nil {
		return err
	}
	// TODO: assert methodName
	_ = methodName
	return proto.Unmarshal(data, message)
}

func (c *rpcChannel) run() {
	for !c.closed {
		_, data, err := c.ws.ReadMessage()
		if err != nil {
			if c.closed {
				return
			}
			fmt.Fprintln(os.Stderr, "ws.ReadMessage:", err)
			continue
		}

		if len(data) <= 3 {
			fmt.Fprintln(os.Stderr, "数据过短", data)
			continue
		}

		if data[0] == messageTypeResponse {
			messageIndex := binary.LittleEndian.Uint16(data[1:3])
			rawRespMessageChan, ok := c.respMessageChanMap.Load(messageIndex)
			if !ok {
				fmt.Fprintln(os.Stderr, "未找到消息", messageIndex)
				continue
			}
			c.respMessageChanMap.Delete(messageIndex)

			respMessageType := reflect.TypeOf(rawRespMessageChan).Elem().Elem()
			respMessage := reflect.New(respMessageType)
			if err := c.unwrapMessage(data[3:], respMessage.Interface().(proto.Message)); err != nil {
				fmt.Fprintln(os.Stderr, "unwrapData:", err)
				reflect.ValueOf(rawRespMessageChan).Close()
				continue
			}
			reflect.ValueOf(rawRespMessageChan).Send(respMessage)
		}
	}
}

func (c *rpcChannel) connect(endpoint string, origin string) error {
	header := http.Header{}
	header.Set("origin", origin) // 模拟来源
	ws, _, err := websocket.DefaultDialer.Dial(endpoint, header)
	if err != nil {
		return err
	}
	c.ws = ws

	go c.run()
	go c.heartbeat()

	return nil
}

func (c *rpcChannel) close() error {
	c.closed = true
	return c.ws.Close()
}

func (c *rpcChannel) send(name string, reqMessage proto.Message, respMessageChan interface{}) error {
	c.Lock()
	defer c.Unlock()

	data, err := c.wrapMessage(name, reqMessage)
	if err != nil {
		return err
	}

	c.messageIndex = (c.messageIndex + 1) % 60007 // from code.js

	messageIndexBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(messageIndexBytes, c.messageIndex)
	messageHead := append([]byte{messageTypeRequest}, messageIndexBytes...)
	if err := c.ws.WriteMessage(websocket.BinaryMessage, append(messageHead, data...)); err != nil {
		return err
	}

	c.respMessageChanMap.Store(c.messageIndex, respMessageChan)
	return nil
}

func (c *rpcChannel) callFastTest(methodName string, reqMessage proto.Message, respMessageChan interface{}) error {
	return c.send(".lq.FastTest."+methodName, reqMessage, respMessageChan)
}

func (c *rpcChannel) callLobby(methodName string, reqMessage proto.Message, respMessageChan interface{}) error {
	return c.send(".lq.Lobby."+methodName, reqMessage, respMessageChan)
}

func (c *rpcChannel) heartbeat() {
	for !c.closed {
		// 吐槽：雀魂的开发把 heart 错写成了 heat
		reqHeartBeat := lq.ReqHeatBeat{}
		respCommonChan := make(chan *lq.ResCommon)
		if err := c.callLobby("heatbeat", &reqHeartBeat, respCommonChan); err != nil {
			fmt.Fprintln(os.Stderr, "heartbeat:", err)
		} else if respCommon := <-respCommonChan; respCommon.GetError() != nil {
			fmt.Fprintln(os.Stderr, "heartbeat:", respCommon.Error)
		}
		time.Sleep(6 * time.Second)
	}
}

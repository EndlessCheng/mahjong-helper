package majsoul

import (
	"github.com/golang/protobuf/proto"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/api"
	"fmt"
	"os"
	"reflect"
	"encoding/binary"
	"strings"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
)

// 若 NotifyMessage 不为空，这改消息为通知，仅包含 NotifyMessage，其余字段为空
type Message struct {
	MethodName      string
	RequestMessage  proto.Message
	ResponseMessage proto.Message

	NotifyMessage proto.Message
}

type MessageReceiver struct {
	originMessageQueue  chan []byte   // 包含所有 WebSocket 发出的消息和收到的消息
	orderedMessageQueue chan *Message // 整理后的 WebSocket 收到的消息（包含请求响应和通知）

	indexToMessageMap map[uint16]*Message
}

func NewMessageReceiver() *MessageReceiver {
	const maxQueueSize = 100
	mr := &MessageReceiver{
		originMessageQueue:  make(chan []byte, maxQueueSize),
		orderedMessageQueue: make(chan *Message, maxQueueSize),
		indexToMessageMap:   map[uint16]*Message{},
	}
	go mr.run()
	return mr
}

func (mr *MessageReceiver) run() {
	for data := range mr.originMessageQueue {
		messageType := data[0]
		switch messageType {
		case api.MessageTypeNotify:
			notifyName, data, err := api.UnwrapData(data[1:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "MessageReceiver.run.api.UnwrapData.NOTIFY", err)
				continue
			}

			notifyName = notifyName[1:] // 移除开头的 .
			mt := proto.MessageType(notifyName)
			if mt == nil {
				fmt.Fprintf(os.Stderr, "MessageReceiver.run 未找到 %s，请检查！\n", notifyName)
				continue
			}
			messagePtr := reflect.New(mt.Elem())
			if err := proto.Unmarshal(data, messagePtr.Interface().(proto.Message)); err != nil {
				fmt.Fprintln(os.Stderr, "MessageReceiver.run.proto.Unmarshal.NOTIFY", notifyName, err)
				continue
			}

			mr.orderedMessageQueue <- &Message{
				NotifyMessage: messagePtr.Interface().(proto.Message),
			}
		case api.MessageTypeRequest:
			messageIndex := binary.LittleEndian.Uint16(data[1:3])

			rawMethodName, data, err := api.UnwrapData(data[1:])
			if err != nil {
				fmt.Fprintln(os.Stderr, "MessageReceiver.run.api.UnwrapData.REQUEST", err)
				continue
			}

			// 通过 rawMethodName 找到请求结构体和请求响应结构体
			splits := strings.Split(rawMethodName[4:], ".") //  移除开头的 .lq.
			clientName, methodName := splits[0], splits[1]
			methodType := lq.FindMethod(clientName, methodName)
			reqType := methodType.In(1)
			respType := methodType.Out(0)

			messagePtr := reflect.New(reqType.Elem())
			if err := proto.Unmarshal(data, messagePtr.Interface().(proto.Message)); err != nil {
				fmt.Fprintln(os.Stderr, "MessageReceiver.run.proto.Unmarshal.REQUEST", rawMethodName, err)
				continue
			}
			reqMessage := messagePtr.Interface().(proto.Message)

			messagePtr = reflect.New(respType.Elem())
			respMessage := messagePtr.Interface().(proto.Message)

			mr.indexToMessageMap[messageIndex] = &Message{
				MethodName:      rawMethodName,
				RequestMessage:  reqMessage,
				ResponseMessage: respMessage,
			}
		case api.MessageTypeResponse:
			// 似乎是有序返回的……
			messageIndex := binary.LittleEndian.Uint16(data[1:3])
			message := mr.indexToMessageMap[messageIndex]
			if err := proto.Unmarshal(data, message.RequestMessage); err != nil {
				fmt.Fprintln(os.Stderr, "MessageReceiver.run.proto.Unmarshal.RESPONSE", message.MethodName, err)
				continue
			}
			mr.orderedMessageQueue <- message
		default:
			panic(fmt.Sprintln("[MessageReceiver] 数据有误", messageType))
		}
	}
}

func (mr *MessageReceiver) Put(data []byte) {
	mr.originMessageQueue <- data
}

func (mr *MessageReceiver) Get() *Message {
	return <-mr.orderedMessageQueue
}

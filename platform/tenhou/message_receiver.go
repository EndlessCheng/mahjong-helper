package tenhou

import (
	"time"
	"fmt"
	"os"
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou/ws"
)

type MessageReceiver struct {
	originMessageQueue  chan []byte
	orderedMessageQueue chan *message
}

func NewMessageReceiver() *MessageReceiver {
	const maxQueueSize = 100
	mr := &MessageReceiver{
		originMessageQueue:  make(chan []byte, maxQueueSize),
		orderedMessageQueue: make(chan *message, maxQueueSize),
	}
	go mr.run()
	return mr
}

// TODO: 合并短时间内的 AGARI 消息？（双响）
func (mr *MessageReceiver) run() {
	for data := range mr.originMessageQueue {
		msg, err := parse(data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if msg == nil {
			continue
		}

		if !isSelfDraw(msg.tag) {
			mr.orderedMessageQueue <- msg
			continue
		}

		// 收到了自家摸牌的消息，则等待一段很短的时间
		time.Sleep(75 * time.Millisecond) // 实际间隔在 3~9ms

		// 未收到新数据
		if len(mr.originMessageQueue) == 0 {
			mr.orderedMessageQueue <- msg
			continue
		}

		// 在短时间内收到了新数据
		// 因为摸牌后肯定要等待玩家操作，正常情况是不会马上有新数据的，所以这说明前端乱序发来了数据
		// 把 data 重新塞回去，这样才是正确的顺序
		mr.originMessageQueue <- data
	}
}

func (mr *MessageReceiver) Put(data []byte) {
	mr.originMessageQueue <- data
}

func (mr *MessageReceiver) Get() (message ws.Message, originJSON string) {
	m := <-mr.orderedMessageQueue
	return m.metadata, m.originJSON
}

func (mr *MessageReceiver) IsEmpty() bool {
	return len(mr.originMessageQueue) == 0 && len(mr.orderedMessageQueue) == 0
}

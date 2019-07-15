package tenhou

import (
	"time"
	"encoding/json"
	"regexp"
)

type MessageReceiver struct {
	originMessageQueue  chan []byte
	orderedMessageQueue chan []byte
}

func NewMessageReceiver() *MessageReceiver {
	const maxQueueSize = 100
	mr := &MessageReceiver{
		originMessageQueue:  make(chan []byte, maxQueueSize),
		orderedMessageQueue: make(chan []byte, maxQueueSize),
	}
	go mr.run()
	return mr
}

var isSelfDraw = regexp.MustCompile("^T[0-9]{1,3}$").MatchString

// TODO: 后续使用 parser 中提供的方法
func (mr *MessageReceiver) isSelfDraw(data []byte) bool {
	d := struct {
		Tag string `json:"tag"`
	}{}
	if err := json.Unmarshal(data, &d); err != nil {
		return false
	}
	return isSelfDraw(d.Tag)
}

func (mr *MessageReceiver) run() {
	for data := range mr.originMessageQueue {
		if !mr.isSelfDraw(data) {
			mr.orderedMessageQueue <- data
			continue
		}

		// 收到了自家摸牌的消息，则等待一段很短的时间
		time.Sleep(75 * time.Millisecond) // 实际间隔在 3~9ms

		// 未收到新数据
		if len(mr.originMessageQueue) == 0 {
			mr.orderedMessageQueue <- data
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

func (mr *MessageReceiver) Get() []byte {
	return <-mr.orderedMessageQueue
}

func (mr *MessageReceiver) IsEmpty() bool {
	return len(mr.originMessageQueue) == 0 && len(mr.orderedMessageQueue) == 0
}

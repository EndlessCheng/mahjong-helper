package tenhou

import (
	"testing"
	"time"
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou/ws"
	"github.com/stretchr/testify/assert"
)

func TestMessageReceiver(t *testing.T) {
	mr := NewMessageReceiver()

	done := make(chan struct{})
	go func() {
		indexes := []int{}
		for i := 0; ; i++ {
			select {
			case <-done:
				assert.EqualValues(t, []int{3, 6, 8}, indexes)
				t.Log("DONE")
				return
			default:
				msg := mr.Get()
				if _, ok := msg.Metadata.(*ws.Draw); ok {
					indexes = append(indexes, i)
				}
			}
		}
	}()

	mr.Put([]byte(`{"Tag":"INIT"}`))
	mr.Put([]byte(`{"Tag":"T1"}`))
	mr.Put([]byte(`{"Tag":"N"}`))
	mr.Put([]byte(`{"Tag":"c"}`))
	time.Sleep(time.Second)

	mr.Put([]byte(`{"Tag":"d"}`))
	mr.Put([]byte(`{"Tag":"e"}`))
	mr.Put([]byte(`{"Tag":"T2"}`))
	time.Sleep(time.Second)

	mr.Put([]byte(`{"Tag":"T3"}`))
	mr.Put([]byte(`{"Tag":"f"}`))
	time.Sleep(time.Second)

	mr.Put([]byte(`{"Tag":"a"}`))
	done <- struct{}{}
	time.Sleep(time.Second)
}

package tenhou

import (
	"testing"
	"time"
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou/ws"
	"github.com/stretchr/testify/assert"
	"sync"
)

func TestMessageReceiver(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	done := make(chan struct{})
	defer close(done)

	mr := NewMessageReceiver()

	go func() {
		wg.Add(1)
		var indexes []int
		for i := 0; ; i++ {
			select {
			case msg := <-mr.GetC():
				if _, ok := msg.Metadata.(*ws.Draw); ok {
					indexes = append(indexes, i)
				}
			case <-done:
				if assert.EqualValues(t, []int{3, 6, 8}, indexes) {
					t.Log("DONE")
				}
				wg.Done()
				return
			}
		}
	}()

	// 模拟发来的消息
	mr.Put([]byte(`{"Tag":"INIT"}`))
	mr.Put([]byte(`{"Tag":"T1"}`))
	mr.Put([]byte(`{"Tag":"N"}`))
	mr.Put([]byte(`{"Tag":"c"}`))
	time.Sleep(500 * time.Millisecond)
	mr.Put([]byte(`{"Tag":"d"}`))
	mr.Put([]byte(`{"Tag":"e"}`))
	mr.Put([]byte(`{"Tag":"T2"}`))
	time.Sleep(500 * time.Millisecond)
	mr.Put([]byte(`{"Tag":"T3"}`))
	mr.Put([]byte(`{"Tag":"f"}`))
	time.Sleep(500 * time.Millisecond)
}

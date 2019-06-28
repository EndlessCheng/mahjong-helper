package tenhou

import (
	"testing"
	"time"
	"fmt"
)

func TestMessageReceiver(t *testing.T) {
	mr := NewMessageReceiver()

	go func() {
		for {
			fmt.Println(string(mr.Get()))
		}
	}()

	mr.Put([]byte(`{"Tag":"a"}`))
	mr.Put([]byte(`{"Tag":"T1"}`))
	mr.Put([]byte(`{"Tag":"b"}`))
	mr.Put([]byte(`{"Tag":"c"}`))
	time.Sleep(time.Second)

	mr.Put([]byte(`{"Tag":"d"}`))
	mr.Put([]byte(`{"Tag":"e"}`))
	mr.Put([]byte(`{"Tag":"T2"}`))
	time.Sleep(time.Second)

	mr.Put([]byte(`{"Tag":"T3"}`))
	mr.Put([]byte(`{"Tag":"f"}`))
	time.Sleep(time.Second)
}

package tool

import "testing"

func TestGetMajsoulWebSocketURL(t *testing.T) {
	u, err := GetMajsoulWebSocketURL()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

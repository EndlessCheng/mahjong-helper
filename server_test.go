package main

import (
	"testing"
	"time"
)

func Test_mjHandler_runAnalysisTenhouMessageTask(t *testing.T) {
	debugMode = true

	h := &mjHandler{
		tenhouMessageQueue: make(chan []byte, 100),
		tenhouRoundData:    &tenhouRoundData{isRoundEnd: true},
	}
	h.tenhouRoundData.roundData = newGame(h.tenhouRoundData)

	h.tenhouMessageQueue <- []byte("{\"tag\":\"T8\"}")
	h.tenhouMessageQueue <- []byte("{\"tag\":\"INIT\",\"seed\":\"1,1,0,2,0,27\",\"ten\":\"318,224,224,234\",\"oya\":\"0\",\"hai\":\"129,90,47,39,4,9,116,53,33,123,69,28,14\"}")
	go h.runAnalysisTenhouMessageTask()

	time.Sleep(time.Second)
}

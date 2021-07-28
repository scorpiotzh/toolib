package toolib

import (
	"fmt"
	"os"
	"testing"
)

func TestExitMonitoring(t *testing.T) {
	exit := make(chan struct{})
	ExitMonitoring(func(sig os.Signal) {
		fmt.Println("sys exit ok ... ")
		exit <- struct{}{}
	})
	<-exit
}

func TestJsonString(t *testing.T) {
	type TestJson struct {
		A string            `json:"a"`
		B int               `json:"b"`
		C []string          `json:"c"`
		D map[string]string `json:"d"`
	}
	var tj TestJson
	tj.A = "A"
	tj.B = 1
	tj.C = []string{"a", "b"}
	tj.D = map[string]string{
		"a": "1",
		"b": "2",
	}
	fmt.Println(JsonString(tj))
	fmt.Println(JsonString(tj.C))
	fmt.Println(JsonString(tj.D))
}

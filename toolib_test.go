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

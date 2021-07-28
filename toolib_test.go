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

func TestNewGormDataBase(t *testing.T) {
	db, err := NewGormDataBase("127.0.0.1:3306", "root", "tzh123456", "das_db", 100, 100)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(db)
}

func TestNewRedisClient(t *testing.T) {
	red, err := NewRedisClient("127.0.0.1:6379", "", 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(red)
}

func TestUnmarshalYamlFile(t *testing.T) {
	type YamlTest struct {
		Server struct {
			Port string `json:"port" yaml:"port"`
		} `json:"server" yaml:"server"`
	}
	var yt YamlTest
	err := UnmarshalYamlFile("test.yaml", &yt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(JsonString(yt))
}

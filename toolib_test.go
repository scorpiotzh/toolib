package toolib

import (
	"fmt"
	"os"
	"testing"
	"time"
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

func TestNewGormDB(t *testing.T) {
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

func TestAllowOriginFunc(t *testing.T) {
	fmt.Println(AllowOriginFunc("http://127.0.0.1:"))
	fmt.Println(AllowOriginFunc("http://localhost:80"))
	fmt.Println(AllowOriginFunc("http://test.com"))
	fmt.Println(AllowOriginFunc("https://a.test.com"))
}

func TestPage(t *testing.T) {
	p := Pagination{
		Page: 2,
		Size: 50,
	}
	fmt.Println(p.GetLimit(), p.GetOffset())
}

func TestAddFileWatcher(t *testing.T) {
	exit := make(chan struct{})
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
	//
	w, err := AddFileWatcher("test.yaml", func() {
		err := UnmarshalYamlFile("test.yaml", &yt)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(JsonString(yt))
	})
	if err != nil {
		t.Fatal(err)
	}
	//
	ExitMonitoring(func(sig os.Signal) {
		if w != nil {
			fmt.Println("close")
			w.Close()
		}
		fmt.Println("sys exit ok ... ")
		exit <- struct{}{}
	})
	//
	<-exit
}

func TestTimeFormat(t *testing.T) {
	fmt.Println(TimeFormat(time.Now()))
}

func TestJwtString(t *testing.T) {
	jwtKey := "toolib"
	token, err := JwtSimple(time.Second*6, "111", "", "", "", jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
	claims, err := JwtVerify(token, jwtKey)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("OK", claims)
	}
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjg1MTAzMTksImlhdCI6MTYyODUxMDMxNCwibmJmIjoxNjI4NTEwMzE0fQ.CmP7kBfc-BKCdFE7jFEVi4jEoMMpjUtIzbAbO1SKtMo
	//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mjg1MTAzMzcsImlhdCI6MTYyODUxMDMzMSwibmJmIjoxNjI4NTEwMzMxfQ.DoX9wM8EOKHETc1hmqNuZtXT8JDs3xqVyFAAOyFXsrE
	//time.Sleep(time.Second * 6)
	//claims, err = JwtVerify(token, jwtKey)
	//if err != nil {
	//	t.Fatal(err)
	//} else {
	//	fmt.Println("OK", claims)
	//}
}

func TestSendEmail(t *testing.T) {
	eh := EmailHelper{
		Host:           "smtpdm.aliyun.com",
		port:           465,
		From:           "test",
		SenderAddress:  "noreply@cctip.io",
		SenderPassword: "",
	}
	err := eh.SendEmail("test", "test \n 111", "duzhihongyi@gmail.com")
	fmt.Println(err)
}

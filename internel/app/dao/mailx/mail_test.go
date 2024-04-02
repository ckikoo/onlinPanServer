package mailx

import (
	"fmt"
	"testing"
)

func TestMail(t *testing.T) {
	Init()
	for i := 0; i < 1; i++ {

		// err := Email.DebugSendMsgwithHtml(context.Background(), "lj_5683@163.com", "验证码", "2")
		// if err != nil {
		// panic(err)
		// }
	}
}

var str = "lj_5683@163.com"

func TestMailNo(t *testing.T) {
	for i := 0; i < 1000; i++ {
		s := "online cloud1111<" + str + ">"
		fmt.Printf("s: %v\n", s)
	}
}

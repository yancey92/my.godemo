package jpushkit

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPush(t *testing.T) {

	cl := NewClient("2abbc85497d262cd67fe10e4", "df4bc534cd1c0f024b994885")

	msg := NewMessage()
	msg.AddExtras("key", "时间标识")

	audience := Audience{}
	ids := make([]string, 0)
	ids = append(ids, "18071adc0308490854e")
	audience.RegistrationID = ids

	//u1 := uuid.NewV4()
	req := NewPushRequest()
	req.SetPushMsg(msg)
	req.SetAudience(&audience)

	content, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(content))

	m, err := cl.Push(req)
	if err != nil {
		fmt.Printf(fmt.Sprintf("cl.push error: %v, return value:%#v", err, m))
	}
	fmt.Printf("成功%v", m)
}

package pusher

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	Start("192.168.33.10:28015", "mercury", "192.168.33.10:6379", true)
}

// func TestFindChannelByHash(t *testing.T) {
// 	res, err := FindChannelByHash("05cf281c050be3da4eecf3bc6e8aac1b")
// 	if err != nil {
// 		t.Fail()
// 	}

// 	fmt.Printf("%#v", res)
// }

// func TestCreateChannel(t *testing.T) {
// 	res, err := CreateChannel("abc", []string{"1", "2"})
// 	if err != nil {
// 		t.Fail()
// 	}

// 	fmt.Printf("%#v", res)
// }

func TestGetChannel(t *testing.T) {
	c, err := GetChannelByMembers([]string{"1", "2"})
	if err != nil {
		t.Fail()
	}

	assert.NoError(t, err, "...")

	fmt.Printf("%#v\n", c)
}

func TestFindMessage(t *testing.T) {
	m, err := FindMessage("05be8e3a-dfbd-49df-b6c4-fdb276ad586b")
	if err != nil {
		t.Fail()
	}

	fmt.Printf("%#v\n", m)
}

func TestCreateMessage(t *testing.T) {
	payload := map[string]interface{}{"name": "hello", "age": 12}
	m := NewMessage("abc", 1, payload, "1", nil)
	r, err := CreateMessage(m)
	if err != nil {
		t.Fail()
	}

	fmt.Printf("%#v\n", r)

	m2, _ := FindMessage(m.Id)
	if m2.CreatedAt != m.CreatedAt {
		fmt.Println(m.CreatedAt)
		fmt.Println(m2.CreatedAt)
		t.Error("CreatedAt is not equal")
	}
}

func TestGetUserIdByToke(t *testing.T) {
	token := "test"
	redis.Cmd("hmset", token, "user_id", 123)

	userId, err := GetUserIdByToken(token)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if userId != 123 {
		fmt.Println(userId)
		t.Fail()
	}
}

func TestJSON(t *testing.T) {
	// s := `{"type":1,"sender_id":"1000","channel_id":"0a677b69-500c-4a95-a7e0-47aafb116248","payload":"{\"id\":\"rWAwnseKWXtd8FyN\",\"type\":1,\"sub_type\":1001,\"chat_id\":\"294\",\"sender_id\":\"1000\",\"ttl\":0,\"sent_at\":\"2015-03-29 19:03:24\",\"body\":{\"mime\":\"text\",\"content\":{\"text\":\"38\"}},\"extra\":{\"sender_name\":\"\\u5ba2\\u670d\",\"sender_avatar\":\"http:\\\/\\\/static.lover1314.me\\\/icons\\\/icon.png\"}}","push_offline":true,"opts":{"apn_alert":"\u5ba2\u670d:38"}}`

	// var msg *Message
	// err := json.Unmarshal([]byte(s), &msg)
	// if err != nil {
	// 	t.Fail()
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Printf("%#v", msg)
}

package pusher

import (
	"fmt"
	"testing"
)

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

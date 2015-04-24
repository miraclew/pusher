package main

import (
	"fmt"
	r "github.com/garyburd/redigo/redis"
	"testing"
)

func TestRedisSadd(t *testing.T) {
	c, err := r.Dial("tcp", "192.168.33.10:6379")
	if err != nil {
		t.Fatal("redis connect failed")
	}

	defer c.Close()

	members := []string{"a", "b", "dd"}
	args := r.Args{}.Add("ttt").AddFlat(members)
	fmt.Printf("%#v", args)

	res, err := c.Do("sadd", args...)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	fmt.Println(res)
}

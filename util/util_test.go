package util

import (
	"log"
	"testing"
)

func TestVersionCompare(t *testing.T) {
	r, err := VersionCompare("2.0", "1.9.4")
	if r <= 0 || err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	r, err = VersionCompare("2.0.5", "2.1.0")
	if r >= 0 || err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	r, err = VersionCompare("2.0", "2.0.0")
	if r != 0 || err != nil {
		log.Println(err.Error())
		t.Fail()
	}

}

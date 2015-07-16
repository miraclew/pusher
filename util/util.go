package util

import (
	"strconv"
	"strings"
)

func VersionCompare(v1 string, v2 string) (int64, error) {
	v1List := strings.Split(v1, ".")
	v2List := strings.Split(v2, ".")
	maxLen := len(v1List)
	if maxLen < len(v2List) {
		maxLen = len(v2List)
	}

	for i := 0; i < maxLen; i++ {
		s1 := "0"
		if i < len(v1List) {
			s1 = v1List[i]
		}

		s2 := "0"
		if i < len(v2List) {
			s2 = v2List[i]
		}

		p1, err1 := strconv.ParseInt(s1, 10, 64)
		p2, err2 := strconv.ParseInt(s2, 10, 64)
		if err1 != nil {
			return 0, err1
		}
		if err2 != nil {
			return 0, err2
		}

		if p1 != p2 {
			return p1 - p2, nil
		}
	}

	return 0, nil
}

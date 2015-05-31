package pusher

import (
	"errors"
	"strconv"
	"strings"
)

func VersionCompare(v1 string, v2 string) (int64, error) {
	v1List := strings.Split(v1, ".")
	v2List := strings.Split(v2, ".")
	if len(v1List) != len(v2List) {
		return 0, errors.New("version length not equal")
	}

	for i := 0; i < len(v1List); i++ {
		p1, err1 := strconv.ParseInt(v1List[i], 10, 64)
		p2, err2 := strconv.ParseInt(v2List[i], 10, 64)
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

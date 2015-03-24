package main

import (
	"fmt"
	"github.com/miraclew/mrs/util"
	// "sort"
	// "strings"
)

// func RemoveDuplicates(xs *[]string) {
// 	found := make(map[string]bool)
// 	j := 0
// 	for i, x := range *xs {
// 		if !found[x] {
// 			found[x] = true
// 			(*xs)[j] = (*xs)[i]
// 			j++
// 		}
// 	}
// 	*xs = (*xs)[:j]
// }

func main() {
	strMembers := "a,d,  b,c, a,"
	members := util.SplitUniqSort(strMembers)
	fmt.Printf("%#v\n", members)
	// members := strings.Split(strMembers, ",")
	// fmt.Printf("%#v\n", members)
	// var newMembers []string
	// for i := 0; i < len(members); i++ {
	// 	m := members[i]
	// 	m = strings.TrimSpace(m)
	// 	if len(m) <= 0 {
	// 		continue
	// 	}
	// 	newMembers = append(newMembers, m)
	// }

	// fmt.Printf("%#v\n", newMembers)
	// RemoveDuplicates(&newMembers)
	// sorted := sort.StringSlice(newMembers)
	// sorted.Sort()
	// fmt.Println(sorted)
}

package util

import "fmt"

func ContainStr(sli []string, v string) bool {
	for _, item := range sli {
		if item == v {
			return true
		}
	}
	return false
}

func PressAnyKeyToExit() {
	var x string
	fmt.Println("Press any key to exit.")
	fmt.Scan(&x)
}

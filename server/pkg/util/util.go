package util

func ContainStr(sli []string, v string) bool {
	for _, item := range sli {
		if item == v {
			return true
		}
	}
	return false
}

package helpers

// StringInSlice : array check string
func StringInSlice(search string, list []string) bool {
	for _, str := range list {
		if str == search {
			return true
		}
	}
	return false
}

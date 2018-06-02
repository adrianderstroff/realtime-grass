package engine

func NullStrings(strings []string) []string {
	var nullStrings []string
	for _, str := range strings {
		nullStrings = append(nullStrings, str)
	}
	return nullStrings
}

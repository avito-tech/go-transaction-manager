package common

// CheckErr throws panic if there is an error.
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

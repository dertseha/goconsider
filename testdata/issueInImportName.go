package testdata

import (
	abcd "fmt" // want `Package alias contains 'abcd', consider rephrasing to something else`
)

func PrintSomething() {
	abcd.Println("")
}

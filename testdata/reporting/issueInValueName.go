package reporting

import "fmt"

const AbcdConstant = 1234 // want `Value name contains 'abcd', consider rephrasing to something else`

var abcdGlobal = 1234 // want `Value name contains 'abcd', consider rephrasing to something else`

func ConstantFunc() {
	const LocalAbcdConstant = "" // want `Value name contains 'abcd', consider rephrasing to something else`
	fmt.Println(LocalAbcdConstant)
}

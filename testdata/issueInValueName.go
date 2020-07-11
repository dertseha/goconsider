package testdata

import "fmt"

const AbcdConstant = 1234

var abcdGlobal = 1234

func ConstantFunc() {
	const LocalAbcdConstant = ""
	fmt.Println(LocalAbcdConstant)
}

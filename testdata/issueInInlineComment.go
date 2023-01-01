package testdata

const (
	// someConstant does things with abcd. // want `Comment contains 'abcd', consider rephrasing to something else`
	someConstant = 123 // It should abcd. // want `Comment contains 'abcd', consider rephrasing to something else`
)

func someFunc(value int) bool {
	return value%2 == 0 // This is abcd. // want `Comment contains 'abcd', consider rephrasing to something else`
}

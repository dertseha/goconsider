package reporting

type AbcdType int // want `Type name contains 'abcd', consider rephrasing to something else`

type WrappedType AbcdType

func (a AbcdType) Function() {
}

func WorkWithIt(a AbcdType) AbcdType {
	result := AbcdType(0)
	return result
}

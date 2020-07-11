package testdata

type AbcdType int

type WrappedType AbcdType

func (a AbcdType) Function() {
}

func WorkWithIt(a AbcdType) AbcdType {
	result := AbcdType(0)
	return result
}

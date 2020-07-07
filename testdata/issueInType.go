package testdata

type AbcdThing struct {
	MemberNamedAbcd int
}

type SpecialAbcdFunc func(abcdParam string) (resultAbcd int)

type AbcdInterface interface {
	AbcdFunc(abcdParam string) (resultAbcd int)
}

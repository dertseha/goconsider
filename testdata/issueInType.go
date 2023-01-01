package testdata

type AbcdThing struct { // want `Type name contains 'abcd', consider rephrasing to something else`
	MemberNamedAbcd int // want `Member name contains 'abcd', consider rephrasing to something else`
}

type SpecialAbcdFunc func()                                        // want `Type name contains 'abcd', consider rephrasing to something else`
type SpecialSafeFuncParam func(abcdParam string)                   // want `Parameter name contains 'abcd', consider rephrasing to something else`
type SpecialSafeFuncResult func(safeParam string) (resultAbcd int) // want `Result name contains 'abcd', consider rephrasing to something else`

type AbcdInterface interface { // want `Type name contains 'abcd', consider rephrasing to something else`
	AbcdFunc(safeParam string) (safeResult int)       // want `Method name contains 'abcd', consider rephrasing to something else`
	SafeFuncParam(abcdParam string) (safeResult int)  // want `Parameter name contains 'abcd', consider rephrasing to something else`
	SafeFuncResult(safeParam string) (abcdResult int) // want `Result name contains 'abcd', consider rephrasing to something else`
}

package testdata

type TypeWithMethod int

func (abcdReceiver TypeWithMethod) safeFuncReceiver() { // want `Function receiver contains 'abcd', consider rephrasing to something else`
}

func (safeReceiver TypeWithMethod) abcdFuncName() { // want `Function name contains 'abcd', consider rephrasing to something else`
}

func (safeReceiver TypeWithMethod) safeFuncParam(abcdParam string) { // want `Parameter name contains 'abcd', consider rephrasing to something else`
	ignored := func(abcdHelper string) int { // want `Parameter name contains 'abcd', consider rephrasing to something else`
		return len(abcdHelper)
	}
	ignored(abcdParam)
}

func (safeReceiver TypeWithMethod) safeFuncName() (abcd string) { // want `Result name contains 'abcd', consider rephrasing to something else`
	ignored := func() (abcd string) { // want `Result name contains 'abcd', consider rephrasing to something else`
		return ""
	}
	ignored()
	return ""
}

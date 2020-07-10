package testdata

type TypeWithMethod int

func (abcdReceiver TypeWithMethod) SpecificAbcdFunc(abcdParam string) (resultAbcd int) {
	processAbcd := func(abcdHelper string) int {
		return len(abcdHelper)
	}
	abcdTemp := processAbcd(abcdParam)
	resultAbcd = abcdTemp * 2

	return
}

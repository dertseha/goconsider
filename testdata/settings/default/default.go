package _default

// MasterFunc showcases that it is found by default settings. // want `Comment contains 'master', consider rephrasing to one of \['primary', 'leader', 'main'\].`
func MasterFunc() {} // want `Function name contains 'master', consider rephrasing to one of \['primary', 'leader', 'main'\].`

// AbcFunc will be ignored by default settings.
func AbcFunc() {}

// XyzFunc will be ignored by default settings.
func XyzFunc() {}

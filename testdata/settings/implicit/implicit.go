package implicit

// MasterFunc will be ignored by implicit settings
func MasterFunc() {}

// AbcFunc will be ignored by implicit settings.
func AbcFunc() {}

// XyzFunc is marked with the implicit settings file, found in the current working directory. // want `Comment contains 'xyz', consider rephrasing to one of \[this, that\]`
func XyzFunc() {} // want `Function name contains 'xyz', consider rephrasing to one of \[this, that\]`

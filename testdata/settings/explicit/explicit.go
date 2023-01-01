package explicit

// MasterFunc will be ignored by explicit settings
func MasterFunc() {}

// AbcFunc shows that settings files can be specified explicitly. // want `Comment contains 'abc', consider rephrasing to one of \[def, ghi\]`
func AbcFunc() {} // want `Function name contains 'abc', consider rephrasing to one of \[def, ghi\]`

// XyzFunc will be ignored by explicit settings.
func XyzFunc() {}

package eflag

// StringList represents a list of strings that can be parsed from a comma-separated string.
type StringList struct {
	p     string   // p is for flag set
	value []string // The list of strings parsed from p, a comma-separated string
}

// setValue sets the value of the StringList by splitting the p with commas.
func (sl *StringList) setValue() {
	if sl.p != "" {
		sl.value = SplitWithComma(sl.p)
	}
}

// Value returns the current list of strings stored in the StringList.
func (sl *StringList) Value() []string {
	return sl.value
}

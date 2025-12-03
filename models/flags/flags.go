package flags

import "strings"

type Flags string

func (f *Flags) Has(flags ...string) bool {
	for _, flag := range flags {
		if !strings.Contains(string(*f), flag) {
			return false
		}
	}
	return true
}

func (f *Flags) Set(flag string) {
	if !strings.Contains(string(*f), flag) {
		*f = Flags(string(*f) + flag)
	}
}

func (f *Flags) Remove(flag string) {
	*f = Flags(strings.ReplaceAll(string(*f), flag, ""))
}

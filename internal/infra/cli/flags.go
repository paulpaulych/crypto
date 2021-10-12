package cli

import (
	"flag"
	"strings"
)

type FlagsSpec struct {
	name           string
	keysWithUsages map[string]string
}

func NewFlagSpec(name string, keysWithUsages map[string]string) *FlagsSpec {
	return &FlagsSpec{
		name:           name,
		keysWithUsages: keysWithUsages,
	}
}

type Flags struct {
	Args  []string
	Flags map[string]*Flag
}

type Flag struct {
	value string
}

func (f *Flag) GetOr(defaultVal string) string {
	if f == nil {
		return defaultVal
	} else {
		return f.value
	}
}
func (f *Flag) Get() *string {
	if f == nil {
		return nil
	} else {
		return &f.value
	}
}

func (spec *FlagsSpec) Parse(args []string) (*Flags, CmdConfError) {
	fs := flag.NewFlagSet(spec.name, flag.ContinueOnError)

	var rawFlags = map[string]*string{}
	for key, usage := range spec.keysWithUsages {
		rawFlags[key] = fs.String(key, "", usage)
	}

	usageBuilder := &strings.Builder{}
	fs.SetOutput(usageBuilder)
	err := fs.Parse(args)
	if err != nil {
		usage := usageBuilder.String()
		return nil, NewCmdConfError(err.Error(), &usage)
	}

	flags := map[string]*Flag{}
	for key, value := range rawFlags {
		if value == nil || *value == "" {
			flags[key] = nil
		} else {
			flags[key] = &Flag{value: *value}
		}
	}
	return &Flags{
		Flags: flags,
		Args:  fs.Args(),
	}, nil
}

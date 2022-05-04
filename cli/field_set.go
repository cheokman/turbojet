package cli

import "fmt"

type FlagSet struct {
	flags []*Flag

	index map[string]*Flag
}

func NewFlagSet() *FlagSet {
	return &FlagSet{
		flags: make([]*Flag, 0),
		index: make(map[string]*Flag),
	}
}

func (fs *FlagSet) Flags() []*Flag {
	return fs.flags
}

func (fs *FlagSet) Add(f *Flag) {
	f.checkValid()
	for _, s := range f.GetFormations() {
		if _, ok := fs.index[s]; ok {
			panic(fmt.Errorf("flag duplicated %s", string(s)))
		}
		fs.index[s] = f
	}
	fs.flags = append(fs.flags, f)
}

func (fs *FlagSet) AddByName(name string) (*Flag, error) {
	if _, ok := fs.index["--"+name]; ok {
		return nil, fmt.Errorf("flag duplcated --%s", name)
	}

	f := &Flag{
		Name: name,
	}

	fs.Add(f)
	return f, nil
}

func (fs *FlagSet) Get(name string) *Flag {
	if f, ok := fs.index["--"+name]; ok {
		f.formation = "--" + name
		return f
	}
	return nil
}

func (fs *FlagSet) GetByShorthand(c rune) *Flag {
	if f, ok := fs.index["-"+string(c)]; ok {
		f.formation = "-" + string(c)
		return f
	}
	return nil
}

func (fs *FlagSet) GetValue(name string) (string, bool) {
	f := fs.Get(name)
	if f == nil {
		return "", false
	}
	return f.GetValue()
}

func (fs *FlagSet) assignedCount() int {
	n := 0
	for _, f := range fs.flags {
		if f.assigned {
			n++
		}
	}
	return n
}

func (fs *FlagSet) mergeWith(from *FlagSet, applier func(f *Flag) bool) *FlagSet {
	if from == nil {
		return fs
	}
	r := NewFlagSet()
	if fs != nil {
		for _, f := range fs.flags {
			r.Add(f)
		}
	}
	for _, rv := range from.Flags() {
		if applier(rv) {
			r.put(rv)
		}
	}
	return r
}

func (fs *FlagSet) put(f *Flag) {

	for _, lv := range fs.flags {
		if lv == f {
			return
		}
	}
	fs.flags = append(fs.flags, f)
	for _, s := range f.GetFormations() {
		fs.index[s] = f
	}
}

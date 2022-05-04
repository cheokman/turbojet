package cli

import (
	"fmt"
	"strconv"
)

type AssignedMode int

const (
	AssignedNone       = AssignedMode(-1)
	AssignedDefault    = AssignedMode(0)
	AssignedOnce       = AssignedMode(1)
	AssignedRepeatable = AssignedMode(9)
)

type Flag struct {
	Name string

	Shorthand rune

	Short string

	Long string

	DefaultValue string

	DeafultValue string

	Required bool

	Aliases []string

	ExcludeWith []string

	Persistent bool

	AssignedMode AssignedMode

	Category string

	Validate func(f *Flag) error

	Fields []Field

	assigned  bool
	value     string
	values    []string
	formation string
}

func (f *Flag) IsAssigned() bool {
	return f.assigned
}

func (f *Flag) SetAssigned(istrue bool) {
	f.assigned = istrue
}

func (f *Flag) SetValue(value string) {
	f.value = value
}

func (f *Flag) GetValue() (string, bool) {
	if f.IsAssigned() {
		return f.value, true
	} else if f.Required {
		return f.DeafultValue, false
	} else {
		return "", false
	}
}

func (f *Flag) GetValues() []string {
	return f.values
}

func (f *Flag) SetValues(values []string) {
	f.values = values
}

func (f *Flag) getField(key string) (*Field, bool) {
	for i, field := range f.Fields {
		if field.Key == key {
			return &(f.Fields[i]), true
		}
	}
	return nil, false
}

func (f *Flag) GetFieldValue(key string) (string, bool) {
	if field, ok := f.getField(key); ok {
		return field.getValue()
	}
	return "", false
}

func (f *Flag) GetFieldValues(key string) []string {
	if field, ok := f.getField(key); ok {
		return field.values
	}
	return make([]string, 0)
}

func (f *Flag) GetStringOrDefault(def string) string {
	if f == nil {
		return def
	}
	if f.assigned {
		return f.value
	}

	return def
}

func (f *Flag) GetIntegerOrDefault(def int) int {
	if f == nil {
		return def
	}

	if f.assigned {
		if i, err := strconv.Atoi(f.value); err == nil {
			return i
		}
	}

	return def
}

func (f *Flag) GetFormations() []string {
	r := make([]string, 0)
	if f.Name != "" {
		r = append(r, "--"+f.Name)
	}
	for _, s := range f.Aliases {
		r = append(r, "--"+s)
	}

	if f.Shorthand != 0 {
		r = append(r, "-"+string(f.Shorthand))
	}

	return r
}

func (f *Flag) setIsAssigned() error {
	if !f.assigned {
		f.assigned = true
	} else {
		if f.AssignedMode != AssignedRepeatable {
			return fmt.Errorf("%s duplicated", f.formation)
		}
	}
	return nil
}

func (f *Flag) needValue() bool {
	switch f.AssignedMode {
	case AssignedNone:
		return false
	case AssignedDefault:
		return f.value == ""
	case AssignedOnce:
		return f.value == ""
	case AssignedRepeatable:
		return true
	default:
		panic(fmt.Errorf("unexpected Flag.AssignedMode %s", strconv.Itoa(int(f.AssignedMode))))
	}
}

func (f *Flag) checkValid() {
	if len(f.Fields) > 0 {
		if f.AssignedMode != AssignedRepeatable {
			panic(fmt.Errorf("flag %s with fields must use AssignedRepeatable", f.Name))
		}
	}
}

func (f *Flag) validate() error {
	if f.AssignedMode == AssignedOnce && f.value == "" {
		return fmt.Errorf("%s must be assigned with value", f.formation)
	}
	return nil
}

func (f *Flag) assign(v string) error {
	if f.AssignedMode == AssignedNone {
		return fmt.Errorf("flag --%s can't be assiged", f.Name)
	}

	f.assigned = true
	f.value = v

	if f.AssignedMode == AssignedRepeatable {
		f.values = append(f.values, v)
		if len(f.Fields) > 0 {
			f.assignField(v)
		}
	}
	return nil
}

func (f *Flag) assignField(s string) error {
	if k, v, ok := SplitStringWithPrefix(s, "="); ok {
		field, ok2 := f.getField(k)
		if ok2 {
			field.assign(v)
		} else {
			return fmt.Errorf("--%s can't assign with %s=", f.Name, k)
		}
	} else {
		field, ok2 := f.getField("")
		if ok2 {
			field.assign(v)
		} else {
			return fmt.Errorf("--%s can't assign with value", f.Name)
		}
	}
	return nil
}

func (f *Flag) checkFields() error {
	if len(f.Fields) == 0 {
		return nil
	}
	for _, field := range f.Fields {
		if err := field.check(); err != nil {
			return fmt.Errorf("bad flag format --%s with field %s", f.Name, err)
		}
	}
	return nil
}

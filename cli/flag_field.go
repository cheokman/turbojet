package cli

import "fmt"

type Field struct {
	Key string

	Required bool

	Repeatable bool

	DefaultValue string

	Short string

	assigned bool

	value  string
	values []string
}

func (f *Field) getValue() (string, bool) {
	if f.assigned {
		return f.value, true
	} else if f.DefaultValue != "" {
		return f.DefaultValue, false
	} else {
		return "", false
	}
}

func (f *Field) assign(v string) {
	f.assigned = true
	f.value = v
	f.values = append(f.values, v)
}

func (f *Field) check() error {
	if f.Required && !f.assigned {
		if f.Key != "" {
			return fmt.Errorf("%s= required", f.Key)
		}
		return fmt.Errorf("value required")

	}
	if !f.Repeatable && len(f.values) > 1 {
		if f.Key != "" {
			return fmt.Errorf("%s= duplicated", f.Key)
		}
		return fmt.Errorf("value duplicated")
	}
	return nil
}

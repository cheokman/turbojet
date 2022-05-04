package content

import "fmt"

type InvalidContentFileError struct {
	Name string
}

func NewInvalidContentFileError(name string) error {
	return &InvalidContentFileError{
		Name: name,
	}
}

func (e *InvalidContentFileError) Error() string {
	return fmt.Sprintf("'%s' is not a valid context file", e.Name)
}

type InvalidCacheFromError struct {
	Name string
}

func NewInvalidCacheFromError(name string) error {
	return &InvalidCacheFromError{
		Name: name,
	}
}

func (e *InvalidCacheFromError) Error() string {
	return fmt.Sprintf("'%s' is not a valid folder", e.Name)
}

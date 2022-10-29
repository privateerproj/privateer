package run

import "fmt"

// RaidError retains an error object and the name of the pack that generated it
type RaidError struct {
	Raid string
	Err  error
}

// RaidErrors holds a list of errors and an Error() method
// so it adheres to the standard Error interface
type RaidErrors struct {
	Errors []RaidError
}

func (e *RaidErrors) Error() string {
	return fmt.Sprintf("Service Pack Errors: %v", e.Errors)
}

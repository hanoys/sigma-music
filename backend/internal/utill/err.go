package utill

import "fmt"

func WrapError(e1, e2 error) error {
	return fmt.Errorf("%w: %w", e1, e2)
}

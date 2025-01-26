package document

import (
	"fmt"
	"log"
)

// appends to a slice, but guaranties to return a new one and not alter the original.
func appendNew[T any](sl []T, elems ...T) []T {
	sl2 := make([]T, len(sl), len(sl)+len(elems))
	copy(sl2, sl)
	sl2 = append(sl2, elems...)
	return sl2
}

func warnOrError(strict bool, pattern string, args ...any) error {
	if strict {
		return fmt.Errorf(pattern, args...)
	}
	log.Printf("WARNING: "+pattern+"\n", args...)
	return nil
}

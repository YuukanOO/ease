package flag

// Check if the given flag is set in the given value.
func IsSet[T ~uint](value T, flag T) bool {
	return value&flag != 0
}

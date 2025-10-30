package ptr

// Addr returns the address of the given value.
func Addr[T any](t T) *T {
	return &t
}

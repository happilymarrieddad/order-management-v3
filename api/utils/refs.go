package utils

func Ref[T any](v T) *T {
	return &v
}

func Deref[T any](v *T) T {
	if v == nil {
		v = new(T)
	}
	return *v
}

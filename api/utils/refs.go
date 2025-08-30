package utils

func Ref[T any](v T) *T {
	return &v
}

func Deref[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}

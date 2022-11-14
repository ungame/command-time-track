package pointer

func New[P any](v P) *P {
	return &v
}

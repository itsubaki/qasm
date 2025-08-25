package visitor

type Option func(*Visitor)

func WithMaxQubits(n int) Option {
	return func(v *Visitor) {
		v.maxQubits = n
	}
}

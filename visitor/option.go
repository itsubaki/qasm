package visitor

// Option is a function that modifies the Visitor.
type Option func(*Visitor)

// WithMaxQubits sets the maximum number of qubits for the visitor.
func WithMaxQubits(n int) Option {
	return func(v *Visitor) {
		v.maxQubits = n
	}
}

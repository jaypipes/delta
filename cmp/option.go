package cmp

// Options contains instructions to the Comparable and Compare functions
// for how to compare subjects.
type Options struct {
	// NoResolveTypeDef turns off behaviour that resolves Go type
	// definitions to their underlying type.
	NoResolveTypeDef bool
	// NoCast turns off behaviour that casts different underlying primitive
	// types to the larget primitive size of the same type before comparing
	// values.
	NoCast bool
	// NoConvert turns off behaviour that converts different underlying
	// primitive types to the same primitive type before comparing values.
	NoConvert bool
}

// Option allows the user to control the behaviour of the Comparable and
// Compare functions.
type Option func(*Options)

// WithNoResolveTypeDef disables resolution of Go type defintions to their
// underlying type.
func WithNoResolveTypeDef() Option {
	return func(o *Options) {
		o.NoResolveTypeDef = true
	}
}

var NoResolveTypeDef = WithNoResolveTypeDef()

// WithNoCast disables conversion of compared subjects to a convertible Go
// type.
func WithNoCast() Option {
	return func(o *Options) {
		o.NoCast = true
	}
}

var NoCast = WithNoCast()

// WithNoConvert disables conversion of compared subjects to a convertible Go
// type.
func WithNoConvert() Option {
	return func(o *Options) {
		o.NoConvert = true
	}
}

var NoConvert = WithNoConvert()

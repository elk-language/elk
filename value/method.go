package value

// Elk Method object
type Method interface {
	// Name of the method
	Name() Symbol
	// Names of all parameters in order
	Parameters() []Symbol
	// The number of parameters in total
	ParameterCount() int
	// The number of optional parameters with default values
	OptionalParameterCount() int
	// The number of parameters that appear after a rest parameter.
	//
	// -1 signals that there is no rest parameter
	//
	// 0 means that there are no more parameters after the rest param
	PostRestParameterCount() int
	// Whether the named rest parameter is present
	NamedRestParameter() bool
}

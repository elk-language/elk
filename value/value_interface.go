package value

// Go interface that represents an Elk value
type ValueInterface interface {
	ToValue() Value
	Class() *Class                         // Return the class of the value
	DirectClass() *Class                   // Return the direct class of this value that will be searched for methods first
	SingletonClass() *Class                // Return the singleton class of this value that holds methods unique to this object
	InstanceVariables() *InstanceVariables // Returns a pointer to the slice of instance vars of this value, nil if value doesn't support instance vars
	Inspect() string                       // Returns the string representation of the value
	Error() string                         // Implements the error interface
}

type ComparableValueInterface interface {
	ValueInterface
	comparable
}

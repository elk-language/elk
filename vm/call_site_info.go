package vm

import (
	"fmt"

	"github.com/elk-language/elk/value"
)

// Contains details like the number of arguments
// or the method name of a particular call site.
type CallSiteInfo struct {
	Name          value.Symbol
	ArgumentCount int
	Cache         [3]CallCacheEntry
}

type CallCache struct {
	Entries [3]CallCacheEntry
}

func LookupMethodInCache(class *value.Class, name value.Symbol, cacheLoc **CallCache) value.Method {
	cache := *cacheLoc
	for i := range len(cache.Entries) {
		cacheEntry := cache.Entries[i]
		if cacheEntry.Class == class {
			return cacheEntry.Method
		}
		if cacheEntry.Class == nil {
			method := class.LookupMethod(name)
			newEntries := cache.Entries
			newEntries[i] = CallCacheEntry{
				Class:  class,
				Method: method,
			}
			*cacheLoc = &CallCache{
				Entries: newEntries,
			}
			return method
		}
	}

	return class.LookupMethod(name)
}

type CallCacheEntry struct {
	Class  *value.Class
	Method value.Method
}

// Create a new CallSiteInfo.
func NewCallSiteInfo(methodName value.Symbol, argCount int) *CallSiteInfo {
	return &CallSiteInfo{
		Name:          methodName,
		ArgumentCount: argCount,
	}
}

func (*CallSiteInfo) Class() *value.Class {
	return nil
}

func (*CallSiteInfo) DirectClass() *value.Class {
	return nil
}

func (*CallSiteInfo) SingletonClass() *value.Class {
	return nil
}

func (*CallSiteInfo) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (c *CallSiteInfo) Copy() value.Reference {
	return &CallSiteInfo{
		Name:          c.Name,
		ArgumentCount: c.ArgumentCount,
	}
}

func (c *CallSiteInfo) ToValue() value.Value {
	return value.Ref(c)
}

func (c *CallSiteInfo) Inspect() string {
	return fmt.Sprintf(
		"CallSiteInfo{&: %p, name: %s, argument_count: %d}",
		c,
		c.Name.Inspect(),
		c.ArgumentCount,
	)
}

func (c *CallSiteInfo) Error() string {
	return c.Inspect()
}

// Contains details like the number of arguments
// or the pointer to a bytecode method of a particular call site.
type BytecodeCallSiteInfo struct {
	Method        *BytecodeFunction
	ArgumentCount int
	TailCall      bool
}

// Create a new BytecodeCallSiteInfo.
func NewBytecodeCallSiteInfo(method *BytecodeFunction, argCount int, tailCall bool) *BytecodeCallSiteInfo {
	return &BytecodeCallSiteInfo{
		Method:        method,
		ArgumentCount: argCount,
		TailCall:      tailCall,
	}
}

func (ci *BytecodeCallSiteInfo) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*BytecodeCallSiteInfo)
	if !ok {
		return false
	}

	return ci.Method == o.Method &&
		ci.ArgumentCount == o.ArgumentCount &&
		ci.TailCall == o.TailCall
}

func (*BytecodeCallSiteInfo) Class() *value.Class {
	return nil
}

func (*BytecodeCallSiteInfo) DirectClass() *value.Class {
	return nil
}

func (*BytecodeCallSiteInfo) SingletonClass() *value.Class {
	return nil
}

func (*BytecodeCallSiteInfo) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (c *BytecodeCallSiteInfo) Copy() value.Reference {
	return &BytecodeCallSiteInfo{
		Method:        c.Method,
		ArgumentCount: c.ArgumentCount,
	}
}

func (c *BytecodeCallSiteInfo) ToValue() value.Value {
	return value.Ref(c)
}

func (c *BytecodeCallSiteInfo) Inspect() string {
	return fmt.Sprintf(
		"BytecodeCallSiteInfo{&: %p, method: %s, argument_count: %d}",
		c,
		c.Method.Inspect(),
		c.ArgumentCount,
	)
}

func (c *BytecodeCallSiteInfo) Error() string {
	return c.Inspect()
}

// Contains details like the number of arguments
// or the pointer to a native method of a particular call site.
type NativeCallSiteInfo struct {
	Method        *NativeMethod
	ArgumentCount int
}

// Create a new NativeCallSiteInfo.
func NewNativeCallSiteInfo(method *NativeMethod, argCount int) *NativeCallSiteInfo {
	return &NativeCallSiteInfo{
		Method:        method,
		ArgumentCount: argCount,
	}
}

func (ci *NativeCallSiteInfo) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*NativeCallSiteInfo)
	if !ok {
		return false
	}

	return ci.Method == o.Method &&
		ci.ArgumentCount == o.ArgumentCount
}

func (*NativeCallSiteInfo) Class() *value.Class {
	return nil
}

func (*NativeCallSiteInfo) DirectClass() *value.Class {
	return nil
}

func (*NativeCallSiteInfo) SingletonClass() *value.Class {
	return nil
}

func (*NativeCallSiteInfo) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (c *NativeCallSiteInfo) Copy() value.Reference {
	return &NativeCallSiteInfo{
		Method:        c.Method,
		ArgumentCount: c.ArgumentCount,
	}
}

func (c *NativeCallSiteInfo) ToValue() value.Value {
	return value.Ref(c)
}

func (c *NativeCallSiteInfo) Inspect() string {
	return fmt.Sprintf(
		"NativeCallSiteInfo{&: %p, method: %s, argument_count: %d}",
		c,
		c.Method.Inspect(),
		c.ArgumentCount,
	)
}

func (c *NativeCallSiteInfo) Error() string {
	return c.Inspect()
}

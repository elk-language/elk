package value

import (
	"fmt"
)

// Contains details like the number of arguments
// or the pointer to a method of a particular call site.
type OptimisedCallSiteInfo struct {
	Method        *Method
	ArgumentCount int
}

// Create a new OptimisedCallSiteInfo.
func NewOptimisedCallSiteInfo(method *Method, argCount int) *OptimisedCallSiteInfo {
	return &OptimisedCallSiteInfo{
		Method:        method,
		ArgumentCount: argCount,
	}
}

func (*OptimisedCallSiteInfo) Class() *Class {
	return nil
}

func (*OptimisedCallSiteInfo) DirectClass() *Class {
	return nil
}

func (*OptimisedCallSiteInfo) SingletonClass() *Class {
	return nil
}

func (*OptimisedCallSiteInfo) InstanceVariables() *InstanceVariables {
	return nil
}

func (c *OptimisedCallSiteInfo) Copy() Reference {
	return &OptimisedCallSiteInfo{
		Method:        c.Method,
		ArgumentCount: c.ArgumentCount,
	}
}

func (c *OptimisedCallSiteInfo) ToValue() Value {
	return Ref(c)
}

func (c *OptimisedCallSiteInfo) Inspect() string {
	return fmt.Sprintf(
		"OptimisedCallSiteInfo{&: %p, method: %s, argument_count: %d}",
		c,
		(*c.Method).Inspect(),
		c.ArgumentCount,
	)
}

func (c *OptimisedCallSiteInfo) Error() string {
	return c.Inspect()
}

// Contains details like the number of arguments
// or the method name of a particular call site.
type CallSiteInfo struct {
	Name          Symbol
	ArgumentCount int
	Cache         [3]CallCacheEntry
}

type CallCache struct {
	Entries [3]CallCacheEntry
}

func LookupMethodInCache(class *Class, name Symbol, cacheLoc **CallCache) Method {
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
	Class  *Class
	Method Method
}

// Create a new CallSiteInfo.
func NewCallSiteInfo(methodName Symbol, argCount int) *CallSiteInfo {
	return &CallSiteInfo{
		Name:          methodName,
		ArgumentCount: argCount,
	}
}

func (*CallSiteInfo) Class() *Class {
	return nil
}

func (*CallSiteInfo) DirectClass() *Class {
	return nil
}

func (*CallSiteInfo) SingletonClass() *Class {
	return nil
}

func (*CallSiteInfo) InstanceVariables() *InstanceVariables {
	return nil
}

func (c *CallSiteInfo) Copy() Reference {
	return &CallSiteInfo{
		Name:          c.Name,
		ArgumentCount: c.ArgumentCount,
	}
}

func (c *CallSiteInfo) ToValue() Value {
	return Ref(c)
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

package value

import (
	"fmt"
)

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
	return c
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

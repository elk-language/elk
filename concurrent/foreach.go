package concurrent

// Execute a function concurrently for each element of a collection.
// Specify the max number of goroutines running at the same time.
func Foreach[E any](concurrencyLimit int, collection []E, f func(E)) {
	sem := make(chan bool, concurrencyLimit)
	for _, element := range collection {
		sem <- true
		go func(element E) {
			f(element)
			<-sem
		}(element)
	}
	for range cap(sem) {
		sem <- true
	}
}

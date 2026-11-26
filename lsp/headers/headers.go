package headers

import (
	"bytes"
	"fmt"

	"github.com/elk-language/elk/ds"
)

type Headers = ds.OrderedMap[string, string]

func New(pairs ...ds.Pair[string, string]) *Headers {
	return ds.NewOrderedMapWithPairs(pairs...)
}

func Encode(headers *Headers) []byte {
	var buf bytes.Buffer

	for key, val := range headers.All() {
		fmt.Fprintf(&buf, "%s: %s\r\n", key, val)
	}

	buf.WriteString("\r\n")
	return buf.Bytes()
}

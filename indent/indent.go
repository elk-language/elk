package indent

import (
	"bufio"
	"io"
	"strings"
)

const indentUnitString = "  "

var indentUnitBytes = []byte(indentUnitString)
var newlineBytes = []byte{'\n'}

func IndentString(out io.Writer, str string, indentLevel int) {
	scanner := bufio.NewScanner(strings.NewReader(str))

	firstIteration := true
	for scanner.Scan() {
		if !firstIteration {
			out.Write(newlineBytes)
		} else {
			firstIteration = false
		}

		for range indentLevel {
			out.Write(indentUnitBytes)
		}
		line := scanner.Bytes()
		out.Write(line)
	}
}

func IndentStringFromSecondLine(out io.Writer, str string, indentLevel int) {
	scanner := bufio.NewScanner(strings.NewReader(str))

	firstIteration := true
	for scanner.Scan() {
		if !firstIteration {
			out.Write(newlineBytes)
			for range indentLevel {
				out.Write(indentUnitBytes)
			}
		} else {
			firstIteration = false
		}

		line := scanner.Bytes()
		out.Write(line)
	}
}

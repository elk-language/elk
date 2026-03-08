package colorize

import (
	"bytes"
	"go/scanner"
	"go/token"

	"github.com/fatih/color"
)

// Colorize returns the Go source code with ANSI color escape codes
func Colorize(srcBytes []byte) []byte {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(srcBytes))

	var s scanner.Scanner
	s.Init(file, srcBytes, nil, scanner.ScanComments)

	kwColor := color.New(color.FgGreen)
	opColor := color.New(color.FgHiMagenta)
	strColor := color.New(color.FgHiYellow)
	intColor := color.New(color.FgHiBlue)
	floatColor := color.New(color.FgHiMagenta)
	comColor := color.New(color.FgHiBlack)

	var buf bytes.Buffer
	currentOffset := 0

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		start := file.Offset(pos)
		if len(lit) == 0 {
			lit = tok.String()
		}
		end := start + len(lit)

		// Emit any whitespace/gap before this token
		if start > currentOffset {
			buf.Write(srcBytes[currentOffset:start])
		}

		currentOffset = end

		// Emit the token with the appropriate color
		switch {
		case tok.IsKeyword():
			buf.WriteString(kwColor.Sprint(lit))
		case tok.IsOperator():
			buf.WriteString(opColor.Sprint(lit))
		case tok == token.STRING || tok == token.CHAR:
			buf.WriteString(strColor.Sprint(lit))
		case tok == token.INT:
			buf.WriteString(intColor.Sprint(lit))
		case tok == token.FLOAT || tok == token.IMAG:
			buf.WriteString(floatColor.Sprint(lit))
		case tok == token.COMMENT:
			buf.WriteString(comColor.Sprint(lit))
		default:
			buf.WriteString(lit)
		}
	}

	// Emit any remaining source (e.g. after last token or on scanner error)
	if currentOffset < len(srcBytes) {
		buf.Write(srcBytes[currentOffset:])
	}

	return buf.Bytes()
}

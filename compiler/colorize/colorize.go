package colorize

import (
	"bytes"
	"go/scanner"
	"go/token"

	"github.com/fatih/color"
)

// Colorize returns the Go source code with ANSI color escape codes
func Colorize(srcBytes []byte) []byte {
	return ColorizeWhen(srcBytes, !color.NoColor)
}

func ColorizeWhen(srcBytes []byte, useColor bool) []byte {
	if !useColor {
		return srcBytes
	}

	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(srcBytes))

	var s scanner.Scanner
	s.Init(file, srcBytes, nil, scanner.ScanComments)

	kwColor := color.New(color.FgGreen)
	kwColor.EnableColor()

	opColor := color.New(color.FgHiMagenta)
	opColor.EnableColor()

	strColor := color.New(color.FgHiYellow)
	strColor.EnableColor()

	intColor := color.New(color.FgHiBlue)
	intColor.EnableColor()

	floatColor := color.New(color.FgHiMagenta)
	floatColor.EnableColor()

	comColor := color.New(color.FgHiBlack)
	comColor.EnableColor()

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

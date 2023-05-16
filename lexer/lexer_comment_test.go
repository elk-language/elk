package lexer

import (
	"testing"

	"github.com/elk-language/elk/token"
)

func TestSingleLineComment(t *testing.T) {
	tests := testTable{
		"discards characters until a new line is reached": {
			input: `3 + # 25 / 3
							5`,
			want: []*token.Token{
				V(P(0, 1, 1, 1), token.DEC_INT, "3"),
				T(P(2, 1, 1, 3), token.PLUS),
				T(P(12, 1, 1, 13), token.NEWLINE),
				V(P(20, 1, 2, 8), token.DEC_INT, "5"),
			},
		},
		"can appear at the beginning of the line": {
			input: `# something awesome
							foo := 3`,
			want: []*token.Token{
				T(P(19, 1, 1, 20), token.NEWLINE),
				V(P(27, 3, 2, 8), token.PUBLIC_IDENTIFIER, "foo"),
				T(P(31, 2, 2, 12), token.COLON_EQUAL),
				V(P(34, 1, 2, 15), token.DEC_INT, "3"),
			},
		},
		"can appear on consecutive lines": {
			input: `
# peace
# and love
# from Elk
println 'Hey'
			`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.NEWLINE),
				T(P(8, 1, 2, 8), token.NEWLINE),
				T(P(19, 1, 3, 11), token.NEWLINE),
				T(P(30, 1, 4, 11), token.NEWLINE),
				V(P(31, 7, 5, 1), token.PUBLIC_IDENTIFIER, "println"),
				V(P(39, 5, 5, 9), token.RAW_STRING, "Hey"),
				T(P(44, 1, 5, 14), token.NEWLINE),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestBlockComment(t *testing.T) {
	tests := testTable{
		"discards characters in the middle of the line": {
			input: `3 + #[25 / 3]# 5`,
			want: []*token.Token{
				V(P(0, 1, 1, 1), token.DEC_INT, "3"),
				T(P(2, 1, 1, 3), token.PLUS),
				V(P(15, 1, 1, 16), token.DEC_INT, "5"),
			},
		},
		"must be terminated": {
			input: `3 + #[25 / 3 5`,
			want: []*token.Token{
				V(P(0, 1, 1, 1), token.DEC_INT, "3"),
				T(P(2, 1, 1, 3), token.PLUS),
				V(P(4, 10, 1, 5), token.ERROR, "unbalanced block comments, expected 1 more block comment ending(s) `]#`"),
			},
		},
		"discards multiple lines": {
			input: `
class String
	#[
		def length: Integer
			len := 0
			self.each -> len += 1
			len
		end
	]#
end
			`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.NEWLINE),
				T(P(1, 5, 2, 1), token.CLASS),
				V(P(7, 6, 2, 7), token.PUBLIC_CONSTANT, "String"),
				T(P(13, 1, 2, 13), token.NEWLINE),
				T(P(93, 1, 9, 4), token.NEWLINE),
				T(P(94, 3, 10, 1), token.END),
				T(P(97, 1, 10, 4), token.NEWLINE),
			},
		},
		"may be nested": {
			input: `
class String
	#[
		def length: Integer
			len := 0
			self.each ->
				len += 1
				#[
					#[ another comment ]#
					println len
				]#
			end
			len
		end
	]#
end
			`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.NEWLINE),
				T(P(1, 5, 2, 1), token.CLASS),
				V(P(7, 6, 2, 7), token.PUBLIC_CONSTANT, "String"),
				T(P(13, 1, 2, 13), token.NEWLINE),
				T(P(162, 1, 15, 4), token.NEWLINE),
				T(P(163, 3, 16, 1), token.END),
				T(P(166, 1, 16, 4), token.NEWLINE),
			},
		},
		"nesting must be balanced": {
			input: `
class String
	#[
		def length: Integer
			len := 0
			self.each ->
				len += 1
				#[
					#[ another comment
					println len
			end
			len
		end
	]#
end
			`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.NEWLINE),
				T(P(1, 5, 2, 1), token.CLASS),
				V(P(7, 6, 2, 7), token.PUBLIC_CONSTANT, "String"),
				T(P(13, 1, 2, 13), token.NEWLINE),
				V(P(15, 145, 3, 2), token.ERROR, "unbalanced block comments, expected 2 more block comment ending(s) `]#`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

func TestDocComment(t *testing.T) {
	tests := testTable{
		"may be used in the middle of the line": {
			input: `3 + ##[25 / 3]## 5`,
			want: []*token.Token{
				V(P(0, 1, 1, 1), token.DEC_INT, "3"),
				T(P(2, 1, 1, 3), token.PLUS),
				V(P(4, 12, 1, 5), token.DOC_COMMENT, "25 / 3"),
				V(P(17, 1, 1, 18), token.DEC_INT, "5"),
			},
		},
		"must be terminated": {
			input: `3 + ##[25 / 3 5`,
			want: []*token.Token{
				V(P(0, 1, 1, 1), token.DEC_INT, "3"),
				T(P(2, 1, 1, 3), token.PLUS),
				V(P(4, 11, 1, 5), token.ERROR, "unbalanced doc comments, expected 1 more doc comment ending(s) `]##`"),
			},
		},
		"may contain multiple lines": {
			input: `
class String
	##[
		def length: Integer
			len := 0
			self.each -> len += 1
			len
		end
	]##
end
			`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.NEWLINE),
				T(P(1, 5, 2, 1), token.CLASS),
				V(P(7, 6, 2, 7), token.PUBLIC_CONSTANT, "String"),
				T(P(13, 1, 2, 13), token.NEWLINE),
				V(P(15, 80, 3, 2), token.DOC_COMMENT, `def length: Integer
	len := 0
	self.each -> len += 1
	len
end`),
				T(P(95, 1, 9, 5), token.NEWLINE),
				T(P(96, 3, 10, 1), token.END),
				T(P(99, 1, 10, 4), token.NEWLINE),
			},
		},
		"trims leading whitespace of each line up to the least indented line's level": {
			input: `##[
		Something
			awesome
				and
		foo
					bar
]##`,
			want: []*token.Token{
				V(P(0, 53, 1, 1), token.DOC_COMMENT, `Something
	awesome
		and
foo
			bar`),
			},
		},
		"trims leading and trailing whitespace when single line": {
			input: `##[   foo + bar = awesome          ]##`,
			want: []*token.Token{
				V(P(0, 38, 1, 1), token.DOC_COMMENT, `foo + bar = awesome`),
			},
		},
		"trims leading and trailing endlines": {
			input: `##[



			foo + bar = awesome


]##`,
			want: []*token.Token{
				V(P(0, 35, 1, 1), token.DOC_COMMENT, `foo + bar = awesome`),
			},
		},
		"may be nested": {
			input: `
class String
	##[
		def length: Integer
			len := 0
			self.each ->
				len += 1
				##[
					##[ another comment ]##
					println len
				]##
			end
			len
		end
	]##
end
			`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.NEWLINE),
				T(P(1, 5, 2, 1), token.CLASS),
				V(P(7, 6, 2, 7), token.PUBLIC_CONSTANT, "String"),
				T(P(13, 1, 2, 13), token.NEWLINE),
				V(P(15, 153, 3, 2), token.DOC_COMMENT, `def length: Integer
	len := 0
	self.each ->
		len += 1
		##[
			##[ another comment ]##
			println len
		]##
	end
	len
end`),
				T(P(168, 1, 15, 5), token.NEWLINE),
				T(P(169, 3, 16, 1), token.END),
				T(P(172, 1, 16, 4), token.NEWLINE),
			},
		},
		"nesting must be balanced": {
			input: `
class String
	##[
		def length: Integer
			len := 0
			self.each ->
				len += 1
				##[
					##[ another comment
					println len
			end
			len
		end
	]##
end
			`,
			want: []*token.Token{
				T(P(0, 1, 1, 1), token.NEWLINE),
				T(P(1, 5, 2, 1), token.CLASS),
				V(P(7, 6, 2, 7), token.PUBLIC_CONSTANT, "String"),
				T(P(13, 1, 2, 13), token.NEWLINE),
				V(P(15, 149, 3, 2), token.ERROR, "unbalanced doc comments, expected 2 more doc comment ending(s) `]##`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

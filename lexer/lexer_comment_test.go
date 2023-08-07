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
				V(S(P(0, 1, 1), P(0, 1, 1)), token.INT, "3"),
				T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
				T(S(P(12, 1, 13), P(12, 1, 13)), token.NEWLINE),
				V(S(P(20, 2, 8), P(20, 2, 8)), token.INT, "5"),
			},
		},
		"can appear at the beginning of the line": {
			input: `# something awesome
							foo := 3`,
			want: []*token.Token{
				T(S(P(19, 1, 20), P(19, 1, 20)), token.NEWLINE),
				V(S(P(27, 2, 8), P(29, 2, 10)), token.PUBLIC_IDENTIFIER, "foo"),
				T(S(P(31, 2, 12), P(32, 2, 13)), token.COLON_EQUAL),
				V(S(P(34, 2, 15), P(34, 2, 15)), token.INT, "3"),
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
				T(S(P(0, 1, 1), P(0, 1, 1)), token.NEWLINE),
				T(S(P(8, 2, 8), P(8, 2, 8)), token.NEWLINE),
				T(S(P(19, 3, 11), P(19, 3, 11)), token.NEWLINE),
				T(S(P(30, 4, 11), P(30, 4, 11)), token.NEWLINE),
				V(S(P(31, 5, 1), P(37, 5, 7)), token.PUBLIC_IDENTIFIER, "println"),
				V(S(P(39, 5, 9), P(43, 5, 13)), token.RAW_STRING, "Hey"),
				T(S(P(44, 5, 14), P(44, 5, 14)), token.NEWLINE),
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
				V(S(P(0, 1, 1), P(0, 1, 1)), token.INT, "3"),
				T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
				V(S(P(15, 1, 16), P(15, 1, 16)), token.INT, "5"),
			},
		},
		"must be terminated": {
			input: `3 + #[25 / 3 5`,
			want: []*token.Token{
				V(S(P(0, 1, 1), P(0, 1, 1)), token.INT, "3"),
				T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
				V(S(P(4, 1, 5), P(13, 1, 14)), token.ERROR, "unbalanced block comments, expected 1 more block comment ending(s) `]#`"),
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
				T(S(P(0, 1, 1), P(0, 1, 1)), token.NEWLINE),
				T(S(P(1, 2, 1), P(5, 2, 5)), token.CLASS),
				V(S(P(7, 2, 7), P(12, 2, 12)), token.PUBLIC_CONSTANT, "String"),
				T(S(P(13, 2, 13), P(13, 2, 13)), token.NEWLINE),
				T(S(P(93, 9, 4), P(93, 9, 4)), token.NEWLINE),
				T(S(P(94, 10, 1), P(96, 10, 3)), token.END),
				T(S(P(97, 10, 4), P(97, 10, 4)), token.NEWLINE),
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
				T(S(P(0, 1, 1), P(0, 1, 1)), token.NEWLINE),
				T(S(P(1, 2, 1), P(5, 2, 5)), token.CLASS),
				V(S(P(7, 2, 7), P(12, 2, 12)), token.PUBLIC_CONSTANT, "String"),
				T(S(P(13, 2, 13), P(13, 2, 13)), token.NEWLINE),
				T(S(P(162, 15, 4), P(162, 15, 4)), token.NEWLINE),
				T(S(P(163, 16, 1), P(165, 16, 3)), token.END),
				T(S(P(166, 16, 4), P(166, 16, 4)), token.NEWLINE),
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
				T(S(P(0, 1, 1), P(0, 1, 1)), token.NEWLINE),
				T(S(P(1, 2, 1), P(5, 2, 5)), token.CLASS),
				V(S(P(7, 2, 7), P(12, 2, 12)), token.PUBLIC_CONSTANT, "String"),
				T(S(P(13, 2, 13), P(13, 2, 13)), token.NEWLINE),
				V(S(P(15, 3, 2), P(159, 16, 3)), token.ERROR, "unbalanced block comments, expected 2 more block comment ending(s) `]#`"),
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
				V(S(P(0, 1, 1), P(0, 1, 1)), token.INT, "3"),
				T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
				V(S(P(4, 1, 5), P(15, 1, 16)), token.DOC_COMMENT, "25 / 3"),
				V(S(P(17, 1, 18), P(17, 1, 18)), token.INT, "5"),
			},
		},
		"must be terminated": {
			input: `3 + ##[25 / 3 5`,
			want: []*token.Token{
				V(S(P(0, 1, 1), P(0, 1, 1)), token.INT, "3"),
				T(S(P(2, 1, 3), P(2, 1, 3)), token.PLUS),
				V(S(P(4, 1, 5), P(14, 1, 15)), token.ERROR, "unbalanced doc comments, expected 1 more doc comment ending(s) `]##`"),
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
				T(S(P(0, 1, 1), P(0, 1, 1)), token.NEWLINE),
				T(S(P(1, 2, 1), P(5, 2, 5)), token.CLASS),
				V(S(P(7, 2, 7), P(12, 2, 12)), token.PUBLIC_CONSTANT, "String"),
				T(S(P(13, 2, 13), P(13, 2, 13)), token.NEWLINE),
				V(S(P(15, 3, 2), P(94, 9, 4)), token.DOC_COMMENT, `def length: Integer
	len := 0
	self.each -> len += 1
	len
end`),
				T(S(P(95, 9, 5), P(95, 9, 5)), token.NEWLINE),
				T(S(P(96, 10, 1), P(98, 10, 3)), token.END),
				T(S(P(99, 10, 4), P(99, 10, 4)), token.NEWLINE),
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
				V(S(P(0, 1, 1), P(52, 7, 3)), token.DOC_COMMENT, `Something
	awesome
		and
foo
			bar`),
			},
		},
		"trims leading and trailing whitespace when single line": {
			input: `##[   foo + bar = awesome          ]##`,
			want: []*token.Token{
				V(S(P(0, 1, 1), P(37, 1, 38)), token.DOC_COMMENT, `foo + bar = awesome`),
			},
		},
		"trims leading and trailing endlines": {
			input: `##[



			foo + bar = awesome


]##`,
			want: []*token.Token{
				V(S(P(0, 1, 1), P(34, 8, 3)), token.DOC_COMMENT, `foo + bar = awesome`),
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
				T(S(P(0, 1, 1), P(0, 1, 1)), token.NEWLINE),
				T(S(P(1, 2, 1), P(5, 2, 5)), token.CLASS),
				V(S(P(7, 2, 7), P(12, 2, 12)), token.PUBLIC_CONSTANT, "String"),
				T(S(P(13, 2, 13), P(13, 2, 13)), token.NEWLINE),
				V(S(P(15, 3, 2), P(167, 15, 4)), token.DOC_COMMENT, `def length: Integer
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
				T(S(P(168, 15, 5), P(168, 15, 5)), token.NEWLINE),
				T(S(P(169, 16, 1), P(171, 16, 3)), token.END),
				T(S(P(172, 16, 4), P(172, 16, 4)), token.NEWLINE),
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
				T(S(P(0, 1, 1), P(0, 1, 1)), token.NEWLINE),
				T(S(P(1, 2, 1), P(5, 2, 5)), token.CLASS),
				V(S(P(7, 2, 7), P(12, 2, 12)), token.PUBLIC_CONSTANT, "String"),
				T(S(P(13, 2, 13), P(13, 2, 13)), token.NEWLINE),
				V(S(P(15, 3, 2), P(163, 16, 3)), token.ERROR, "unbalanced doc comments, expected 2 more doc comment ending(s) `]##`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

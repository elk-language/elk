package lexer

import "testing"

func TestSingleLineComment(t *testing.T) {
	tests := testTable{
		"discards characters until a new line is reached": {
			input: `3 + # 25 / 3
							5`,
			want: []*Token{
				V(P(0, 1, 1, 1), DecIntToken, "3"),
				T(P(2, 1, 1, 3), PlusToken),
				T(P(12, 1, 1, 13), EndLineToken),
				V(P(20, 1, 2, 8), DecIntToken, "5"),
			},
		},
		"can appear at the beginning of the line": {
			input: `# something awesome
							foo := 3`,
			want: []*Token{
				T(P(19, 1, 1, 20), EndLineToken),
				V(P(27, 3, 2, 8), PublicIdentifierToken, "foo"),
				T(P(31, 2, 2, 12), ColonEqualToken),
				V(P(34, 1, 2, 15), DecIntToken, "3"),
			},
		},
		"can appear on consecutive lines": {
			input: `
# peace
# and love
# from Elk
println 'Hey'
			`,
			want: []*Token{
				T(P(0, 1, 1, 1), EndLineToken),
				T(P(8, 1, 2, 8), EndLineToken),
				T(P(19, 1, 3, 11), EndLineToken),
				T(P(30, 1, 4, 11), EndLineToken),
				V(P(31, 7, 5, 1), PublicIdentifierToken, "println"),
				V(P(39, 5, 5, 9), RawStringToken, "Hey"),
				T(P(44, 1, 5, 14), EndLineToken),
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
			want: []*Token{
				V(P(0, 1, 1, 1), DecIntToken, "3"),
				T(P(2, 1, 1, 3), PlusToken),
				V(P(15, 1, 1, 16), DecIntToken, "5"),
			},
		},
		"must be terminated": {
			input: `3 + #[25 / 3 5`,
			want: []*Token{
				V(P(0, 1, 1, 1), DecIntToken, "3"),
				T(P(2, 1, 1, 3), PlusToken),
				V(P(4, 10, 1, 5), ErrorToken, "unbalanced block comments, expected 1 more block comment ending(s) `]#`"),
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
			want: []*Token{
				T(P(0, 1, 1, 1), EndLineToken),
				T(P(1, 5, 2, 1), ClassToken),
				V(P(7, 6, 2, 7), PublicConstantToken, "String"),
				T(P(13, 1, 2, 13), EndLineToken),
				T(P(93, 1, 9, 4), EndLineToken),
				T(P(94, 3, 10, 1), EndToken),
				T(P(97, 1, 10, 4), EndLineToken),
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
			want: []*Token{
				T(P(0, 1, 1, 1), EndLineToken),
				T(P(1, 5, 2, 1), ClassToken),
				V(P(7, 6, 2, 7), PublicConstantToken, "String"),
				T(P(13, 1, 2, 13), EndLineToken),
				T(P(162, 1, 15, 4), EndLineToken),
				T(P(163, 3, 16, 1), EndToken),
				T(P(166, 1, 16, 4), EndLineToken),
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
			want: []*Token{
				T(P(0, 1, 1, 1), EndLineToken),
				T(P(1, 5, 2, 1), ClassToken),
				V(P(7, 6, 2, 7), PublicConstantToken, "String"),
				T(P(13, 1, 2, 13), EndLineToken),
				V(P(15, 145, 3, 2), ErrorToken, "unbalanced block comments, expected 2 more block comment ending(s) `]#`"),
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
			want: []*Token{
				V(P(0, 1, 1, 1), DecIntToken, "3"),
				T(P(2, 1, 1, 3), PlusToken),
				V(P(4, 12, 1, 5), DocCommentToken, "25 / 3"),
				V(P(17, 1, 1, 18), DecIntToken, "5"),
			},
		},
		"must be terminated": {
			input: `3 + ##[25 / 3 5`,
			want: []*Token{
				V(P(0, 1, 1, 1), DecIntToken, "3"),
				T(P(2, 1, 1, 3), PlusToken),
				V(P(4, 11, 1, 5), ErrorToken, "unbalanced doc comments, expected 1 more doc comment ending(s) `]##`"),
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
			want: []*Token{
				T(P(0, 1, 1, 1), EndLineToken),
				T(P(1, 5, 2, 1), ClassToken),
				V(P(7, 6, 2, 7), PublicConstantToken, "String"),
				T(P(13, 1, 2, 13), EndLineToken),
				V(P(15, 80, 3, 2), DocCommentToken, `def length: Integer
	len := 0
	self.each -> len += 1
	len
end`),
				T(P(95, 1, 9, 5), EndLineToken),
				T(P(96, 3, 10, 1), EndToken),
				T(P(99, 1, 10, 4), EndLineToken),
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
			want: []*Token{
				V(P(0, 53, 1, 1), DocCommentToken, `Something
	awesome
		and
foo
			bar`),
			},
		},
		"trims leading and trailing whitespace when single line": {
			input: `##[   foo + bar = awesome          ]##`,
			want: []*Token{
				V(P(0, 38, 1, 1), DocCommentToken, `foo + bar = awesome`),
			},
		},
		"trims leading and trailing endlines": {
			input: `##[



			foo + bar = awesome


]##`,
			want: []*Token{
				V(P(0, 35, 1, 1), DocCommentToken, `foo + bar = awesome`),
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
			want: []*Token{
				T(P(0, 1, 1, 1), EndLineToken),
				T(P(1, 5, 2, 1), ClassToken),
				V(P(7, 6, 2, 7), PublicConstantToken, "String"),
				T(P(13, 1, 2, 13), EndLineToken),
				V(P(15, 153, 3, 2), DocCommentToken, `def length: Integer
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
				T(P(168, 1, 15, 5), EndLineToken),
				T(P(169, 3, 16, 1), EndToken),
				T(P(172, 1, 16, 4), EndLineToken),
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
			want: []*Token{
				T(P(0, 1, 1, 1), EndLineToken),
				T(P(1, 5, 2, 1), ClassToken),
				V(P(7, 6, 2, 7), PublicConstantToken, "String"),
				T(P(13, 1, 2, 13), EndLineToken),
				V(P(15, 149, 3, 2), ErrorToken, "unbalanced doc comments, expected 2 more doc comment ending(s) `]##`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

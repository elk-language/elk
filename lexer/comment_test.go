package lexer

import "testing"

func TestSingleLineComment(t *testing.T) {
	tests := testTable{
		"discards characters until a new line is reached": {
			input: `3 + # 25 / 3
							5`,
			want: []*Token{
				V(DecIntToken, "3", 0, 1, 1, 1),
				T(PlusToken, 2, 1, 1, 3),
				T(EndLineToken, 12, 1, 1, 13),
				V(DecIntToken, "5", 20, 1, 2, 8),
			},
		},
		"can appear at the beginning of the line": {
			input: `# something awesome
							foo := 3`,
			want: []*Token{
				T(EndLineToken, 19, 1, 1, 20),
				V(PublicIdentifierToken, "foo", 27, 3, 2, 8),
				T(ColonEqualToken, 31, 2, 2, 12),
				V(DecIntToken, "3", 34, 1, 2, 15),
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
				T(EndLineToken, 0, 1, 1, 1),
				T(EndLineToken, 8, 1, 2, 8),
				T(EndLineToken, 19, 1, 3, 11),
				T(EndLineToken, 30, 1, 4, 11),
				V(PublicIdentifierToken, "println", 31, 7, 5, 1),
				V(RawStringToken, "Hey", 39, 5, 5, 9),
				T(EndLineToken, 44, 1, 5, 14),
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
				V(DecIntToken, "3", 0, 1, 1, 1),
				T(PlusToken, 2, 1, 1, 3),
				V(DecIntToken, "5", 15, 1, 1, 16),
			},
		},
		"must be terminated": {
			input: `3 + #[25 / 3 5`,
			want: []*Token{
				V(DecIntToken, "3", 0, 1, 1, 1),
				T(PlusToken, 2, 1, 1, 3),
				V(ErrorToken, "unbalanced block comments, expected 1 more block comment ending(s) `]#`", 4, 10, 1, 5),
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
				T(EndLineToken, 0, 1, 1, 1),
				T(ClassToken, 1, 5, 2, 1),
				V(PublicConstantToken, "String", 7, 6, 2, 7),
				T(EndLineToken, 13, 1, 2, 13),
				T(EndLineToken, 93, 1, 9, 4),
				T(EndToken, 94, 3, 10, 1),
				T(EndLineToken, 97, 1, 10, 4),
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
				T(EndLineToken, 0, 1, 1, 1),
				T(ClassToken, 1, 5, 2, 1),
				V(PublicConstantToken, "String", 7, 6, 2, 7),
				T(EndLineToken, 13, 1, 2, 13),
				T(EndLineToken, 162, 1, 15, 4),
				T(EndToken, 163, 3, 16, 1),
				T(EndLineToken, 166, 1, 16, 4),
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
				T(EndLineToken, 0, 1, 1, 1),
				T(ClassToken, 1, 5, 2, 1),
				V(PublicConstantToken, "String", 7, 6, 2, 7),
				T(EndLineToken, 13, 1, 2, 13),
				V(ErrorToken, "unbalanced block comments, expected 2 more block comment ending(s) `]#`", 15, 145, 3, 2),
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
				V(DecIntToken, "3", 0, 1, 1, 1),
				T(PlusToken, 2, 1, 1, 3),
				V(DocCommentToken, "25 / 3", 4, 12, 1, 5),
				V(DecIntToken, "5", 17, 1, 1, 18),
			},
		},
		"must be terminated": {
			input: `3 + ##[25 / 3 5`,
			want: []*Token{
				V(DecIntToken, "3", 0, 1, 1, 1),
				T(PlusToken, 2, 1, 1, 3),
				V(ErrorToken, "unbalanced doc comments, expected 1 more doc comment ending(s) `]##`", 4, 11, 1, 5),
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
				T(EndLineToken, 0, 1, 1, 1),
				T(ClassToken, 1, 5, 2, 1),
				V(PublicConstantToken, "String", 7, 6, 2, 7),
				T(EndLineToken, 13, 1, 2, 13),
				V(DocCommentToken, `def length: Integer
	len := 0
	self.each -> len += 1
	len
end`, 15, 80, 3, 2),
				T(EndLineToken, 95, 1, 9, 5),
				T(EndToken, 96, 3, 10, 1),
				T(EndLineToken, 99, 1, 10, 4),
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
				V(DocCommentToken, `Something
	awesome
		and
foo
			bar`, 0, 53, 1, 1),
			},
		},
		"trims leading and trailing whitespace when single line": {
			input: `##[   foo + bar = awesome          ]##`,
			want: []*Token{
				V(DocCommentToken, `foo + bar = awesome`, 0, 38, 1, 1),
			},
		},
		"trims leading and trailing endlines": {
			input: `##[



			foo + bar = awesome


]##`,
			want: []*Token{
				V(DocCommentToken, `foo + bar = awesome`, 0, 35, 1, 1),
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
				T(EndLineToken, 0, 1, 1, 1),
				T(ClassToken, 1, 5, 2, 1),
				V(PublicConstantToken, "String", 7, 6, 2, 7),
				T(EndLineToken, 13, 1, 2, 13),
				V(DocCommentToken, `def length: Integer
	len := 0
	self.each ->
		len += 1
		##[
			##[ another comment ]##
			println len
		]##
	end
	len
end`, 15, 153, 3, 2),
				T(EndLineToken, 168, 1, 15, 5),
				T(EndToken, 169, 3, 16, 1),
				T(EndLineToken, 172, 1, 16, 4),
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
				T(EndLineToken, 0, 1, 1, 1),
				T(ClassToken, 1, 5, 2, 1),
				V(PublicConstantToken, "String", 7, 6, 2, 7),
				T(EndLineToken, 13, 1, 2, 13),
				V(ErrorToken, "unbalanced doc comments, expected 2 more doc comment ending(s) `]##`", 15, 149, 3, 2),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}

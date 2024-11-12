package position

// Represents something that contains a span.
type SpanInterface interface {
	Span() *Span
	SetSpan(*Span)
}

// Represents a span of text in a string/file.
type Span struct {
	StartPos *Position
	EndPos   *Position
}

// Create a new Span.
func NewSpan(startPos, endPos *Position) *Span {
	return &Span{
		StartPos: startPos,
		EndPos:   endPos,
	}
}

var DefaultSpan = NewSpanFromPosition(Default)

// Create a new span from a single position.
func NewSpanFromPosition(pos *Position) *Span {
	return &Span{
		StartPos: pos,
		EndPos:   pos,
	}
}

func (s *Span) Span() *Span {
	return s
}

func (s *Span) SetSpan(other *Span) {
	s.StartPos = other.StartPos
	s.EndPos = other.EndPos
}

// Join two spans into one.
// Works properly when the receiver is nil or the argument is nil.
func (left *Span) Join(right *Span) *Span {
	if right == nil {
		return left
	}
	if left == nil {
		return right
	}

	return &Span{
		StartPos: left.StartPos,
		EndPos:   right.EndPos,
	}
}

// Joins the given position with the last element of the given collection.
func JoinSpanOfLastElement[Element SpanInterface](left *Span, rightCollection []Element) *Span {
	if len(rightCollection) > 0 {
		return left.Join(rightCollection[len(rightCollection)-1].Span())
	}

	return left
}

// Join the position of the first element of a collection with the last one.
func JoinSpanOfCollection[Element SpanInterface](collection []Element) *Span {
	if len(collection) < 1 {
		return nil
	}

	left := collection[0].Span()
	if len(collection) == 1 {
		return left
	}

	return JoinSpanOfLastElement(left, collection)
}

// Retrieve the position of the last element of a collection.
func SpanOfLastElement[Element SpanInterface](collection []Element) *Span {
	if len(collection) > 0 {
		return collection[len(collection)-1].Span()
	}

	return nil
}

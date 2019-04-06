package span

import (
	engospan "github.com/EngoEngine/math/span"
)

// Span represents an interval.
type Span struct {
	Min, Max float32
}

func toIntern(x engospan.Span) Span {
	return Span{Min: x.Min, Max: x.Max}
}

func toExtern(i Span) engospan.Span {
	return engospan.Span{Min: i.Min, Max: i.Max}
}

// Add 2 span togheter
//	[a, b] + [c, d] = [a+c, b+d]
func (s0 Span) Add(s1 Span) Span {
	return toIntern(toExtern(s0).Add(toExtern(s1)))
}

// Sub 2 span togheter
//	[a, b] - [c, d] = [a-c, b-d]
func (s0 Span) Sub(s1 Span) Span {
	return toIntern(toExtern(s0).Sub(toExtern(s1)))
}

// Mul multiply this these 2 span togheter
//	[a, b] * [c, d] = [min(ac, ad, bc, bd), max(ac, ad, bc, bd)]
func (s0 Span) Mul(s1 Span) Span {
	return toIntern(toExtern(s0).Mul(toExtern(s1)))
}

// Div returns s0/s1
func (s0 Span) Div(s1 Span) Span {
	return toIntern(toExtern(s0).Div(toExtern(s1)))
}

// Abs return the absolute of the given span.
func Abs(s Span) Span {
	return toIntern(engospan.Abs(toExtern(s)))
}

package spells

import "github.com/ireina7/summoner"

func Summon[A any]() (A, error) {
	return summoner.Summon[A]()
}

func Given[A any](a A) error {
	return summoner.Given(a)
}

type Display[A any] interface {
	Display(a A) string
}

type Debug[A any] interface {
	Debug(a A) string
}

type Default[A any] interface {
	Default() A
}

type Monoid[A any] interface {
	Zero() A
	Plus(A, A) A
}

type Fold[A any] struct {
	Monoid Monoid[A]
}

func (self *Fold[A]) FoldLeft(xs []A) A {
	m := self.Monoid
	ans := m.Zero()
	for _, x := range xs {
		ans = m.Plus(ans, x)
	}
	return ans
}

type FunctorSlice[A, B any] struct{}

func (self *FunctorSlice[A, B]) Map(xs []A, f func(A) B) []B {
	ys := []B{}
	for _, x := range xs {
		ys = append(ys, f(x))
	}
	return ys
}

func init() {
	// summoner.Given[int](0)
}

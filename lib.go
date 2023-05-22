package summoner

import (
	"reflect"

	"github.com/ireina7/summoner/summoner"
)

func TypeOf[A any]() reflect.Type {
	return summoner.TypeOf[A]()
}

func IsRule[A any]() bool {
	return summoner.IsRule[A]()
}

func Summon[I any]() (I, error) {
	return summoner.Summon[I]()
}

func SummonType(t reflect.Type) (any, error) {
	return summoner.SummonType(t)
}

func Given[I any](instance I) error {
	return summoner.Given(instance)
}

func GivenType(instance any, t reflect.Type) error {
	return summoner.GivenType(instance, t)
}

func Transfrom[A, B any](s *summoner.Summoner[A]) *summoner.Summoner[B] {
	return summoner.Transfrom[A, B](s)
}

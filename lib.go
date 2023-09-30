package summoner

import (
	"reflect"
)

func TypeOf[A any]() reflect.Type {
	return reflect.TypeOf((*A)(nil)).Elem()
}

func IsRule[A any]() bool {
	return TypeOf[A]().Kind() == reflect.Struct
}

func Rules() string {
	return global.Rules()
}

func Summon[I any]() (I, error) {
	return Transfrom[any, I](&global).Summon()
}

func SummonType(t reflect.Type) (any, error) {
	return global.SummonType(t)
}

func Given[I any](instance I) error {
	return Transfrom[any, I](&global).Given(instance)
}

func GivenType(instance any, t reflect.Type) error {
	return global.GivenType(instance, t)
}

func Transfrom[A, B any](s *Summoner[A]) *Summoner[B] {
	return &Summoner[B]{
		instances: s.instances,
		rules:     s.rules,
	}
}

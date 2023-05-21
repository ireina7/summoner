package summoner

import (
	"reflect"
)

type Typeclass[A any] interface {
	Given(A) error
	GivenType(any, reflect.Type) error

	Summon() (A, error)
	SummonType(reflect.Type) (any, error)
}

type Summoner[A any] struct {
	instances map[reflect.Type]any
}

func TypeOf[A any]() reflect.Type {
	return reflect.TypeOf((*A)(nil)).Elem()
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

func Transfrom[A, B any](s *Summoner[A]) *Summoner[B] {
	return &Summoner[B]{
		instances: s.instances,
	}
}

func (self *Summoner[A]) tryBuild(t reflect.Type) (any, error) {
	// t := TypeOf[A]()
	// fmt.Println("Try build", t)
	r := reflect.New(t)
	i := 0
	var a A
	for i < r.Elem().NumField() {
		field := r.Elem().Field(i)
		instance, err := self.SummonType(field.Type())
		if err != nil {
			return a, err
		}
		dependency := reflect.ValueOf(instance)
		field.Set(dependency)
		i += 1
	}
	ans := r.Elem().Interface()
	return ans, nil
}

func (self *Summoner[A]) Summon() (A, error) {
	t := TypeOf[A]()
	// fmt.Println("Summon", t)
	ev, ok := self.instances[t]
	if ok {
		return ev.(A), nil
	}
	var a A
	// return a, fmt.Errorf("Instance of %v not found", t)
	x, err := self.tryBuild(t)
	if err != nil {
		return a, err
	}
	return x.(A), nil

}

func (self *Summoner[A]) SummonType(t reflect.Type) (any, error) {
	// fmt.Println("Summon", t)
	ev, ok := self.instances[t]
	if ok {
		return ev, nil
	}
	// return nil, fmt.Errorf("Instance of %v not found", t)
	var a A
	x, err := self.tryBuild(t)
	if err != nil {
		return a, err
	}
	return x, nil
}

func (self *Summoner[A]) Given(ev A) error {
	t := TypeOf[A]()
	return self.GivenType(ev, t)
}

func (self *Summoner[A]) GivenType(ev any, t reflect.Type) error {
	self.instances[t] = ev
	return nil
}

var global Summoner[any] = Summoner[any]{
	instances: map[reflect.Type]any{},
}

type RType = reflect.Type
type Type struct {
	RType
	params []Type
}

func fromReflect(t reflect.Type) Type {
	return Type{
		RType:  t,
		params: []Type{},
	}
}

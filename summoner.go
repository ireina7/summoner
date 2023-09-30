package summoner

import (
	"fmt"
	"reflect"
	// "github.com/traefik/yaegi/stdlib"
)

type Typeclass[A any] interface {
	Given(A) error
	GivenType(any, reflect.Type) error

	Summon() (A, error)
	SummonType(reflect.Type) (any, error)
}

type Transmute[A any] interface {
	Transform() A
}

// func ValueFromTypeName(tname string) (any, error) {
// 	i := interp.New(interp.Options{GoPath: "/Users/comcx/go"})
// 	// i.CompilePath(".")
// 	_, err := i.Eval(`.`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	x, err := i.Eval(fmt.Sprintf(`%s{}`, tname))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return x.Interface(), nil
// }

type Summoner[A any] struct {
	instances map[reflect.Type]any
	rules     map[reflect.Type]any
}

type summonError struct {
	want reflect.Type
}

func (err *summonError) Error() string {
	return fmt.Sprintf("Summon error: expected type %v", err.want)
}

func (self *Summoner[A]) tryBuild(t reflect.Type) (any, error) {
	if t.Kind() == reflect.Interface {
		return nil, &summonError{t}
	}
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
	// Cache
	self.GivenType(ans, t)
	return ans, nil
}

func (self *Summoner[A]) Summon() (A, error) {
	t := TypeOf[A]()
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

func (self *Summoner[A]) Inspect() string {
	devils := ""
	for k, v := range self.instances {
		devils += fmt.Sprintf("\t%v:\t%v\n", k, v)
	}
	return fmt.Sprintf("Devils[%d] {\n%v}",
		len(self.instances),
		devils,
	)
}

func (self *Summoner[A]) Rules() string {
	rules := ""
	for k, v := range self.rules {
		rules += fmt.Sprintf("\t%v: %v\n", k, v)
	}
	return fmt.Sprintf("Rules[%d] {\n%v}",
		len(self.rules),
		rules,
	)
}

var global Summoner[any] = Summoner[any]{
	instances: map[reflect.Type]any{},
	rules:     map[reflect.Type]any{},
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

func (self *Summoner[A]) Inject(structPtr any) error {
	t := reflect.TypeOf(structPtr)
	if t.Kind() != reflect.Pointer {
		return fmt.Errorf("t.Kind(%v) is not reflect.Pointer", t.Kind())
	}
	v := reflect.ValueOf(structPtr).Elem()
	for i := 0; i < v.NumField(); i += 1 {
		field := v.Field(i)
		ft := field.Type()
		tag := t.Elem().Field(i).Tag.Get("summon")
		if len(tag) == 0 {
			switch ft.Kind() {
			case reflect.Pointer:
				self.Inject(field.Interface()) //recursively
			case reflect.Struct:
				self.Inject(field.Addr().Interface()) //recursively
			default:
				//pass
			}
			continue
		}

		//It's leave node
		switch ft.Kind() {
		case reflect.Struct:
			x, err := self.tryBuild(field.Type())
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(x))
		case reflect.Interface:
			x, err := self.SummonType(field.Type())
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(x))
		case reflect.Pointer:
			x, err := self.tryBuild(field.Type().Elem())
			if err != nil {
				return err
			}
			p := reflect.New(ft.Elem())
			val := reflect.ValueOf(x)
			p.Elem().Set(val)
			field.Set(p)
		default:
			x, err := self.SummonType(field.Type())
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(x))
		}
	}
	return nil
}

func main() {
	fmt.Println("Hello, summoner")
}

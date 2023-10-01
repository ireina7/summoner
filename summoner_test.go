package summoner

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSummoner(t *testing.T) {

	// Given[Show[Person]](new(ShowPerson))
	GivenType(new(ShowPerson), TypeOf[Show[Person]]())
	Given[Show[int]](new(ShowInt))
	Given[Show[string]](new(ShowString))
	test()
	global.GivenRule(
		TypeOf[Debug[A]](),
		TypeOf[Debug[A]](),
	)

	fmt.Println(Rules())

	ev, err := Summon[Show[Person]]()
	if err != nil {
		panic(err)
	}
	t.Log(ev.Show(Person{0, "Tom", 10}))

	ee, err := Summon[Show[int]]()
	if err != nil {
		panic(err)
	}
	t.Log(ee.Show(7))

	es, err := Summon[Debug[Person]]()
	if err != nil {
		panic(err)
	}
	t.Log(es.Debug(Person{0, "Tom", 10}))

	ed, err := Summon[Recursive[Person]]()
	if err != nil {
		panic(err)
	}
	ed.Log(Person{1, "Jack", 14})

	ea, err := SummonType(TypeOf[App[Person]]())
	if err != nil {
		panic(err)
	}
	xx := ea.(App[Person])
	xx.Execute(Person{1, "Jack", 14})

	// t.Log(global.Inspect())

	// x, err := ValueFromTypeName("Person")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%#v\n", x)
}

type Person struct {
	id   int
	name string
	age  int
}

type Show[A any] interface {
	Show(a A) string
}

type ShowPerson struct{}

func (ev *ShowPerson) Show(p Person) string {
	return fmt.Sprintf("Person: %v", p)
}

type ShowInt struct{}

func (ev *ShowInt) Show(i int) string {
	return fmt.Sprintf("Integer: %v", i)
}

type ShowString struct{}

func (ev *ShowString) Show(s string) string {
	return "[str] " + s
}

type Debug[A any] struct {
	Show Show[A] `showtag`
}

func (self *Debug[A]) Debug(a A) string {
	return fmt.Sprintf("Debug: %s", self.Show.Show(a))
}

type Recursive[A any] struct {
	Debugger Debug[A]
}

func (self *Recursive[A]) Log(a A) {
	fmt.Println(self.Debugger.Debug(a))
	fmt.Println(self.Debugger.Show.Show(a))
}

type App[A any] struct {
	Show      Show[string]
	Recursive Recursive[A]
}

func (app *App[A]) Execute(a A) {
	fmt.Println("---------------------")
	fmt.Println(app.Show.Show("This is a test!"))
	app.Recursive.Log(a)
}

type Monoid[A any] interface {
	Zero() A
	Plus(A, A) A
}

type ForAll[A any] struct{}

type MonoidSlice[A any] struct{}

func (self *MonoidSlice[A]) Zero() []A {
	return []A{}
}

func (self *MonoidSlice[A]) Plus(a, b []A) []A {
	c := a
	c = append(c, b...)
	return c
}

// t: the rule name
func (self *Summoner[A]) GivenRule(from reflect.Type, to reflect.Type) error {
	self.rules[from] = to
	return nil
}

func test() {
	global.GivenRule(
		TypeOf[MonoidSlice[A]](),
		TypeOf[Monoid[[]A]](),
	)
}

type Service struct {
	Version  string
	Debugger *Debug[string] `summon:"type"`
	Device   *Device
}

type Device struct {
	Id   int `summon:"type"`
	Name string
	Show Show[string] `summon:"type"`
}

func TestInject(t *testing.T) {
	Given[Show[string]](new(ShowString))
	Given[int](-7)
	device := &Device{
		Name: "sp",
	}
	service := &Service{
		Version: "0.1.0",
		Device:  device,
	}
	err := global.Inject(service)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Service: %#v", service)
	t.Log(service.Debugger.Debug("sss"))
	t.Log(service.Device.Show.Show("ss"))
}

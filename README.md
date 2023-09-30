# Summoner
Typeclass system and dependency injection for Golang.

Let's summon magic!

![banner](banner.jpg)


## Getting started
### Get me
```shell
go get -u github.com/ireina7/summoner
```

### Usage
Summoner has two basic functions: `Summon` and `Given`.
`Summon[T]()` summons a value of type `T`,
`Given[T](T)` injects a value of type `T`,
`Inject(any)` injects any struct pointer recursively while still persist existing fields.
```go
import "fmt"
import summoner "github.com/ireina7/summoner"

func Summon[A any]() (A, error) {
    return summoner.Summon[A]()
}

func Given[A any](a A) error {
    return summoner.Given(a)
}


func main() {
    Given[int](0) // Default integer!

    fmt.Println(Summon[int]())
}
```



## Basic dependency injection
You can use `summoner` to achieve simple typeclass function.
For example, if you want to have a typeclass `Show`, 
and implement it for types `int`, `string` and `Person`:
```go
type Show[A any] interface {
    Show(a A) string
}

type ShowInt struct{}
func (self *ShowInt) Show(i int) string {
    fmt.Printf("Show int %v", i)
}

type ShowString struct{}
func (self *ShowString) Show(s string) string {
    fmt.Printf("Show str %v", s)
}

type Person struct {
    id int
    name string
    age int
}
type ShowPerson struct{}
func (self *ShowPerson) Show(p Person) string {
    fmt.Printf("Show person %v", p)
}


func main() {
    Given[Show[int]](new(ShowInt))
    Given[Show[string]](new(ShowString))
    Given[Show[Person]](new(ShowPerson))

    si, err := Summon[Show[int]]()
    if err != nil {
        panic(err)
    }
    fmt.Println(si.Show(7))

    ss, err := Summon[Show[String]]()
    if err != nil {
        panic(err)
    }
    fmt.Println(ss.Show("str"))

    sp, err := Summon[Show[Person]]()
    if err != nil {
        panic(err)
    }
    fmt.Println(sp.Show(Person{0, "Tom", 10}))
}
```

## Summon rules to summon inferred value
Suppose you want to build instance from some existed instance,
for example, 
- *if a type has implemented typeclass `Show`, then it should implement `Debug` automatically* or
- *For all types implemented `Show`, they also implement `Debug` dynamically*

Here's the code:
```go
type Debug[A any] struct {
    Show Show[A] // If implemented `Show`, then implement `Debug`
}

func (self *Debug[A]) Debug(a A) string {
    return fmt.Sprintf("Debug: %s", self.Show.Show(a))
}


func main() {
    Given[Show[int]](new(ShowInt))
    Given[Show[String]](new(ShowString))
    Given[Show[Person]](new(ShowPerson))

    si, err := Summon[Debug[int]]() // Get Debug[int] for free
    if err != nil {
        panic(err)
    }
    fmt.Println(si.Debug(7))

    sp, err := Summon[Debug[Person]]() // Get Debug[Person] for free
    if err != nil {
        panic(err)
    }
    fmt.Println(sp.Debug(Person{0, "Tom", 10}))
}

```

### Inject existing value
To inject some fields(e.g., summon only subset of a struct field set),
Tag fields to be summoned with `summon:"true"`.
```go
type Service struct {
	Version  string
	Debugger *Debug[string] `summon:"true"`
	Device   Device
}

type Device struct {
	Id   int `summon:"true"`
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
		Device:  *device,
	}
	err := global.Inject(service)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Service: %#v", service)
	t.Log(service.Debugger.Debug("sss"))
}
```

## An example for `Monoid` and `Fold`
```go
type Monoid[A any] interface {
    Zero() A
    Plus(a, b A) A
}

type MonoidString struct{}

func (self *MonoidString) Zero() string {
    return ""
}

func (self *MonoidString) Plus(a, b string) string {
    return a + " :+: " + b
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


func main() {
    Given[Monoid[string]](new(MonoidString))
    ss, err := Summon[Monoid[string]]()
    if err != nil {
        panic(err)
    }
    fmt.Println(ss.Plus("Hello", "world"))

    sf, err := Summon[Fold[string]]() // Get `Fold[string]` for free!
    if err != nil {
        panic(err)
    }
    fmt.Println(sf.FoldLeft([]string{"Hello", "world", "fst", "snd"}))
}
```

## API for golang compiler version before 1.18
Summoner also has API for non-generic style:
```go
func SummonType(t reflect.Type) (any, error)
func GivenType(instance any, t reflect.Type) error
```
Without generics, one need to pass `reflect.Type` value.


## Todo
- Better error message
- Optimization for runtime
- Forall type constraint a little verbose
- RPC server

## Contribution
Any new ideas?
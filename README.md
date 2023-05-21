# summoner
Typeclass system and dependency injection for Golang.
![banner](banner.jpg)

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
    Given[Show[int]](new(ShowString))
    Given[Show[int]](new(ShowPerson))

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
    Given[Show[int]](new(ShowString))
    Given[Show[int]](new(ShowPerson))

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
package summoner

type TypeVar interface {
	GetName() string
}

type A struct{}
type B struct{}
type C struct{}
type D struct{}
type E struct{}
type F struct{}

func (t A) GetName() string {
	return "A"
}
func (t B) GetName() string {
	return "B"
}
func (t C) GetName() string {
	return "C"
}
func (t D) GetName() string {
	return "D"
}
func (t E) GetName() string {
	return "E"
}
func (t F) GetName() string {
	return "F"
}

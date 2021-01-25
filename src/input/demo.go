package input

type DemoInput struct{}

func (server *DemoInput) Entity(name string) (*Entity, error) {
	var e Entity

	e.Name = name
	return &e, nil
}
func (server *DemoInput) Entities(pattern string) (*[]Entity, error) {
	return nil, nil
}

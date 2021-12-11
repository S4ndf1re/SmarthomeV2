package main

type Container struct {
	List []*Child `json:"list"`
}

func NewContainer() *Container {
	container := new(Container)
	container.List = make([]*Child, 0)
	return container
}

func (container *Container) Add(child *Child) {
	container.List = append(container.List, child)
}

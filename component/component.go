package component

import (
	"github.com/trustmaster/goflow"
)

const (
	ComponentModePool int8 = flow.ComponentModePool
	ComponentModeSync      = flow.ComponentModeSync
)

//ProcessFn is run for every input item
type ProcessFn func(val interface{}, out chan<- interface{}) error

//Component defines a dynamic flow component
type Component struct {
	flow.Component // component "superclass" embedded
	Name           string
	In             <-chan interface{} // input port
	Out            chan<- interface{} // output port
	Process        ProcessFn
}

// OnIn reacts a new input item
func (c *Component) OnIn(val interface{}) {
	c.Process(val, c.Out)
}

//NewComponent creates a new component
func NewComponent(name string, process ProcessFn) *Component {
	return &Component{Name: name, Process: process}
}

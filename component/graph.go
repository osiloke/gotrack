package component

import (
	"github.com/trustmaster/goflow"
)

//Graph extends flow graph, holds a set of components
type Graph struct {
	flow.Graph
}

//NewLinearGraph creates a linear graph which executes components in parallel
func NewLinearGraph(in chan interface{}, components ...*Component) *Graph {
	n := new(Graph)    // creates the object in heap
	n.InitGraphState() // allocates memory for the graph
	// Add processes to the network
	var (
		first = components[0]
		last  *Component
	)
	for _, component := range components {
		n.Add(component, component.Name)
		if last != nil {
			n.Connect(last.Name, "Out", component.Name, "In")
		}
		last = component
	}
	n.MapInPort("In", first.Name, "In")
	n.SetInPort("In", in)
	// Run the net
	flow.RunNet(n)
	return n
}

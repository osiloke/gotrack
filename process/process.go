package process

// https://rclayton.silvrback.com/pipelines-in-golang
//http://stackoverflow.com/questions/6395076/in-golang-using-reflect-how-do-you-set-the-value-of-a-struct-field

// _ "github.com/trustmaster/goflow"

type Process interface {
	Do(in chan interface{}) chan interface{}
}

type ProcessLine struct {
	head chan interface{}
	tail chan interface{}
}

func (p *ProcessLine) Enqueue(item string) {
	p.head <- item
}

func (p *ProcessLine) Dequeue(handler func(interface{})) {
	for i := range p.tail {
		handler(i)
	}
}

func (p *ProcessLine) Close() {
	close(p.head)
}

func NewProcessLine(processes ...Process) ProcessLine {
	head := make(chan interface{})
	var next_chan chan interface{}
	for _, process := range processes {
		if next_chan == nil {
			next_chan = process.Do(head)
		} else {
			next_chan = process.Do(next_chan)
		}
	}
	return ProcessLine{head: head, tail: next_chan}
}

package process

import (
	"time"

	"github.com/everdev/mack"
)

//https://www.snip2code.com/Snippet/415861/Golang-port-of-underscore-js
func debounce(interval time.Duration, input chan interface{}) (output chan interface{}) {
	// input := make(chan interface{})

	go func() {
		var buffer interface{}
		var ok bool

		// We do not start waiting for interval until called at least once
		buffer, ok = <-input
		// If channel closed exit, we could also close output
		if !ok {
			return
		}

		// We start waiting for an interval
		for {
			select {
			case buffer, ok = <-input:
				// If channel closed exit, we could also close output
				if !ok {
					return
				}

			case <-time.After(interval):
				// Interval has passed and we have data, so send it
				output <- buffer
				// Wait for data again before starting waiting for an interval
				buffer, ok = <-input
				if !ok {
					return
				}
				// If channel is not closed we have more data and start waiting for interval
			}
		}
	}()

	return input
}

//PowerControl monitors a power button
type PowerControl struct {
	status bool
	in     chan interface{}
}

//SetIn sets the input channel for this control process
///common can be extracted to base class
func (p PowerControl) SetIn(in chan interface{}) {
	p.in = in
}

//In sends an input to the input channel
func (p PowerControl) In(val interface{}) {
	p.in <- val
}

//Do is the control loop
func (p PowerControl) Do() {
	go func() {
		var (
			interval = 1 * time.Second
		)
		_in := debounce(interval, p.in)
		for {
			select {
			case <-_in:
				//out <- val
				p.status = !p.status
				if p.status {
					mack.Say("power on. Welcome? o c?!", "Victoria") //, "-v", "Victoria")
				} else {
					mack.Say("power off. Buy?buy! o c!?", "Victoria") //, "-v", "Victoria")
				}
			}
		}
	}()
}

//NewPowerControl creates a power control
func NewPowerControl(status bool) *PowerControl {
	p := PowerControl{status: status}
	p.in = make(chan interface{})
	p.Do()
	return &p
}

//DoFunction gets executed for every input in a single process input channel
type DoFunction func(c *SingleProcess, val interface{}) error

//SingleProcess is a chain in a process pipeline that performs an action on a
//sensor input
type SingleProcess struct {
	in   chan interface{}
	doFn DoFunction
}

//SetIn sets the input channel
func (p *SingleProcess) SetIn(in chan interface{}) {
	p.in = in
}

//In sends input to the input channel
func (p *SingleProcess) In(val interface{}) {
	p.in <- val
}

////////////

func (p *SingleProcess) do() {
	go func() {
		var (
			interval = 1 * time.Second
		)
		_in := debounce(interval, p.in)
		for {
			select {
			case v := <-_in:
				if p.doFn(p, v) != nil {
					return
				}
			}
		}
	}()
}

//NewSingleProcess create a single process
func NewSingleProcess(doFn DoFunction) *SingleProcess {
	p := SingleProcess{}
	p.SetIn(make(chan interface{}))
	p.doFn = doFn
	p.do()
	return &p
}

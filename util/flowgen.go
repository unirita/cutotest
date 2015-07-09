package util

import ()

type Flowgen struct {
	command
}

func NewFlowgen() *Flowgen {
	f := new(Flowgen)
	f.GeneratePath("flowgen")
	return f
}

func (f *Flowgen) Run(params ...string) (int, error) {
	return f.Exec(params...)
}

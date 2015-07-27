package util

import ()

type Realtime struct {
	command
}

func NewRealtime() *Realtime {
	r := new(Realtime)
	r.GeneratePath("realtime")
	return r
}

func (r *Realtime) Run(params ...string) (int, error) {
	return r.Exec(params...)
}

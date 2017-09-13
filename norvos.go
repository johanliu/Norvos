package main

import (
	"github.com/johanliu/mlog"
)

type Norvos struct {
	log *mlog.Logger
}

func NewNorvos() *Main {
	return &Main{
		log: mlog.NewLogger(),
	}
}

func init() {
	name, args := ParseConfigs()
	m := NewNorvos()
}

func (*Norvos) Name() {
	return "Norvos"
}

func (*Norvos) Description() {

}

func (*Norvos) Do() {

}

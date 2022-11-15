package ioext

import (
	"io"
	"log"
)

func Close(closer io.Closer) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			log.Printf("Unable to close %T: %s\n", closer, err.Error())
		}
	}
}

type CloserGroup struct {
	group []func()
}

func NewCloserGroup(closers ...func()) *CloserGroup {
	group := make([]func(), 0, len(closers))
	for _, closer := range closers {
		group = append(group, closer)
	}
	return &CloserGroup{group: group}
}

func (g *CloserGroup) Add(closer func()) {
	if closer != nil {
		g.group = append(g.group, closer)
	}
}

func (g *CloserGroup) Close() {
	for _, closer := range g.group {
		closer()
	}
}

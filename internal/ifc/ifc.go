package ifc

import (
	"github.com/monopole/mdparse/internal/usegold/loader"
)

type Marker interface {
	Load(set *loader.MyContrivedFolder) error
	Render() (string, error)
	Dump()
}

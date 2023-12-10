package ifc

import "github.com/monopole/mdrip/base"

type Marker interface {
	Load(set *base.DataSet) error
	Render() (string, error)
	Dump()
}

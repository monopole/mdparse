package ifc

type Marker interface {
	Render() (string, error)
	Dump()
}

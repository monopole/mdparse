package ifc

type Marker interface {
	Parse(bytes []byte) error
	Render() (string, error)
	Dump()
}

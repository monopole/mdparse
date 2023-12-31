package ifc

type Marker interface {
	Load(rawData []byte) error
	Render() (string, error)
	Dump()
}

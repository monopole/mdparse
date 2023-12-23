package ifc

type Marker interface {
	//Load(set *loader.MyContrivedFolder) error
	Render() (string, error)
	Dump()
}

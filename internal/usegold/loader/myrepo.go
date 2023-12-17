package loader

// MyRepo is a named group of files and folders.
type MyRepo struct {
	// name is the URL of the repo, e.g.
	// https://githu.com/monopole/mdrip
	name   string
	folder *MyFolder
}

var _ MyTreeItem = &MyRepo{}

func (r *MyRepo) Accept(v TreeVisitor) {
	v.VisitRepo(r)
}

func (r *MyRepo) Parent() MyTreeItem {
	return nil
}

func (r *MyRepo) FullName() string {
	return r.name
}

func (r *MyRepo) DirName() string {
	return ""
}

func (r *MyRepo) Name() string {
	return r.name
}

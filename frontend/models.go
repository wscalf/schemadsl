package frontend

type Service struct {
	Name   string
	Assets []*Asset
}

type Asset struct {
	Name        string
	IsPublic    bool
	Permissions []*Permission
}

type Permission struct {
	IsPublic bool
	Name     string
}

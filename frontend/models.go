package frontend

type Service struct {
	Name    string
	Imports map[string]string
	Assets  map[string]*Asset
}

type Asset struct {
	Name         string
	IsPublic     bool
	Permissions  map[string]*Permission
	Dependencies map[string]AssetReference
}

type Permission struct {
	IsPublic bool
	Name     string
}

type AssetReference struct {
	Service string
	Asset   string
}

type PermissionReference struct {
	AssetReference
	Permission string
}

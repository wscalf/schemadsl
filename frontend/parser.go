package frontend

import (
	"os"
	"schemadsl/gen"

	"github.com/antlr4-go/antlr/v4"
)

func LoadFile(path string) *Service {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	input := antlr.NewInputStream(string(data))
	lexer := gen.NewauthzLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)
	parser := gen.NewauthzParser(stream)

	tree := parser.File()

	visitor := &dslVisitor{
		service: &Service{
			Imports: map[string]string{},
			Assets:  map[string]*Asset{},
		},
	}

	if start, ok := tree.(*gen.FileContext); ok {
		visitor.VisitFile(start)
		return visitor.service
	} else {
		panic("Unrecognized root element type")
	}
}

type dslVisitor struct {
	service *Service
}

func (v *dslVisitor) VisitFile(ctx *gen.FileContext) {
	for _, statement := range ctx.AllStatement() {
		v.VisitStatement(statement)
	}
}

func (v *dslVisitor) VisitStatement(ctx gen.IStatementContext) {
	switch val := ctx.(type) {
	case *gen.VersionContext:
		break //Skip for now
	case *gen.ServiceContext:
		v.service.Name = v.VisitServiceStatement(val)
	case *gen.AssetContext:
		asset := v.VisitAssetStatement(val)
		v.service.Assets[asset.Name] = asset
	case *gen.ImportContext:
		alias, serviceName := v.VisitImportStatement(val)
		v.service.Imports[alias] = serviceName
	}
}

func (v *dslVisitor) VisitVersionStatement(ctx *gen.VersionContext) string {
	return ctx.VERSIONNO().GetText()
}

func (v *dslVisitor) VisitServiceStatement(ctx *gen.ServiceContext) string {
	return ctx.NAME().GetText()
}

func (v *dslVisitor) VisitAssetStatement(ctx *gen.AssetContext) *Asset {
	asset := &Asset{
		Permissions:  map[string]*Permission{},
		Dependencies: map[string]AssetReference{},
	}

	asset.Name = ctx.NAME().GetText()

	access := ctx.ACCESS()
	if access == nil {
		asset.IsPublic = true
	} else {
		if access.GetText() == "public" { //Is there a way to check against the actual keyword?
			asset.IsPublic = true
		}
	}

	for _, dependency := range ctx.AllDependency() {
		alias, reference := v.VisitDependency(dependency)
		asset.Dependencies[alias] = reference
	}

	for _, permission := range ctx.AllPermission() {
		permission := v.VisitPermission(permission)
		asset.Permissions[permission.Name] = permission
	}

	for _, computed := range ctx.AllComputedPermission() {
		//Handle computed permissions
		_ = computed
	}

	return asset
}

func (v *dslVisitor) VisitPermission(ctx gen.IPermissionContext) *Permission {
	permission := &Permission{}

	access := ctx.ACCESS()
	if access == nil {
		permission.IsPublic = true
	} else {
		if access.GetText() == "public" { //Is there a way to check against the actual keyword?
			permission.IsPublic = true
		}
	}

	permission.Name = ctx.NAME().GetText()

	return permission
}

func (v *dslVisitor) VisitImportStatement(ctx *gen.ImportContext) (string, string) {
	service := ctx.GetService().GetText()
	alias := ctx.GetAlias()

	if alias == nil {
		return service, service
	} else {
		return alias.GetText(), service
	}
}

func (v *dslVisitor) VisitDependency(ctx gen.IDependencyContext) (string, AssetReference) {
	serviceName := ctx.GetService().GetText()
	assetName := ctx.GetAsset().GetText()
	alias := ctx.GetAlias().GetText()

	return alias, AssetReference{
		Service: serviceName,
		Asset:   assetName,
	}
}

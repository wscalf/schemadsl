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
		service: &Service{},
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
		v.service.Assets = append(v.service.Assets, v.VisitAssetStatement(val))
	}
}

func (v *dslVisitor) VisitVersionStatement(ctx *gen.VersionContext) string {
	return ctx.VERSIONNO().GetText()
}

func (v *dslVisitor) VisitServiceStatement(ctx *gen.ServiceContext) string {
	return ctx.NAME().GetText()
}

func (v *dslVisitor) VisitAssetStatement(ctx *gen.AssetContext) *Asset {
	asset := &Asset{}

	asset.Name = ctx.NAME().GetText()

	access := ctx.ACCESS()
	if access == nil {
		asset.IsPublic = true
	} else {
		if access.GetText() == "public" { //Is there a way to check against the actual keyword?
			asset.IsPublic = true
		}
	}

	for _, permission := range ctx.AllPermission() {
		asset.Permissions = append(asset.Permissions, v.VisitPermission(permission))
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

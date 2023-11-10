package backend

import (
	"fmt"
	"os"
	"schemadsl/frontend"
	"strings"
)

const (
	rolePermsPlaceholder        = "// <GENERATED_ROLE_PERMS_HERE>"
	roleBindingPermsPlaceholder = "// <GENERATED_ROLEBINDING_PERMS_HERE>"
	workspacePermsPlaceholder   = "// <GENERATED_WORKSPACE_PERMS_HERE>"
	definitionsPlaceholder      = "// <GENERATED_DEFINITIONS_HERE>"
)

func WriteOutput(outputFile string, services []*frontend.Service) {
	var rolePermsBuilder strings.Builder
	var roleBindingPermsBuilder strings.Builder
	var workspacePermsBuilder strings.Builder
	var definitionsBuilder strings.Builder

	for _, service := range services {
		for _, asset := range service.Assets {
			definitionsBuilder.WriteString(fmt.Sprintf("definition %s/%s {\n", service.Name, asset.Name))
			definitionsBuilder.WriteString("\trelation workspace: workspace\n")
			definitionsBuilder.WriteString("\n")

			for _, permission := range asset.Permissions {
				fullname := fmt.Sprintf("%s_%s_%s", service.Name, asset.Name, permission.Name)
				definitionsBuilder.WriteString(fmt.Sprintf("\tpermission %s = workspace->%s\n", permission.Name, fullname))

				rolePermsBuilder.WriteString(fmt.Sprintf("\trelation %s: user:*\n", fullname))
				roleBindingPermsBuilder.WriteString(fmt.Sprintf("\tpermission %s = subject & granted->%s\n", fullname, fullname))
				workspacePermsBuilder.WriteString(fmt.Sprintf("\tpermission %s = user_grant->%s + parent->%s\n", fullname, fullname, fullname))
			}

			definitionsBuilder.WriteString("}\n\n")
		}
	}

	original, err := os.ReadFile("template.zed")
	if err != nil {
		panic(err)
	}
	text := string(original)

	text = strings.ReplaceAll(text, rolePermsPlaceholder, rolePermsBuilder.String())
	text = strings.ReplaceAll(text, roleBindingPermsPlaceholder, roleBindingPermsBuilder.String())
	text = strings.ReplaceAll(text, workspacePermsPlaceholder, workspacePermsBuilder.String())
	text = strings.ReplaceAll(text, definitionsPlaceholder, definitionsBuilder.String())

	os.WriteFile(outputFile, []byte(text), os.FileMode(0644))
}

package module

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl"
	"io/ioutil"
)

type terraModuleManifests struct {
	Modules []TerraModuleManifest `hcl:"module"`
}

type TerraModuleManifest struct {
	Name string `hcl:",key"`
	Path string `hcl:"path,omitempty"`
}

func LoadModuleDependencies(dependenciesManifestPath string) (*map[string]TerraModuleManifest, error) {
	bytes, err := ioutil.ReadFile(dependenciesManifestPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading dependencies manifest file %s. %#v", dependenciesManifestPath, err))
	}

	manifests, err := parseDependenciesManifest(string(bytes))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing dependencies manifest file %s. %#v", dependenciesManifestPath, err))
	}

	manifestMap := make(map[string]TerraModuleManifest)
	for _, moduleManifest := range (*manifests).Modules {
		manifestMap[moduleManifest.Name] = moduleManifest
	}

	return &manifestMap, nil
}

func parseDependenciesManifest(hclText string) (*terraModuleManifests, error) {
	result := &terraModuleManifests{}

	hclParseTree, err := hcl.Parse(hclText)
	if err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&result, hclParseTree); err != nil {
		return nil, err
	}

	return result, nil
}

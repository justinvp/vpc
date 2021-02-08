// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	dotnetgen "github.com/pulumi/pulumi/pkg/v2/codegen/dotnet"
	gogen "github.com/pulumi/pulumi/pkg/v2/codegen/go"
	nodejsgen "github.com/pulumi/pulumi/pkg/v2/codegen/nodejs"
	pygen "github.com/pulumi/pulumi/pkg/v2/codegen/python"
	"github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
)

const Tool = "pulumi-gen-vpc"

// Language is the SDK language.
type Language string

const (
	DotNet Language = "dotnet"
	Go     Language = "go"
	NodeJS Language = "nodejs"
	Python Language = "python"
	Schema Language = "schema"
)

func main() {
	printUsage := func() {
		fmt.Printf("Usage: %s <language> <out-dir> [schema-file] [version]\n", os.Args[0])
	}

	args := os.Args[1:]
	if len(args) < 2 {
		printUsage()
		os.Exit(1)
	}

	language, outdir := Language(args[0]), args[1]

	var schemaFile string
	var version string
	if language != Schema {
		if len(args) < 4 {
			printUsage()
			os.Exit(1)
		}
		schemaFile, version = args[2], args[3]
	}

	switch language {
	case DotNet:
		genDotNet(readSchema(schemaFile, version), outdir)
	case Go:
		genGo(readSchema(schemaFile, version), outdir)
	case NodeJS:
		genNodeJS(readSchema(schemaFile, version), outdir)
	case Python:
		genPython(readSchema(schemaFile, version), outdir)
	case Schema:
		pkgSpec := generateSchema()
		mustWritePulumiSchema(pkgSpec, outdir)
	default:
		panic(fmt.Sprintf("Unrecognized language %q", language))
	}
}

const awsVersion = "v3.28.0"

func awsRef(ref string) string {
	return fmt.Sprintf("/aws/%s/schema.json%s", awsVersion, ref)
}

// nolint: lll
func generateSchema() schema.PackageSpec {
	return schema.PackageSpec{
		Name:        "vpc",
		Description: "Pulumi Amazon Web Services (AWS) VPC component",
		License:     "Apache-2.0",
		Keywords:    []string{"pulumi", "aws", "vpc"},
		Homepage:    "https://pulumi.com",
		Repository:  "https://github.com/pulumi/pulumi", // TODO
		Resources: map[string]schema.ResourceSpec{
			"vpc:index:Vpc": {
				IsComponent: true,
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Properties: map[string]schema.PropertySpec{
						"vpcId": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "The ID of the VPC resource.",
						},
						"vpcResource": {
							TypeSpec:    schema.TypeSpec{Ref: awsRef("#/resources/aws:ec2%2Fvpc:Vpc")},
							Description: "The VPC resource.",
						},
					},
					Required: []string{
						"vpcId",
						"vpcResource",
					},
				},
				InputProperties: map[string]schema.PropertySpec{
					"assignGeneratedIpv6CidrBlock": {
						TypeSpec: schema.TypeSpec{Type: "boolean"},
						Description: "Requests an Amazon-provided IPv6 CIDR block with a /56 prefix length for the " +
							"VPC. You cannot specify the range of IP addresses, or the size of the CIDR block. " +
							"Default is `false`.  If set to `true`, then subnets created will default to " +
							"`assignIpv6AddressOnCreation: true` as well.",
					},
					"cidrBlock": {
						TypeSpec:    schema.TypeSpec{Type: "string"},
						Description: "The CIDR block for the VPC. Defaults to \"10.0.0.0/16\" if unspecified.",
					},
					"enableClassiclink": {
						TypeSpec: schema.TypeSpec{Type: "boolean"},
						Description: "A boolean flag to enable/disable ClassicLink for the VPC. Only valid in " +
							"regions and accounts that support EC2 Classic. See the [ClassicLink documentation][1] " +
							"for more information. Defaults false.",
					},
					"enableClassiclinkDnsSupport": {
						TypeSpec: schema.TypeSpec{Type: "boolean"},
						Description: "A boolean flag to enable/disable ClassicLink DNS Support for the VPC. Only " +
							"valid in regions and accounts that support EC2 Classic.",
					},
					"enableDnsHostnames": {
						TypeSpec: schema.TypeSpec{Type: "boolean"},
						Description: "A boolean flag to enable/disable DNS hostnames in the VPC. Defaults to true if " +
							"unspecified.",
					},
					"enableDnsSupport": {
						TypeSpec: schema.TypeSpec{Type: "boolean"},
						Description: "A boolean flag to enable/disable DNS support in the VPC. Defaults true if " +
							"unspecified.",
					},
					"instanceTenancy": {
						TypeSpec: schema.TypeSpec{Type: "string"},
						Description: "A tenancy option for instances launched into the VPC. Defaults to \"default\" " +
							"if unspecified.",
					},
					"tags": {
						TypeSpec: schema.TypeSpec{
							Type:                 "object",
							AdditionalProperties: &schema.TypeSpec{Type: "string"},
						},
						Description: "A mapping of tags to assign to the resource.",
					},
				},
			},
		},

		Types: map[string]schema.ComplexTypeSpec{},

		Language: map[string]json.RawMessage{
			"csharp": rawMessage(map[string]interface{}{
				"packageReferences": map[string]string{
					"Pulumi":     "2.*",
					"Pulumi.Aws": "3.*",
				},
			}),
			"go": rawMessage(map[string]interface{}{
				"generateResourceContainerTypes": true,
				"importBasePath":                 "github.com/justinvp/vpc/sdk/go/vpc",
			}),
			"nodejs": rawMessage(map[string]interface{}{
				"dependencies": map[string]string{
					"@pulumi/aws": "^3.28.0",
				},
				"devDependencies": map[string]string{
					"typescript": "^3.7.0",
				},
			}),
			"python": rawMessage(map[string]interface{}{
				"requires": map[string]string{
					"pulumi":     ">=2.20.0,<3.0.0",
					"pulumi-aws": ">=3.28.0,<4.0.0",
				},
				"readme": "Pulumi Amazon Web Services (AWS) VPC component.",
			}),
		},
	}
}

func rawMessage(v interface{}) json.RawMessage {
	bytes, err := json.Marshal(v)
	contract.Assert(err == nil)
	return bytes
}

func readSchema(schemaPath string, version string) *schema.Package {
	// Read in, decode, and import the schema.
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		panic(err)
	}

	var pkgSpec schema.PackageSpec
	if err = json.Unmarshal(schemaBytes, &pkgSpec); err != nil {
		panic(err)
	}
	pkgSpec.Version = version

	pkg, err := schema.ImportSpec(pkgSpec, nil)
	if err != nil {
		panic(err)
	}
	return pkg
}

func genDotNet(pkg *schema.Package, outdir string) {
	files, err := dotnetgen.GeneratePackage(Tool, pkg, map[string][]byte{})
	if err != nil {
		panic(err)
	}
	mustWriteFiles(outdir, files)
}

func genGo(pkg *schema.Package, outdir string) {
	files, err := gogen.GeneratePackage(Tool, pkg)
	if err != nil {
		panic(err)
	}
	mustWriteFiles(outdir, files)
}

func genNodeJS(pkg *schema.Package, outdir string) {
	files, err := nodejsgen.GeneratePackage(Tool, pkg, map[string][]byte{})
	if err != nil {
		panic(err)
	}
	mustWriteFiles(outdir, files)
}

func genPython(pkg *schema.Package, outdir string) {
	files, err := pygen.GeneratePackage(Tool, pkg, map[string][]byte{})
	if err != nil {
		panic(err)
	}
	mustWriteFiles(outdir, files)
}

func mustWriteFiles(rootDir string, files map[string][]byte) {
	for filename, contents := range files {
		mustWriteFile(rootDir, filename, contents)
	}
}

func mustWriteFile(rootDir, filename string, contents []byte) {
	outPath := filepath.Join(rootDir, filename)

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		panic(err)
	}
	err := ioutil.WriteFile(outPath, contents, 0600)
	if err != nil {
		panic(err)
	}
}

func mustWritePulumiSchema(pkgSpec schema.PackageSpec, outdir string) {
	schemaJSON, err := json.MarshalIndent(pkgSpec, "", "    ")
	if err != nil {
		panic(errors.Wrap(err, "marshaling Pulumi schema"))
	}
	mustWriteFile(outdir, "schema.json", schemaJSON)
}

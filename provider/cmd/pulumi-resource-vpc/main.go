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

//go:generate go run ./generate.go

package main

import (
	"github.com/pkg/errors"

	"github.com/pulumi/pulumi/sdk/v2/go/common/util/cmdutil"
	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

var providerName = "vpc"
var version = "0.0.1"

func main() {
	provider := pulumi.ProviderArgs{
		Name:    providerName,
		Version: version,
		Schema:  pulumiSchema,
		ConstructF: func(ctx *pulumi.Context, typ, name string, inputs *pulumi.ConstructInputs,
			options pulumi.ResourceOption) (pulumi.ConstructResult, error) {

			switch typ {
			case "vpc:index:Vpc":
				return constructVpc(ctx, typ, name, inputs, options)
			default:
				return pulumi.ConstructResult{}, errors.Errorf("unknown resource type %s", typ)
			}
		},
	}

	err := pulumi.ProviderMain(provider)

	if err != nil {
		cmdutil.ExitError(err.Error())
	}
}

func constructVpc(ctx *pulumi.Context, typ, name string, inputs *pulumi.ConstructInputs,
	options pulumi.ResourceOption) (pulumi.ConstructResult, error) {
	contract.Assert(typ == "vpc:index:Vpc")

	args := &VpcArgs{}
	err := inputs.SetArgs(args)
	if err != nil {
		return pulumi.ConstructResult{}, errors.Wrap(err, "setting args")
	}

	vpc, err := NewVpc(ctx, name, args, options)
	if err != nil {
		return pulumi.ConstructResult{}, errors.Wrap(err, "creating vpc")
	}

	return pulumi.ConstructResult{
		URN: vpc.URN(),
		State: pulumi.Map{
			"vpcId":       vpc.ID,
			"vpcResource": vpc.Vpc,
		},
	}, nil
}

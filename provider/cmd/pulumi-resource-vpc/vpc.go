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
	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

type Vpc struct {
	pulumi.ResourceState

	ID  pulumi.IDOutput
	Vpc *ec2.Vpc
}

func NewVpc(ctx *pulumi.Context,
	name string, args *VpcArgs, opts ...pulumi.ResourceOption) (*Vpc, error) {
	if args == nil {
		args = &VpcArgs{}
	}

	component := &Vpc{}
	err := ctx.RegisterComponentResource("vpc:index:Vpc", name, component, opts...)
	if err != nil {
		return nil, err
	}

	assignGeneratedIpv6CidrBlock := args.AssignGeneratedIpv6CidrBlock
	if assignGeneratedIpv6CidrBlock == nil {
		assignGeneratedIpv6CidrBlock = pulumi.Bool(false)
	}

	cidrBlock := args.CidrBlock
	if cidrBlock == nil {
		cidrBlock = pulumi.String("10.0.0.0/16")
	}

	enableDnsHostnames := args.EnableDnsHostnames
	if enableDnsHostnames == nil {
		enableDnsHostnames = pulumi.Bool(true)
	}

	enableDnsSupport := args.EnableDnsSupport
	if enableDnsSupport == nil {
		enableDnsSupport = pulumi.Bool(true)
	}

	instanceTenancy := args.InstanceTenancy
	if instanceTenancy == nil {
		instanceTenancy = pulumi.String("default")
	}

	component.Vpc, err = ec2.NewVpc(ctx, name, &ec2.VpcArgs{
		AssignGeneratedIpv6CidrBlock: assignGeneratedIpv6CidrBlock,
		CidrBlock:                    cidrBlock,
		EnableClassiclink:            args.EnableClassiclink,
		EnableClassiclinkDnsSupport:  args.EnableClassiclinkDnsSupport,
		EnableDnsHostnames:           enableDnsHostnames,
		EnableDnsSupport:             enableDnsSupport,
		InstanceTenancy:              instanceTenancy,
		Tags:                         args.Tags,
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}
	component.ID = component.Vpc.ID()

	err = ctx.RegisterResourceOutputs(component, pulumi.Map{
		"vpcId": component.Vpc.ID(),
	})
	if err != nil {
		return nil, err
	}

	return component, nil
}

// The set of arguments for constructing a Vpc component.
type VpcArgs struct {
	// TODO
	// /**
	//  * The information about what subnets to create per availability zone.  Defaults to one public and
	//  * one private subnet if unspecified.
	//  */
	//  subnets?: VpcSubnetArgs[];
	//  /**
	//   * The maximum number of availability zones to use in the current region.  Defaults to `2` if
	//   * unspecified.  Use `"all"` to use all the availability zones in the current region.
	//   */
	//  numberOfAvailabilityZones?: number | "all";
	//  /**
	//   * The max number of NAT gateways to create if there are any private subnets created.  A NAT
	//   * gateway enables instances in a private subnet to connect to the internet or other AWS
	//   * services, but prevent the internet from initiating a connection with those instances. A
	//   * minimum of '1' gateway is needed if an instance is to be allowed connection to the internet.
	//   *
	//   * If this is not set, a nat gateway will be made for each availability zone in the current
	//   * region. The first public subnet for that availability zone will be the one used to place the
	//   * nat gateway in.  If less gateways are requested than availability zones, then only that many
	//   * nat gateways will be created.
	//   *
	//   * Private subnets in an availability zone that contains a nat gateway will route through that
	//   * gateway.  Private subnets in an availability zone that does not contain a nat gateway will be
	//   * routed to the other nat gateways in a round-robin fashion.
	//   *
	//   * See https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-gateway.html for more details.
	//   *
	//   * Defaults to [numberOfAvailabilityZones].
	//   */
	//  numberOfNatGateways?: number;

	// Requests an Amazon-provided IPv6 CIDR block with a /56 prefix length for the VPC. You cannot
	// specify the range of IP addresses, or the size of the CIDR block. Default is `false`.  If set
	// to `true`, then subnets created will default to `assignIpv6AddressOnCreation: true` as well.
	AssignGeneratedIpv6CidrBlock pulumi.BoolInput `pulumi:"assignGeneratedIpv6CidrBlock"`

	// The CIDR block for the VPC. Defaults to "10.0.0.0/16" if unspecified.
	CidrBlock pulumi.StringInput `pulumi:"cidrBlock"`

	// A boolean flag to enable/disable ClassicLink
	// for the VPC. Only valid in regions and accounts that support EC2 Classic.
	// See the [ClassicLink documentation][1] for more information. Defaults false.
	EnableClassiclink pulumi.BoolInput `pulumi:"enableClassiclink"`

	// A boolean flag to enable/disable ClassicLink DNS Support for the VPC.
	// Only valid in regions and accounts that support EC2 Classic.
	EnableClassiclinkDnsSupport pulumi.BoolInput `pulumi:"enableClassiclinkDnsSupport"`

	// A boolean flag to enable/disable DNS hostnames in the VPC. Defaults to true if unspecified.
	EnableDnsHostnames pulumi.BoolInput `pulumi:"enableDnsHostnames"`

	// A boolean flag to enable/disable DNS support in the VPC. Defaults true if unspecified.
	EnableDnsSupport pulumi.BoolInput `pulumi:"enableDnsSupport"`

	// A tenancy option for instances launched into the VPC. Defaults to "default" if unspecified.
	InstanceTenancy pulumi.StringInput `pulumi:"instanceTenancy"` // TODO enum "default" | "dedicated"?

	// A mapping of tags to assign to the resource.
	Tags pulumi.StringMapInput `pulumi:"tags"`
}

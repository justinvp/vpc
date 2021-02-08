// *** WARNING: this file was generated by pulumi-gen-vpc. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "./utilities";

import * as aws from "@pulumi/aws";

export class Vpc extends pulumi.ComponentResource {
    /** @internal */
    public static readonly __pulumiType = 'vpc:index:Vpc';

    /**
     * Returns true if the given object is an instance of Vpc.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is Vpc {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === Vpc.__pulumiType;
    }

    /**
     * The ID of the VPC resource.
     */
    public /*out*/ readonly vpcId!: pulumi.Output<string>;
    /**
     * The VPC resource.
     */
    public /*out*/ readonly vpcResource!: pulumi.Output<aws.ec2.Vpc>;

    /**
     * Create a Vpc resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args?: VpcArgs, opts?: pulumi.ComponentResourceOptions) {
        let inputs: pulumi.Inputs = {};
        if (!(opts && opts.id)) {
            inputs["assignGeneratedIpv6CidrBlock"] = args ? args.assignGeneratedIpv6CidrBlock : undefined;
            inputs["cidrBlock"] = args ? args.cidrBlock : undefined;
            inputs["enableClassiclink"] = args ? args.enableClassiclink : undefined;
            inputs["enableClassiclinkDnsSupport"] = args ? args.enableClassiclinkDnsSupport : undefined;
            inputs["enableDnsHostnames"] = args ? args.enableDnsHostnames : undefined;
            inputs["enableDnsSupport"] = args ? args.enableDnsSupport : undefined;
            inputs["instanceTenancy"] = args ? args.instanceTenancy : undefined;
            inputs["tags"] = args ? args.tags : undefined;
            inputs["vpcId"] = undefined /*out*/;
            inputs["vpcResource"] = undefined /*out*/;
        } else {
            inputs["vpcId"] = undefined /*out*/;
            inputs["vpcResource"] = undefined /*out*/;
        }
        if (!opts) {
            opts = {}
        }

        if (!opts.version) {
            opts.version = utilities.getVersion();
        }
        super(Vpc.__pulumiType, name, inputs, opts, true /*remote*/);
    }
}

/**
 * The set of arguments for constructing a Vpc resource.
 */
export interface VpcArgs {
    /**
     * Requests an Amazon-provided IPv6 CIDR block with a /56 prefix length for the VPC. You cannot specify the range of IP addresses, or the size of the CIDR block. Default is `false`.  If set to `true`, then subnets created will default to `assignIpv6AddressOnCreation: true` as well.
     */
    readonly assignGeneratedIpv6CidrBlock?: pulumi.Input<boolean>;
    /**
     * The CIDR block for the VPC. Defaults to "10.0.0.0/16" if unspecified.
     */
    readonly cidrBlock?: pulumi.Input<string>;
    /**
     * A boolean flag to enable/disable ClassicLink for the VPC. Only valid in regions and accounts that support EC2 Classic. See the [ClassicLink documentation][1] for more information. Defaults false.
     */
    readonly enableClassiclink?: pulumi.Input<boolean>;
    /**
     * A boolean flag to enable/disable ClassicLink DNS Support for the VPC. Only valid in regions and accounts that support EC2 Classic.
     */
    readonly enableClassiclinkDnsSupport?: pulumi.Input<boolean>;
    /**
     * A boolean flag to enable/disable DNS hostnames in the VPC. Defaults to true if unspecified.
     */
    readonly enableDnsHostnames?: pulumi.Input<boolean>;
    /**
     * A boolean flag to enable/disable DNS support in the VPC. Defaults true if unspecified.
     */
    readonly enableDnsSupport?: pulumi.Input<boolean>;
    /**
     * A tenancy option for instances launched into the VPC. Defaults to "default" if unspecified.
     */
    readonly instanceTenancy?: pulumi.Input<string>;
    /**
     * A mapping of tags to assign to the resource.
     */
    readonly tags?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
}

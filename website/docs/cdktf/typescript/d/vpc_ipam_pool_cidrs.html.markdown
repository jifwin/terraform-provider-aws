---
subcategory: "VPC IPAM (IP Address Manager)"
layout: "aws"
page_title: "AWS: aws_vpc_ipam_pool_cidrs"
description: |-
    Returns cidrs provisioned into an IPAM pool.
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_vpc_ipam_pool_cidrs

`aws_vpc_ipam_pool_cidrs` provides details about an IPAM pool.

This resource can prove useful when an ipam pool was shared to your account and you want to know all (or a filtered list) of the CIDRs that are provisioned into the pool.

## Example Usage

Basic usage:

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { Token, TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsVpcIpamPool } from "./.gen/providers/aws/data-aws-vpc-ipam-pool";
import { DataAwsVpcIpamPoolCidrs } from "./.gen/providers/aws/data-aws-vpc-ipam-pool-cidrs";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    const p = new DataAwsVpcIpamPool(this, "p", {
      filter: [
        {
          name: "description",
          values: ["*mypool*"],
        },
        {
          name: "address-family",
          values: ["ipv4"],
        },
      ],
    });
    new DataAwsVpcIpamPoolCidrs(this, "c", {
      ipamPoolId: Token.asString(p.id),
    });
  }
}

```

Filtering:

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { Token, TerraformIterator, Fn, TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsVpcIpamPoolCidrs } from "./.gen/providers/aws/data-aws-vpc-ipam-pool-cidrs";
import { Ec2ManagedPrefixList } from "./.gen/providers/aws/ec2-managed-prefix-list";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    const c = new DataAwsVpcIpamPoolCidrs(this, "c", {
      filter: [
        {
          name: "cidr",
          values: ["10.*"],
        },
      ],
      ipamPoolId: "ipam-pool-123",
    });
    const mycidrs =
      "${[ for cidr in ${" +
      c.ipamPoolCidrs +
      '} : cidr.cidr if cidr.state == "provisioned"]}';
    /*In most cases loops should be handled in the programming language context and 
    not inside of the Terraform context. If you are looping over something external, e.g. a variable or a file input
    you should consider using a for loop. If you are looping over something only known to Terraform, e.g. a result of a data source
    you need to keep this like it is.*/
    const plsDynamicIterator0 = TerraformIterator.fromList(
      Token.asAny(mycidrs)
    );
    new Ec2ManagedPrefixList(this, "pls", {
      addressFamily: "IPv4",
      maxEntries: Token.asNumber(Fn.lengthOf(mycidrs)),
      name: "IPAM Pool (${" + test.id + "}) Cidrs",
      entry: plsDynamicIterator0.dynamic({
        cidr: plsDynamicIterator0.value,
        description: plsDynamicIterator0.value,
      }),
    });
  }
}

```

## Argument Reference

The arguments of this data source act as filters for querying the available
VPCs in the current region. The given filters must match exactly one
VPC whose data will be exported as attributes.

* `ipamPoolId` - ID of the IPAM pool you would like the list of provisioned CIDRs.
* `filter` - Custom filter block as described below.

## Attribute Reference

All of the argument attributes except `filter` blocks are also exported as
result attributes. This data source will complete the data by populating
any fields that are not included in the configuration with the data for
the selected IPAM Pool CIDRs.

The following attribute is additionally exported:

* `ipamPoolCidrs` - The CIDRs provisioned into the IPAM pool, described below.

### ipam_pool_cidrs

* `cidr` - A network CIDR.
* `state` - The provisioning state of that CIDR.

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

- `read` - (Default `1m`)

<!-- cache-key: cdktf-0.20.0 input-b9fe0569a3bb59535a5bd318b59bb8f37d4a198fb5bfa62af131cfaa8ffd8980 -->
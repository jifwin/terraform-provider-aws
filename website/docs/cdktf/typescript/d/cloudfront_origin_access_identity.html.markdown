---
subcategory: "CloudFront"
layout: "aws"
page_title: "AWS: aws_cloudfront_origin_access_identity"
description: |-
  Use this data source to retrieve information for an Amazon CloudFront origin access identity.
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_cloudfront_origin_access_identity

Use this data source to retrieve information for an Amazon CloudFront origin access identity.

## Example Usage

The following example below creates a CloudFront origin access identity.

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsCloudfrontOriginAccessIdentity } from "./.gen/providers/aws/data-aws-cloudfront-origin-access-identity";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new DataAwsCloudfrontOriginAccessIdentity(this, "example", {
      id: "EDFDVBD632BHDS5",
    });
  }
}

```

## Argument Reference

* `id` (Required) -  The identifier for the distribution. For example: `EDFDVBD632BHDS5`.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `callerReference` - Internal value used by CloudFront to allow future
   updates to the origin access identity.
* `cloudfrontAccessIdentityPath` - A shortcut to the full path for the
   origin access identity to use in CloudFront, see below.
* `comment` - An optional comment for the origin access identity.
* `etag` - Current version of the origin access identity's information.
   For example: `E2QWRUHAPOMQZL`.
* `iamArn` - Pre-generated ARN for use in S3 bucket policies (see below).
   Example: `arn:aws:iam::cloudfront:user/CloudFront Origin Access Identity
   E2QWRUHAPOMQZL`.
* `s3CanonicalUserId` - The Amazon S3 canonical user ID for the origin
   access identity, which you use when giving the origin access identity read
   permission to an object in Amazon S3.

<!-- cache-key: cdktf-0.20.0 input-6ee07d67e88c83135a682bc01ecae9a6413d07992377fed93295e34cd1b0ceec -->
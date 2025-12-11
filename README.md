# ARN

A Go library for building AWS ARN (Amazon Resource Name) strings with automatic region and account ID detection.

## Features

- Automatic AWS region detection from configuration
- Automatic AWS account ID retrieval via STS
- Service-specific ARN formatting for:
  - CloudWatch Logs (with wildcard suffixes)
  - S3 (bucket and object ARNs)
  - IAM (global resources)
  - And more standard services
- Support for multiple AWS partitions (aws, aws-cn, aws-us-gov)
- Type-safe API with proper error handling

## Installation

```bash
go get github.com/JamesWoolfenden/arn
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/JamesWoolfenden/arn"
)

func main() {
    // Create an ARN for an EC2 instance
    resource := "instance/i-1234567890abcdef0"
    arnBuilder := &arn.AwsArn{
        Service:  "ec2",
        Resource: &resource,
    }

    // Builder() automatically detects region and account ID
    arns, err := arnBuilder.Builder()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(arns[0])
    // Output: arn:aws:ec2:us-east-1:123456789012:instance/i-1234567890abcdef0
}
```

### S3 Bucket ARNs

```go
bucket := "my-bucket"
arnBuilder := &arn.AwsArn{
    Service:  "s3",
    Resource: &bucket,
}

arns, err := arnBuilder.Builder()
if err != nil {
    log.Fatal(err)
}

// S3 returns both bucket and object ARNs
fmt.Println(arns[0]) // arn:aws:s3:::my-bucket
fmt.Println(arns[1]) // arn:aws:s3:::my-bucket/*
```

### CloudWatch Logs

```go
logGroup := "log-group:my-app"
arnBuilder := &arn.AwsArn{
    Service:  "logs",
    Resource: &logGroup,
}

arns, err := arnBuilder.Builder()
if err != nil {
    log.Fatal(err)
}

// CloudWatch Logs returns base ARN and wildcard variant
fmt.Println(arns[0]) // arn:aws:logs:us-east-1:123456789012:log-group:my-app
fmt.Println(arns[1]) // arn:aws:logs:us-east-1:123456789012:log-group:my-app:*
```

### Explicit Configuration

```go
partition := "aws-us-gov"
region := "us-gov-west-1"
account := "123456789012"
resource := "role/MyRole"

arnBuilder := &arn.AwsArn{
    Partition: &partition,
    Service:   "iam",
    Region:    &region,
    Account:   &account,
    Resource:  &resource,
}

arns, err := arnBuilder.Builder()
if err != nil {
    log.Fatal(err)
}

fmt.Println(arns[0])
// Output: arn:aws-us-gov:iam::123456789012:role/MyRole
```

## AWS Credentials

The library requires AWS credentials to automatically detect region and account ID. Credentials can be provided via:

- Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
- Shared credentials file (`~/.aws/credentials`)
- IAM role (when running on EC2, ECS, Lambda, etc.)

If you provide explicit `Region` and `Account` values, credentials are not required.

## API

### Type: `AwsArn`

```go
type AwsArn struct {
    Partition *string  // AWS partition (default: "aws")
    Service   string   // AWS service (e.g., "s3", "ec2", "iam")
    Region    *string  // AWS region (auto-detected if nil)
    Account   *string  // AWS account ID (auto-detected if nil)
    Resource  *string  // Resource identifier
}
```

### Method: `Builder()`

```go
func (m *AwsArn) Builder() ([]string, error)
```

Builds one or more ARN strings. Returns multiple ARNs for services that require variants (S3, CloudWatch Logs).

### Method: `GetRegion()`

```go
func (m *AwsArn) GetRegion() (*string, error)
```

Retrieves the AWS region from default configuration.

### Method: `GetAccountId()`

```go
func (m *AwsArn) GetAccountId() (*string, error)
```

Retrieves the AWS account ID using STS GetCallerIdentity.

## Testing

```bash
go test -v
```

Tests that require AWS credentials will be skipped if credentials are not configured.

## License

See LICENSE file for details.

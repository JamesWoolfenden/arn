// Package arn provides utilities for building AWS ARN (Amazon Resource Name) strings.
// It supports automatic detection of AWS region and account ID, and handles
// service-specific ARN formatting for various AWS services.
package arn

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// AwsArn represents the components of an AWS ARN.
// ARN format: arn:partition:service:region:account-id:resource-id
type AwsArn struct {
	// Partition is the AWS partition (e.g., "aws", "aws-cn", "aws-us-gov").
	// If nil, defaults to "aws".
	Partition *string
	// Service is the AWS service identifier (e.g., "s3", "iam", "ec2").
	Service string
	// Region is the AWS region (e.g., "us-east-1", "eu-west-2").
	// If nil, automatically detected from AWS config.
	Region *string
	// Account is the AWS account ID.
	// If nil, automatically detected via STS GetCallerIdentity.
	Account *string
	// Resource is the resource identifier specific to the service.
	Resource *string
}

// Builder constructs one or more ARN strings based on the configured fields.
// It automatically detects region and account ID if not provided.
// Some services (e.g., CloudWatch Logs, S3) may return multiple ARN variants.
// Returns a slice of ARN strings and an error if construction fails.
func (m *AwsArn) Builder() ([]string, error) {
	var arns []string
	var err error

	defaultPartition := "aws"

	// Initialize nil fields
	if m.Resource == nil {
		m.Resource = new(string)
	}

	if m.Partition == nil {
		m.Partition = &defaultPartition
	}

	if m.Region == nil {
		m.Region, err = m.GetRegion()
		if err != nil {
			return nil, fmt.Errorf("failed to get region: %w", err)
		}
	}

	if m.Account == nil {
		m.Account, err = m.GetAccountId()
		if err != nil {
			return nil, fmt.Errorf("failed to get account ID: %w", err)
		}
	}

	// Normalize service name to lowercase
	if m.Service != "" {
		m.Service = strings.ToLower(m.Service)
	}

	// Build base ARN using the configured partition
	baseArn := fmt.Sprintf("arn:%s:%s:%s:%s:%s",
		*m.Partition, m.Service, *m.Region, *m.Account, *m.Resource)
	arns = append(arns, baseArn)

	// Service-specific ARN handling
	switch m.Service {
	case "logs":
		// CloudWatch Logs requires wildcard suffix for log streams
		arns = append(arns, baseArn+":*")
	case "s3":
		// S3 buckets need both bucket and object ARNs
		if *m.Resource != "" {
			// Bucket ARN (no region/account for S3)
			bucketArn := fmt.Sprintf("arn:%s:s3:::%s", *m.Partition, *m.Resource)
			// Object ARN with wildcard
			objectArn := fmt.Sprintf("arn:%s:s3:::%s/*", *m.Partition, *m.Resource)
			arns = []string{bucketArn, objectArn}
		}
	case "iam":
		// IAM resources are global (no region)
		arns = []string{fmt.Sprintf("arn:%s:%s::%s:%s",
			*m.Partition, m.Service, *m.Account, *m.Resource)}
	}

	return arns, nil
}

// GetRegion retrieves the AWS region from the default AWS configuration.
// It uses the AWS SDK v2 config loader which checks environment variables,
// shared config files, and EC2 instance metadata.
func (m *AwsArn) GetRegion() (*string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed loading config: %w", err)
	}

	return &cfg.Region, nil
}

// GetAccountId retrieves the AWS account ID using STS GetCallerIdentity.
// This requires valid AWS credentials to be configured.
// Returns the account ID and an error if the API call fails.
func (m *AwsArn) GetAccountId() (*string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed loading config: %w", err)
	}

	svc := sts.NewFromConfig(cfg)
	result, err := svc.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %w", err)
	}

	return result.Account, nil
}

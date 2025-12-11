package arn

import (
	"os"
	"reflect"
	"testing"
)

// Test_awsArn_getAccountId tests the GetAccountId method.
// This test requires valid AWS credentials to be configured.
// It will be skipped in CI environments without credentials.
func Test_awsArn_getAccountId(t *testing.T) {
	// Skip if no AWS credentials are available
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" && os.Getenv("AWS_PROFILE") == "" {
		t.Skip("Skipping test: AWS credentials not configured")
	}

	m := &AwsArn{}
	got, err := m.GetAccountId()
	if err != nil {
		t.Fatalf("GetAccountId() error = %v", err)
	}
	if got == nil || *got == "" {
		t.Errorf("GetAccountId() returned empty account ID")
	}
	// Account IDs should be 12 digits
	if len(*got) != 12 {
		t.Errorf("GetAccountId() = %v, expected 12-digit account ID", *got)
	}
}

// Test_awsArn_getRegion tests the GetRegion method.
// This test requires AWS configuration to be set up.
// It will be skipped if no region is configured.
func Test_awsArn_getRegion(t *testing.T) {
	// Skip if no AWS region is configured
	if os.Getenv("AWS_REGION") == "" && os.Getenv("AWS_DEFAULT_REGION") == "" {
		t.Skip("Skipping test: AWS region not configured")
	}

	m := &AwsArn{}
	got, err := m.GetRegion()
	if err != nil {
		t.Fatalf("GetRegion() error = %v", err)
	}
	if got == nil || *got == "" {
		t.Errorf("GetRegion() returned empty region")
	}
}

// TestAwsArn_Builder tests the Builder method with various configurations.
// These tests use explicit values to avoid requiring AWS credentials.
func TestAwsArn_Builder(t *testing.T) {
	t.Parallel()
	type fields struct {
		Partition *string
		Service   string
		Region    *string
		Account   *string
		Resource  *string
	}

	empty := ""
	partition := "aws"
	region := "eu-west-2"
	account := "123456789012"
	resource := "my-resource"
	s3Bucket := "my-bucket"

	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "SSM with all fields specified",
			fields: fields{
				Partition: &partition,
				Service:   "ssm",
				Region:    &region,
				Account:   &account,
				Resource:  &resource,
			},
			want:    []string{"arn:aws:ssm:eu-west-2:123456789012:my-resource"},
			wantErr: false,
		},
		{
			name: "CloudWatch Logs with wildcards",
			fields: fields{
				Partition: &partition,
				Service:   "logs",
				Region:    &region,
				Account:   &account,
				Resource:  &resource,
			},
			want:    []string{"arn:aws:logs:eu-west-2:123456789012:my-resource", "arn:aws:logs:eu-west-2:123456789012:my-resource:*"},
			wantErr: false,
		},
		{
			name: "S3 bucket and object ARNs",
			fields: fields{
				Partition: &partition,
				Service:   "s3",
				Region:    &region,
				Account:   &account,
				Resource:  &s3Bucket,
			},
			want:    []string{"arn:aws:s3:::my-bucket", "arn:aws:s3:::my-bucket/*"},
			wantErr: false,
		},
		{
			name: "IAM resource (no region)",
			fields: fields{
				Partition: &partition,
				Service:   "iam",
				Region:    &region,
				Account:   &account,
				Resource:  &resource,
			},
			want:    []string{"arn:aws:iam::123456789012:my-resource"},
			wantErr: false,
		},
		{
			name: "Empty CloudWatch Logs",
			fields: fields{
				Partition: nil,
				Service:   "logs",
				Region:    &empty,
				Account:   &empty,
				Resource:  nil,
			},
			want:    []string{"arn:aws:logs:::", "arn:aws:logs::::*"},
			wantErr: false,
		},
		{
			name: "AWS GovCloud partition",
			fields: fields{
				Partition: strPtr("aws-us-gov"),
				Service:   "ec2",
				Region:    strPtr("us-gov-west-1"),
				Account:   &account,
				Resource:  &resource,
			},
			want:    []string{"arn:aws-us-gov:ec2:us-gov-west-1:123456789012:my-resource"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &AwsArn{
				Partition: tt.fields.Partition,
				Service:   tt.fields.Service,
				Region:    tt.fields.Region,
				Account:   tt.fields.Account,
				Resource:  tt.fields.Resource,
			}
			got, err := m.Builder()
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder() = %v, want %v", got, tt.want)
			}
		})
	}
}

// strPtr is a helper function to create string pointers for tests.
func strPtr(s string) *string {
	return &s
}

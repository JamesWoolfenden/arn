package arn

import (
	"reflect"
	"testing"
)

func Test_awsArn_builder(t *testing.T) {

	partition := "aws"
	region := "eu-west-2"
	account := "680235478471"
	type args struct {
		partition *string
		service   string
		region    *string
		account   *string
		resource  *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Pass",
			args{partition: nil, service: "ssm", region: nil, account: nil, resource: nil},
			"awsArn:aws:ssm:eu-west-2:680235478471:",
			false},
		{"Pass 2",
			args{partition: &partition, service: "ssm", region: &region, account: &account, resource: nil},
			"awsArn:aws:ssm:eu-west-2:680235478471:",
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AwsArn{}
			got, err := m.Builder(tt.args.partition, tt.args.service, tt.args.region, tt.args.account, tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("builder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("builder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_awsArn_getAccountId(t *testing.T) {

	want := "680235478471"
	tests := []struct {
		name string
		want *string
	}{
		{"Pass", &want},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AwsArn{}
			if got := m.GetAccountId(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAccountId() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func Test_awsArn_getRegion(t *testing.T) {
	region := "eu-west-2"
	tests := []struct {
		name    string
		want    *string
		wantErr bool
	}{
		{"Pass", &region, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AwsArn{}
			got, err := m.GetRegion()
			if (err != nil) != tt.wantErr {
				t.Errorf("getRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRegion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

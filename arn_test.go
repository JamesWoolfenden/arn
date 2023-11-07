package arn

import (
	"reflect"
	"testing"
)

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

func TestAwsArn_Builder(t *testing.T) {

	type fields struct {
		partition *string
		service   string
		region    *string
		account   *string
		resource  *string
	}

	partition := "aws"
	region := "eu-west-2"
	account := "680235478471"

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"Pass",
			fields{&partition, "ssm", &region, &account, nil},
			"awsArn:aws:ssm:eu-west-2:680235478471:",
			false},
		{"Pass 2",
			fields{partition: nil, service: "ssm", region: nil, account: nil, resource: nil},
			"awsArn:aws:ssm:eu-west-2:680235478471:",
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AwsArn{
				Partition: tt.fields.partition,
				Service:   tt.fields.service,
				Region:    tt.fields.region,
				Account:   tt.fields.account,
				Resource:  tt.fields.resource,
			}
			got, err := m.Builder()
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Builder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

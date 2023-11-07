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
	account := "680235478471"
	want := []string{"arn:aws:ssm:eu-west-2:680235478471:"}
	wantEmpty := []string{"arn:aws:logs:::", "arn:aws:logs::::*"}

	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{"Pass",
			fields{&partition, "ssm", &region, &account, nil},
			want,
		},
		{"Pass 2",
			fields{nil, "ssm", nil, nil, nil},
			want,
		},
		{"Pass 3",
			fields{nil, "logs", &empty, &empty, nil},
			wantEmpty,
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
			if got := m.Builder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Builder() = %v, want %v", got, tt.want)
			}
		})
	}
}

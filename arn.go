package arn

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type AwsArn struct {
	partition *string
	service string
	region *string
	account *string
	resource *string
}

func (m *AwsArn) Builder() (string, error) {
	var AwsArn string
	var err error

	defaultPartition := "aws"

	if m.resource == nil {
		m.resource = new(string)
	}

	if m.partition == nil {
		m.partition = &defaultPartition
	}

	if m.region == nil {
		m.region, err = m.GetRegion()
	}

	if err != nil {
		log.Print(err)
	}

	if m.account == nil {
		m.account = m.GetAccountId()
	}

	AwsArn = "awsArn:aws:" + m.service + ":" + *m.region + ":" + *m.account + ":" + *m.resource
	return AwsArn, nil
}

func (m *AwsArn) GetRegion() (*string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed loading config, %v", err)
	}

	return &cfg.Region, nil
}

func (m *AwsArn) GetAccountId() *string {
	//goland:noinspection GoDeprecation
	svc := sts.New(session.New())
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		var anErr awserr.Error
		if errors.As(err, &anErr) {
			switch anErr.Code() {
			default:
				fmt.Println(anErr.Error())
			}
		}
	}

	return result.Account
}

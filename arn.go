package arn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type AwsArn struct {
	Partition *string
	Service   string
	Region    *string
	Account   *string
	Resource  *string
}

func (m *AwsArn) Builder() []string {
	var AwsArn []string
	var err error

	defaultPartition := "aws"

	if m.Resource == nil {
		m.Resource = new(string)
	}

	if m.Partition == nil {
		m.Partition = &defaultPartition
	}

	if m.Region == nil {
		m.Region, err = m.GetRegion()
		if err != nil {
			log.Printf("region %s", err)
		}
	}

	if m.Account == nil {
		m.Account = m.GetAccountId()
	}

	if m.Service != "" {
		m.Service = strings.ToLower(m.Service)
	}

	AwsArn = append(AwsArn, "arn:aws:"+m.Service+":"+*m.Region+":"+*m.Account+":"+*m.Resource)

	switch m.Service {
	case "logs":
		AwsArn = append(AwsArn, "arn:aws:"+m.Service+":"+*m.Region+":"+*m.Account+":"+*m.Resource+":*")
	}

	//todo is this a service that needs more than one arn in its resource string e.g. s3

	return AwsArn
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

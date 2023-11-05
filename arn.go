package arnBuilder

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

type awsArn struct {
}

func (m *awsArn) builder(partition *string, service string, region *string, account *string, resource *string) (string, error) {
	var awsArn string
	var err error

	defaultPartition := "aws"

	if resource == nil {
		resource = new(string)
	}

	if partition == nil {
		partition = &defaultPartition
	}

	if region == nil {
		region, err = m.getRegion()
	}

	if err != nil {
		log.Print(err)
	}

	if account == nil {
		account = m.getAccountId()
	}

	awsArn = "awsArn:aws:" + service + ":" + *region + ":" + *account + ":" + *resource
	return awsArn, nil
}

func (m *awsArn) getRegion() (*string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed loading config, %v", err)
	}

	return &cfg.Region, nil
}

func (m *awsArn) getAccountId() *string {
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

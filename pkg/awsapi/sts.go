package awsapi

//go:generate mockgen -destination=mocks/mock_stsiface.go -package=mock_awsapi github.com/aws/aws-sdk-go/service/sts/stsiface STSAPI
//go:generate mockgen -destination=mocks/mock_sts.go -package=mock_awsapi github.com/keikoproj/iam-manager/pkg/awsapi STSIface

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/keikoproj/iam-manager/pkg/log"
)

type STSIface interface {
	GetAccountID(ctx context.Context) (string, error)
}

type STS struct {
	Client stsiface.STSAPI
}

func NewSTS() *STS {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-west-2")})
	if err != nil {
		panic(err)
	}
	return &STS{
		Client: sts.New(sess),
	}
}

// GetAccountID loads aws accountID from sts caller identity
func (i *STS) GetAccountID(ctx context.Context) (string, error) {
	log := log.Logger(context.Background(), "awsapi", "iam", "GetAccountID")

	// get caller identity in order to fetch aws account ID
	result, err := i.Client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Error(err, "failed to get caller's identity")
		return "", err
	}
	return *result.Account, nil
}

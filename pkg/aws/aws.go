package aws_services

import (
	pkg_constants "ResiSync/pkg/constants"
	"ResiSync/pkg/models"
	"ResiSync/pkg/security"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func CreateNewS3Session() (*s3.S3, error) {

	var awsModel models.Aws
	err := viper.UnmarshalKey(pkg_constants.ConfigSectionAWS, &awsModel)
	if err != nil {
		log.Println("Error while unmarshalling aws config", err)
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(awsModel); err != nil {
		log.Println("Error while validating aws config", err)
		return nil, err
	}

	accessKey, err := security.DecryptPassword(awsModel.EncryptedSecretKey, awsModel.SecretKeyNonce)
	if err != nil {
		log.Println("Error while decrypting aws secret", err)
		return nil, err
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsModel.Region),
		Credentials: credentials.NewStaticCredentials(awsModel.AccessKeyId, accessKey, awsModel.Token),
	})
	if err != nil {
		log.Println("Error while creating aws session", err)
		return nil, err
	}

	return s3.New(sess), nil
}

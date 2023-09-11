package shared_utils

import (
	"ResiSync/pkg/api"
	pkg_models "ResiSync/pkg/models"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nyaruka/phonenumbers"
	"go.uber.org/zap"
)

func NowInUTC() time.Time {
	return time.Now().UTC()
}

func IsValidEmail(email string) bool {

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)

}

func IsValidContact(contact, region string) bool {
	number, err := phonenumbers.Parse(contact, "IN")
	if err != nil {
		fmt.Println(err)
		return false
	}

	return phonenumbers.IsValidNumber(number)
}

func GetPresignedS3Url(requestContext pkg_models.ResiSyncRequestContext, bucket, key string, duration time.Duration) string {
	span := api.AddTrace(&requestContext, "info", "GetPresignedS3Url")
	defer span.End()

	s3Session := api.ApplicationContext.S3Session

	req, _ := s3Session.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(duration)
	if err != nil {
		requestContext.Log.Error("Failed to sign object",
			zap.String("bucket", bucket), zap.String("key", key), zap.Error(err))
	}

	return urlStr
}

func DeleteObjectFromS3(requestContext pkg_models.ResiSyncRequestContext, bucket, key string) error {
	span := api.AddTrace(&requestContext, "info", "DeleteObjectFromS3")
	defer span.End()

	s3Session := api.ApplicationContext.S3Session

	_, err := s3Session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		requestContext.Log.Error("Failed to delete object request",
			zap.String("bucket", bucket), zap.String("key", key), zap.Error(err))
		return err
	}

	return nil
}

func GenerateOTP() string {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}
	return n.String()
}

package shared_utils

import (
	"fmt"
	"regexp"
	"time"

	"github.com/nyaruka/phonenumbers"
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

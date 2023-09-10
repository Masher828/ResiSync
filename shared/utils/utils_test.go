package shared_utils_test

import (
	"ResiSync/pkg/config"
	shared_utils "ResiSync/shared/utils"
	"fmt"
	"testing"
)

func init() {
	config.LoadEnv()

	config.LoadConfig()
}

func Test_IsValidContact(t *testing.T) {
	number := "xxxxxxxx"
	fmt.Println(shared_utils.IsValidContact(number, "in"))
}

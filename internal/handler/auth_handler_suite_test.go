// internal/handler/auth_handler_suite_test.go
package handler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuthHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthHandler Suite")
}

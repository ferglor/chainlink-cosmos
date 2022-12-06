package migration_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"

	. "github.com/onsi/ginkgo/v2"
)

func Test_Suite(t *testing.T) {
	actions.GinkgoSuite()
	RunSpecs(t, "Migration")
}

package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateIssueAction(t *testing.T) {
	t.Run("fails on nil action", func(t *testing.T) {
		_, _, _, _, err := validateIssueAction(nil)
		assert.Error(t, err)
	})
}

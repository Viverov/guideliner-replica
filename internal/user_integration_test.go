// +build integration

package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegrationSuccess(t *testing.T) {
	assert.Equal(t, 1, 1)
}

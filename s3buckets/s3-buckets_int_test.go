package s3buckets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoesExistIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	execReqId := "non-existing-exec-request-id"
	actual, _ := DoesExist(context.Background(), execReqId)
	assert.Equal(t, false, actual, "should not exist")
}

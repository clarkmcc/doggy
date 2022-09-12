package doggy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetParentFunc(t *testing.T) {
	assert.Equal(t, "github.com/clarkmcc/doggy.sampleFunc", sampleFunc())
}

func sampleFunc() string {
	return getFunctionName()
}

func getFunctionName() string {
	_, _, name := getParentFunc(1)
	return name
}

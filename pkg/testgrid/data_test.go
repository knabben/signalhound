package testgrid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RenderStatuses(t *testing.T) {
	test := Test{
		ShortTexts: []string{
			"F",
			"T",
		},
		Messages: []string{
			"Build failed outside of test results",
			"Build did not complete within 24 hours",
		},
	}
	text, failures, firstFailure := test.RenderStatuses([]int64{
		1730221517000, 1728666129000,
	})
	assert.Len(t, text, 142)
	assert.Equal(t, 2, failures)
	assert.Equal(t, firstFailure, 0)
}

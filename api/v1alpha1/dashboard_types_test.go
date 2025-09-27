package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderStatuses(t *testing.T) {
	message := "kubetest --timeout triggered"
	tests := []struct {
		name            string
		inputTest       Test
		inputTimestamps []int64
		expectedOutput  string
		expectedCount   int
		expectedIndex   int
	}{
		{
			name: "all short texts must match timestamp",
			inputTest: Test{
				ShortTexts: []string{"", "", "F", "", "F"},
				Messages:   []string{"", "", message, "", message},
			},
			inputTimestamps: []int64{1758974631000, 1758967371000, 1758960111000, 1758952851000, 1758945591000},
			expectedOutput:  formatTestStatus("F", 1758960111000, message) + formatTestStatus("F", 1758945591000, message),
			expectedIndex:   2,
			expectedCount:   2,
		},
		{
			name: "no statuses to render",
			inputTest: Test{
				ShortTexts: []string{"", "", ""},
				Messages:   []string{"", "", ""},
			},
			inputTimestamps: []int64{1620000000, 1620003600, 1620007200},
			expectedOutput:  "",
			expectedCount:   0,
			expectedIndex:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, failureCount, firstFailureIndex := tt.inputTest.RenderStatuses(tt.inputTimestamps)
			assert.Equal(t, tt.expectedOutput, output)
			assert.Equal(t, tt.expectedCount, failureCount)
			assert.Equal(t, tt.expectedIndex, firstFailureIndex)
		})
	}
}

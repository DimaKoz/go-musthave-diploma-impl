package cooldown

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsAccrualReady(t *testing.T) {
	tests := []struct {
		name          string
		requestedTime int64
		sleepSecTime  int64
		want          bool
	}{
		{
			name:          "zero",
			requestedTime: 0,
			sleepSecTime:  0,
			want:          true,
		},

		{
			name:          "ready",
			requestedTime: 1,
			sleepSecTime:  2,
			want:          true,
		},
		{
			name:          "notReady",
			requestedTime: 2,
			sleepSecTime:  1,
			want:          false,
		},
	}
	for _, testCase := range tests {
		test := testCase
		t.Run(test.name, func(t *testing.T) {
			if test.requestedTime != 0 {
				NeedAccrualCooldown(test.requestedTime)
			}
			if test.sleepSecTime != 0 {
				time.Sleep(time.Duration(test.sleepSecTime) * time.Second)
			}
			got := IsAccrualReady()
			assert.Equal(t, test.want, got)
		})
	}
}

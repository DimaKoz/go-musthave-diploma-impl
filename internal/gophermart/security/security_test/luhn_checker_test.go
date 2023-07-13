package security_test

import (
	"testing"

	"github.com/DimaKoz/go-musthave-diploma-impl/internal/gophermart/security"
	"github.com/stretchr/testify/assert"
)

func TestIsValidLuhnNumber(t *testing.T) {
	type args struct {
		luhnNumber string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "bad luhn number",
			args: args{
				luhnNumber: "465530212",
			},
			want: false,
		},
		{
			name: "good luhn number",
			args: args{
				luhnNumber: "79927398713",
			},
			want: true,
		},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			got := security.IsValidLuhnNumber(tt.args.luhnNumber)
			if tt.want {
				assert.True(t, got)
			} else {
				assert.False(t, got)
			}
		})
	}
}

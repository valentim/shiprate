package shipping

import "testing"

func TestMoneyString(t *testing.T) {
	tests := []struct {
		amount Money
		want   string
	}{
		{0, "0.00"},
		{5, "0.05"},
		{EUR(18), "18.00"},
		{1805, "18.05"},
		{-250, "-2.50"},
	}
	for _, tt := range tests {
		if got := tt.amount.String(); got != tt.want {
			t.Errorf("Money(%d).String() = %q, want %q", int64(tt.amount), got, tt.want)
		}
	}
}

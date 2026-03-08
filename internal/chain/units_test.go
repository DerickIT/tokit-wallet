package chain

import (
	"math/big"
	"testing"
)

func TestParseUnits(t *testing.T) {
	tests := []struct {
		name     string
		amount   string
		decimals uint8
		want     string
		wantErr  bool
	}{
		{name: "whole amount", amount: "1", decimals: 18, want: "1000000000000000000"},
		{name: "fractional amount", amount: "0.5", decimals: 18, want: "500000000000000000"},
		{name: "trim spaces", amount: " 1.25 ", decimals: 6, want: "1250000"},
		{name: "dot prefix", amount: ".5", decimals: 6, want: "500000"},
		{name: "too many decimals", amount: "1.1234567", decimals: 6, wantErr: true},
		{name: "negative amount", amount: "-1", decimals: 18, wantErr: true},
		{name: "invalid text", amount: "abc", decimals: 18, wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseUnits(tc.amount, tc.decimals)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tc.want {
				t.Fatalf("unexpected value: got %s, want %s", got.String(), tc.want)
			}
		})
	}
}

func TestFormatUnits(t *testing.T) {
	tests := []struct {
		name      string
		value     *big.Int
		decimals  uint8
		precision int
		want      string
	}{
		{
			name:      "nil value",
			value:     nil,
			decimals:  18,
			precision: 6,
			want:      "0",
		},
		{
			name:      "whole value",
			value:     big.NewInt(1230000000000000000),
			decimals:  18,
			precision: 6,
			want:      "1.23",
		},
		{
			name:      "precision trimming",
			value:     big.NewInt(123456789),
			decimals:  8,
			precision: 4,
			want:      "1.2345",
		},
		{
			name:      "no decimals",
			value:     big.NewInt(42),
			decimals:  0,
			precision: 6,
			want:      "42",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := FormatUnits(tc.value, tc.decimals, tc.precision)
			if got != tc.want {
				t.Fatalf("unexpected value: got %s, want %s", got, tc.want)
			}
		})
	}
}

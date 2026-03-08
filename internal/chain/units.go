package chain

import (
	"fmt"
	"math/big"
	"strings"
)

func ParseUnits(amount string, decimals uint8) (*big.Int, error) {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return nil, fmt.Errorf("amount is required")
	}
	if strings.HasPrefix(amount, "-") {
		return nil, fmt.Errorf("amount must be positive")
	}

	parts := strings.Split(amount, ".")
	if len(parts) > 2 {
		return nil, fmt.Errorf("invalid amount: %s", amount)
	}

	wholePart := parts[0]
	if wholePart == "" {
		wholePart = "0"
	}
	if !isDigits(wholePart) {
		return nil, fmt.Errorf("invalid amount: %s", amount)
	}

	fractionalPart := ""
	if len(parts) == 2 {
		fractionalPart = parts[1]
		if fractionalPart == "" {
			return nil, fmt.Errorf("invalid amount: %s", amount)
		}
		if !isDigits(fractionalPart) {
			return nil, fmt.Errorf("invalid amount: %s", amount)
		}
		if len(fractionalPart) > int(decimals) {
			return nil, fmt.Errorf("amount has too many decimal places; max supported is %d", decimals)
		}
	}

	base := pow10(decimals)
	wholeValue, ok := new(big.Int).SetString(wholePart, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount: %s", amount)
	}
	wholeValue.Mul(wholeValue, base)

	if fractionalPart == "" {
		return wholeValue, nil
	}

	fractionalPart += strings.Repeat("0", int(decimals)-len(fractionalPart))
	fractionalValue, ok := new(big.Int).SetString(fractionalPart, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount: %s", amount)
	}

	return wholeValue.Add(wholeValue, fractionalValue), nil
}

func FormatUnits(value *big.Int, decimals uint8, precision int) string {
	if value == nil {
		return "0"
	}

	if decimals == 0 {
		return value.String()
	}

	base := pow10(decimals)
	whole := new(big.Int)
	fractional := new(big.Int)
	whole.QuoRem(value, base, fractional)

	fractionalText := fractional.Text(10)
	if len(fractionalText) < int(decimals) {
		fractionalText = strings.Repeat("0", int(decimals)-len(fractionalText)) + fractionalText
	}
	fractionalText = strings.TrimRight(fractionalText, "0")
	if precision > 0 && len(fractionalText) > precision {
		fractionalText = strings.TrimRight(fractionalText[:precision], "0")
	}

	if fractionalText == "" {
		return whole.String()
	}

	return fmt.Sprintf("%s.%s", whole.String(), fractionalText)
}

func pow10(decimals uint8) *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
}

func isDigits(value string) bool {
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

package services

import (
	"encoding/json"
	"fmt"
	"sort"
)

func ParseJSONStringArray(raw []byte) ([]string, error) {
	if len(raw) == 0 {
		return []string{}, nil
	}

	var values []string
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, fmt.Errorf("parse json string array: %w", err)
	}

	return values, nil
}

func NormalizeOptions(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" {
			normalized = append(normalized, value)
		}
	}

	sort.Strings(normalized)
	return normalized
}

func CalculateAwardedMarks(chosenOptions, correctOptions []string, marks int) int {
	normalizedChosen := NormalizeOptions(chosenOptions)
	normalizedCorrect := NormalizeOptions(correctOptions)

	if len(normalizedChosen) != len(normalizedCorrect) {
		return 0
	}

	for index := range normalizedChosen {
		if normalizedChosen[index] != normalizedCorrect[index] {
			return 0
		}
	}

	return marks
}

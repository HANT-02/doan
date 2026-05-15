package skills

import "strings"

func NormalizeCodes(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	for _, value := range values {
		code := strings.ToUpper(strings.TrimSpace(value))
		if code == "" {
			continue
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		normalized = append(normalized, code)
	}

	return normalized
}

func MissingRequiredCodes(have []string, required []string) []string {
	requiredNormalized := NormalizeCodes(required)
	if len(requiredNormalized) == 0 {
		return []string{}
	}

	haveSet := make(map[string]struct{}, len(have))
	for _, value := range NormalizeCodes(have) {
		haveSet[value] = struct{}{}
	}

	missing := make([]string, 0, len(requiredNormalized))
	for _, code := range requiredNormalized {
		if _, ok := haveSet[code]; ok {
			continue
		}
		missing = append(missing, code)
	}

	return missing
}

func HumanizeCode(code string) string {
	normalized := strings.TrimSpace(code)
	if normalized == "" {
		return ""
	}
	return strings.ReplaceAll(normalized, "_", " ")
}

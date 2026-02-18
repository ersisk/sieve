package search

import (
	"regexp"
	"strings"

	"github.com/ersanisk/sieve/pkg/logentry"
)

// RegexMatch performs regex matching against log entries.
// Returns matches sorted by the number of field matches.
func RegexMatch(entries []logentry.Entry, pattern string) ([]SearchResult, error) {
	if pattern == "" {
		return nil, nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var results []SearchResult

	for _, entry := range entries {
		matchedFields := regexSearchEntry(entry, re)
		if len(matchedFields) > 0 {
			score := float64(len(matchedFields)) / 5.0
			if score > 1.0 {
				score = 1.0
			}

			results = append(results, SearchResult{
				Entry:   entry,
				Score:   score,
				Matched: matchedFields,
			})
		}
	}

	return sortByScore(results), nil
}

func regexSearchEntry(entry logentry.Entry, re *regexp.Regexp) []string {
	var matched []string

	if re.MatchString(entry.Message) {
		matched = append(matched, "message: "+highlightMatch(entry.Message, re))
	}

	if entry.Caller != "" && re.MatchString(entry.Caller) {
		matched = append(matched, "caller: "+highlightMatch(entry.Caller, re))
	}

	for key, value := range entry.Fields {
		strVal, ok := valueToString(value)
		if ok && re.MatchString(strVal) {
			matched = append(matched, key+": "+highlightMatch(strVal, re))
		}
	}

	return matched
}

// highlightMatch wraps regex matches with markers.
func highlightMatch(text string, re *regexp.Regexp) string {
	matches := re.FindAllStringIndex(text, -1)
	if len(matches) == 0 {
		return text
	}

	var result strings.Builder
	lastEnd := 0

	for _, match := range matches {
		start, end := match[0], match[1]
		result.WriteString(text[lastEnd:start])
		result.WriteString("[")
		result.WriteString(text[start:end])
		result.WriteString("]")
		lastEnd = end
	}

	result.WriteString(text[lastEnd:])
	return result.String()
}

// RegexCaseInsensitiveMatch performs case-insensitive regex matching.
func RegexCaseInsensitiveMatch(entries []logentry.Entry, pattern string) ([]SearchResult, error) {
	if pattern == "" {
		return nil, nil
	}

	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return nil, err
	}

	return RegexMatch(entries, re.String())
}

// RegexFieldMatch performs regex matching on a specific field.
func RegexFieldMatch(entries []logentry.Entry, fieldName, pattern string) ([]SearchResult, error) {
	if pattern == "" {
		return nil, nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var results []SearchResult

	for _, entry := range entries {
		var value string
		var found bool

		switch fieldName {
		case "message", "msg":
			value = entry.Message
			found = re.MatchString(value)
		case "caller", "source":
			value = entry.Caller
			found = re.MatchString(value)
		default:
			if val, ok := entry.GetField(fieldName); ok {
				value, ok = valueToString(val)
				if ok {
					found = re.MatchString(value)
				}
			}
		}

		if found {
			results = append(results, SearchResult{
				Entry:   entry,
				Score:   1.0,
				Matched: []string{fieldName + ": " + highlightMatch(value, re)},
			})
		}
	}

	return results, nil
}

// RegexMultiMatch performs regex matching with multiple patterns (OR logic).
func RegexMultiMatch(entries []logentry.Entry, patterns []string) ([]SearchResult, error) {
	if len(patterns) == 0 {
		return nil, nil
	}

	regexes := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		regexes = append(regexes, re)
	}

	var results []SearchResult

	for _, entry := range entries {
		allMatched := make(map[string]string)

		for _, re := range regexes {
			matchedFields := regexSearchEntry(entry, re)
			for _, field := range matchedFields {
				allMatched[field] = field
			}
		}

		if len(allMatched) > 0 {
			matchedFields := make([]string, 0, len(allMatched))
			for field := range allMatched {
				matchedFields = append(matchedFields, field)
			}

			score := float64(len(matchedFields)) / 10.0
			if score > 1.0 {
				score = 1.0
			}

			results = append(results, SearchResult{
				Entry:   entry,
				Score:   score,
				Matched: matchedFields,
			})
		}
	}

	return sortByScore(results), nil
}

// RegexExcludeMatch performs regex matching and excludes entries that match the pattern.
func RegexExcludeMatch(entries []logentry.Entry, pattern string) ([]SearchResult, error) {
	if pattern == "" {
		return nil, nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var results []SearchResult

	for _, entry := range entries {
		matchedFields := regexSearchEntry(entry, re)
		if len(matchedFields) == 0 {
			results = append(results, SearchResult{
				Entry:   entry,
				Score:   1.0,
				Matched: []string{},
			})
		}
	}

	return results, nil
}

// RegexAndMatch performs regex matching with multiple patterns (AND logic).
func RegexAndMatch(entries []logentry.Entry, patterns []string) ([]SearchResult, error) {
	if len(patterns) == 0 {
		return nil, nil
	}

	regexes := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		regexes = append(regexes, re)
	}

	var results []SearchResult

	for _, entry := range entries {
		allPatternsMatched := true
		allMatched := make(map[string]string)

		for _, re := range regexes {
			matchedFields := regexSearchEntry(entry, re)
			if len(matchedFields) == 0 {
				allPatternsMatched = false
				break
			}
			for _, field := range matchedFields {
				allMatched[field] = field
			}
		}

		if allPatternsMatched {
			matchedFields := make([]string, 0, len(allMatched))
			for field := range allMatched {
				matchedFields = append(matchedFields, field)
			}

			results = append(results, SearchResult{
				Entry:   entry,
				Score:   1.0,
				Matched: matchedFields,
			})
		}
	}

	return sortByScore(results), nil
}

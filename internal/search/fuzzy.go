package search

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ersanisk/sieve/pkg/logentry"
)

// SearchResult represents a search match in a log entry.
type SearchResult struct {
	Entry    logentry.Entry
	Position int
	Matched  []string // matched field values
	Score    float64  // match score (0-1, higher is better)
}

// FuzzyMatch performs fuzzy matching against log entries.
// Returns matches sorted by relevance score.
func FuzzyMatch(entries []logentry.Entry, query string) []SearchResult {
	if query == "" {
		return nil
	}

	query = strings.ToLower(query)
	queryRunes := []rune(query)

	var results []SearchResult

	for _, entry := range entries {
		score, matched := fuzzySearchEntry(entry, query, queryRunes)
		if score > 0 {
			results = append(results, SearchResult{
				Entry:   entry,
				Score:   score,
				Matched: matched,
			})
		}
	}

	return sortByScore(results)
}

func fuzzySearchEntry(entry logentry.Entry, query string, queryRunes []rune) (float64, []string) {
	var maxScore float64
	var bestMatched []string

	messageScore, msgMatched := fuzzyStringMatch(entry.Message, query, queryRunes)
	if messageScore > maxScore {
		maxScore = messageScore
		bestMatched = msgMatched
	}

	if entry.Caller != "" {
		callerScore, callerMatched := fuzzyStringMatch(entry.Caller, query, queryRunes)
		if callerScore > maxScore {
			maxScore = callerScore
			bestMatched = callerMatched
		}
	}

	for key, value := range entry.Fields {
		var strVal string
		switch v := value.(type) {
		case string:
			strVal = v
		case float64:
			strVal = formatFloat(v)
		case int:
			strVal = formatInt(v)
		case bool:
			strVal = formatBool(v)
		default:
			strVal = ""
		}

		if strVal != "" {
			fieldScore, fieldMatched := fuzzyStringMatch(strVal, query, queryRunes)
			if fieldScore > maxScore {
				maxScore = fieldScore
				bestMatched = []string{key + ": " + fieldMatched[0]}
			}
		}
	}

	return maxScore, bestMatched
}

func fuzzyStringMatch(text, query string, queryRunes []rune) (float64, []string) {
	if text == "" {
		return 0, nil
	}

	textLower := strings.ToLower(text)
	textRunes := []rune(textLower)

	if strings.Contains(textLower, query) {
		return 1.0, []string{text}
	}

	if queryRunes[0] == textRunes[0] {
		return fuzzyMatchRunes(textRunes, queryRunes)
	}

	return 0, nil
}

func fuzzyMatchRunes(text, query []rune) (float64, []string) {
	if len(query) > len(text) {
		return 0, nil
	}

	tI := 0
	qI := 0
	matched := 0
	matchedPositions := make([]int, 0, len(query))

	for tI < len(text) && qI < len(query) {
		if text[tI] == query[qI] {
			_ = append(matchedPositions, tI)
			matched++
			qI++
		}
		tI++
	}

	if matched != len(query) {
		return 0, nil
	}

	score := float64(matched) / float64(len(text))
	return score * 0.8, []string{string(text)}
}

func sortByScore(results []SearchResult) []SearchResult {
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	return results
}

func formatFloat(f float64) string {
	return strings.TrimRight(strings.TrimRight(strings.ReplaceAll(
		strings.ReplaceAll(
			strconv.FormatFloat(f, 'f', -1, 64),
			".0", ""),
		"-", ""),
		"0"), ".")
}

func formatInt(i int) string {
	var b []byte
	uid := i
	neg := i < 0
	if neg {
		uid = -uid
	}

	for uid > 0 {
		b = append([]byte{byte('0' + uid%10)}, b...)
		uid /= 10
	}

	if neg {
		b = append([]byte{'-'}, b...)
	}

	if len(b) == 0 {
		b = []byte{'0'}
	}

	return string(b)
}

func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// SmartMatch combines fuzzy matching with case-insensitive exact matching.
func SmartMatch(entries []logentry.Entry, query string) []SearchResult {
	if query == "" {
		return nil
	}

	query = strings.ToLower(query)
	var results []SearchResult

	for _, entry := range entries {
		matchedFields := smartSearchEntry(entry, query)
		if len(matchedFields) > 0 {
			score := calculateRelevanceScore(entry, matchedFields)
			results = append(results, SearchResult{
				Entry:   entry,
				Score:   score,
				Matched: matchedFields,
			})
		}
	}

	return sortByScore(results)
}

func smartSearchEntry(entry logentry.Entry, query string) []string {
	var matched []string

	if strings.Contains(strings.ToLower(entry.Message), query) {
		matched = append(matched, "message: "+entry.Message)
	}

	if entry.Caller != "" && strings.Contains(strings.ToLower(entry.Caller), query) {
		matched = append(matched, "caller: "+entry.Caller)
	}

	for key, value := range entry.Fields {
		strVal, ok := valueToString(value)
		if ok && strings.Contains(strings.ToLower(strVal), query) {
			matched = append(matched, key+": "+strVal)
		}
	}

	return matched
}

func valueToString(value any) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	case float64:
		return strings.TrimRight(strings.TrimRight(strings.ReplaceAll(
			strings.ReplaceAll(
				strconv.FormatFloat(v, 'f', -1, 64),
				".0", ""),
			"-", ""),
			"0"), "."), true
	case int:
		return formatInt(v), true
	case bool:
		return formatBool(v), true
	default:
		return "", false
	}
}

func calculateRelevanceScore(entry logentry.Entry, matched []string) float64 {
	var score float64

	for _, m := range matched {
		if strings.HasPrefix(m, "message: ") {
			score += 0.5
		} else if strings.HasPrefix(m, "caller: ") {
			score += 0.3
		} else {
			score += 0.1
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// TokenizeQuery splits a search query into tokens.
func TokenizeQuery(query string) []string {
	var tokens []string
	var current strings.Builder

	query = strings.TrimSpace(query)

	for i := 0; i < len(query); i++ {
		r, size := utf8.DecodeRuneInString(query[i:])
		i += size - 1

		if unicode.IsSpace(r) {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else if r == '"' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

			i++
			for i < len(query) && query[i] != '"' {
				current.WriteByte(query[i])
				i++
			}

			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

package parser

import (
	"strings"
	"testing"
)

func BenchmarkParseLine(b *testing.B) {
	p := NewParser()
	line := `{"level":"info","msg":"Server starting","ts":"2024-01-15T10:00:00Z","service":"api","port":8080,"request_id":"req-123","duration_ms":45}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.ParseLine(line, i+1)
	}
}

func BenchmarkParseLine_Simple(b *testing.B) {
	p := NewParser()
	line := `{"level":"info","msg":"test"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.ParseLine(line, i+1)
	}
}

func BenchmarkParseLine_Complex(b *testing.B) {
	p := NewParser()
	line := `{"level":"warn","msg":"Complex log entry","ts":"2024-01-15T10:00:00.123456789Z","service":"api","request_id":"req-123-abc-456","user_id":"usr_789","method":"POST","path":"/api/v1/users","status":201,"duration_ms":23.45,"metadata":{"client_ip":"192.168.1.100","user_agent":"Mozilla/5.0","referer":"https://example.com"},"tags":["important","authenticated"],"nested":{"deeply":{"value":42}}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.ParseLine(line, i+1)
	}
}

func BenchmarkParseLine_NestedJSON(b *testing.B) {
	p := NewParser()
	line := `{"level":"info","msg":"nested data","data":{"users":[{"id":1,"name":"Alice","email":"alice@example.com","roles":["admin","user"]},{"id":2,"name":"Bob","email":"bob@example.com","roles":["user"]}],"total":2,"page":1}}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.ParseLine(line, i+1)
	}
}

func BenchmarkParseLine_PlainText(b *testing.B) {
	p := NewParser()
	line := `2024-01-15 10:00:00 INFO This is a plain text log line with no JSON structure`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.ParseLine(line, i+1)
	}
}

func BenchmarkParseLines(b *testing.B) {
	p := NewParser()

	input := strings.Repeat(`{"level":"info","msg":"test","ts":"2024-01-15T10:00:00Z"}
`, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(input)
		_, _ = p.ParseLines(r)
	}
}

func BenchmarkDetectFormat(b *testing.B) {
	input := strings.Repeat(`{"level":"info","msg":"test"}
`, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(input)
		DetectFormat(r)
	}
}

func BenchmarkDetectFormat_Mixed(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 50; i++ {
		sb.WriteString(`{"level":"info","msg":"json line"}
`)
		sb.WriteString(`plain text line
`)
	}

	input := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(input)
		DetectFormat(r)
	}
}

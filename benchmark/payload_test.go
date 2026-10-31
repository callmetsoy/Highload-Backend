package benchmark

import (
	"encoding/json"
	"testing"

	"google.golang.org/protobuf/proto"

	"sis2/gen"
)

func BenchmarkPayloadSize(b *testing.B) {
	restReq := map[string]string{
		"long_url":        "https://example.com/very/long/path?param=1",
		"idempotency_key": "6f1c1234-5678-90ab-cdef-1234567890ab",
	}
	protoReq := &gen.CreateUrlRequest{
		LongUrl:        "https://example.com/very/long/path?param=1",
		IdempotencyKey: "6f1c1234-5678-90ab-cdef-1234567890ab",
	}

	b.Run("JSON", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = json.Marshal(restReq)
		}
	})

	b.Run("Protobuf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = proto.Marshal(protoReq)
		}
	})
}

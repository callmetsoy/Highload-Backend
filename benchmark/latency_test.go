package benchmark

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"sis2/gen"
	"sis2/internal/grpcsvc"
)

func BenchmarkLatencyGRPC(b *testing.B) {
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close()

	client := grpcsvc.NewClient(gen.NewUrlServiceClient(conn))
	hist := hdrhistogram.New(1, 10_000_000, 3)
	ctx := context.Background()

	for i := 0; i < 200; i++ {
		_, _ = client.CreateWithRetry(ctx, &gen.CreateUrlRequest{
			LongUrl:        "https://example.com/warmup",
			IdempotencyKey: fmt.Sprintf("warm-%d", i),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &gen.CreateUrlRequest{
			LongUrl:        "https://example.com/very/long",
			IdempotencyKey: fmt.Sprintf("key-%d", time.Now().UnixNano()),
		}
		start := time.Now()
		_, err := client.CreateWithRetry(ctx, req)
		elapsed := time.Since(start).Microseconds()
		if err == nil {
			hist.RecordValue(elapsed)
		}
	}

	b.ReportMetric(float64(hist.ValueAtQuantile(50)), "p50_μs")
	b.ReportMetric(float64(hist.ValueAtQuantile(95)), "p95_μs")
	b.ReportMetric(float64(hist.ValueAtQuantile(99)), "p99_μs")
}

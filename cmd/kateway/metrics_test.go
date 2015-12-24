package main

import (
	"testing"
	"time"

	"github.com/funkygao/gafka/ctx"
)

func init() {
	ctx.LoadConfig("/etc/kateway.cf")
}

func BenchmarkMetricsCounterWithoutLock(b *testing.B) {
	s := NewServerMetrics(time.Hour)
	for i := 0; i < b.N; i++ {
		s.TotalConns.Inc(1)
	}
}

func BenchmarkMetricsPubOkCounter(b *testing.B) {
	p := NewPubMetrics()
	for i := 0; i < b.N; i++ {
		p.pubOk("appid", "topic", "ver")
	}
}

func BenchmarkMetricsQpsMeter(b *testing.B) {
	p := NewPubMetrics()
	for i := 0; i < b.N; i++ {
		p.PubQps.Mark(1)
	}
}

func BenchmarkMetricsLatencyHistogram(b *testing.B) {
	p := NewPubMetrics()
	for i := 0; i < b.N; i++ {
		p.PubLatency.Update(5)
	}
}

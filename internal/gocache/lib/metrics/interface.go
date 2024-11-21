package metrics

import "github.com/tianlin0/plat-lib/internal/gocache/lib/codec"

// MetricsInterface represents the metrics interface for all available providers
type MetricsInterface interface {
	RecordFromCodec(codec codec.CodecInterface)
}

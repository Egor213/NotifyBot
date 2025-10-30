package worker

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Worker interface {
	ProcessMsg(ctx context.Context, msg kafka.Message)
}

package broker

import "context"

type Consumer interface {
	Run(ctx context.Context) error
	Close() error
}

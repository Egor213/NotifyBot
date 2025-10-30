package broker

import "context"

type Broker interface {
	Run(ctx context.Context) error
	Close() error
}

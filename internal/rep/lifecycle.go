package rep

import "context"

// Lifecycle incidates that struct can be ran as an application component
type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

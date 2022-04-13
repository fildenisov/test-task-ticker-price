package rep

import "github.com/fildenisov/test-task-ticker-price/models"

// PriceStreamSubscriber is an original interface from test task
type PriceStreamSubscriber interface {
	SubscribePriceStream(models.Ticker) (chan models.TickerPrice, chan error)
}

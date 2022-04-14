package rep

import "github.com/fildenisov/test-task-ticker-price/models"

type Aggregator interface {
	Lifecycle
	PriceStreamSubscriber
	GetBars(t models.Ticker, max int) ([]models.Bar, bool)
}

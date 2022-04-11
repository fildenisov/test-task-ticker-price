package rep

import "price_aggregator/models"

type PriceStreamSubscriber interface {
	SubscribePriceStream(models.Ticker) (chan models.TickerPrice, chan error)
}

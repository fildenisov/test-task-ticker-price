# Test task: Ticker Price

## Makefile commands
### Run project
Use `make run` command to start project localy.

Default config is located in `configs/local/config.yaml`. 

Running with default config will start server at 127.0.0.1:8080

Every 60 seconds 150 mock stream events with be generated. 100 for `BTC_USD` and 50 for `ETH_USD` ticker. 

Wait some time and visit urls below to get aggregated indexes for tickers:

http://localhost:8080/v1/tickers/BTC_USD/bars/100 to get `BTC_USB` indexes.

http://localhost:8080/v1/tickers/ETH_USD/bars/100 to get `ETH_USD` indexes.

### ![image](https://raw.githubusercontent.com/fildenisov/test-task-ticker-price/29f092d2d4956b7d32546b0c12c33514aa9dfe46/coverage/coverage.svg)

Use `make cover` command to generate coverage report and coverage badge.

Click to view coverage report-> [coverage report](https://htmlpreview.github.io/?https://github.com/fildenisov/test-task-ticker-price/blob/master/coverage/index.html#file0)

### Testing
Use `make test` command to run tests.

`make race` command will run tests with race detection.

### Linting
Use `make check` to run linter.

## Project layout

Application consists of several component:
- internal/app - an instance of an application
- internal/rep - interface repository to separate application components
- delivery/http - a web server that allows user to get aggregated index information
- domain/aggregator - handles core logic for aggregating and storing index data 
- mocks/steam - is a mock for a test data stream
- models - stores models which are shared between application components
- consts - stores application constans

## Aggregator solution explanation

Here are 3 structs to store all data:
- bar
- bars
- Aggregator

#### bar
Represents single index bar (timestamp and aggregated price).

bar has only 1 method - update. It takes index price string (based on the task requirement), converts it fo float64 and returns error if price string is incorrect.

Aggregated index price is calculated as an average price = sum of all passed prices / total count of all passed prices.

#### bars
Stores a circular slice of `[]bar`, which means that if we reach the capacity limit (sets in config) next `bar` will be added to `bars[0]`. 

All stored bars are assosiated with exact `Ticker` (e.g. BTC_USD).

`bars` methods are goroutine safe which is proved by `make race`. 

`bars` stores a `bar` for **every** interval which means that if no value was passed in current interval it will be filled will an empty `bar` (count=0, val=0) value automaticaly. That guaranties that if value will be streamed with a delay `bars` will have a `bar` to update.

#### Aggregator
Stores mapping between `Ticker` and `bars`.

Goroutine safe (proved by `make race`).

Implements `PriceStreamSubscriber` interface given in task.

`SubscribePriceStream` is an entry point for index data. It allows `Aggregator` to subscribe on a new index stream which it gathers information from. 
Every `PriceStream` has it's own sepparated channel to communicate with `Aggregator`. 

`GetBars` methos can provide `Token` `bars` to any application components. In current solution `delivery/http` consumes this method to provide user an API output. 

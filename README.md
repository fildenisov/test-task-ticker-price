# Test task: Ticker Price

## Makefile commands
### Run project
Use `make run` command to start project localy.

Default config is located in `configs/local/config.yaml`. 

Running with default config will start server at 127.0.0.1:8080

Every 60 seconds 150 mock stream events with be generated. 100 for `BTC_USD` and 50 for `ETH_USD` ticker. 

Wait some time and visit:

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
- delivery/http - a web server that allowes user to get aggregated index information
- domain/aggregator - handles core logic for aggregating and storing index data 
- mocks/steam - is a mock for a test data stream
- models - stores models which are shared between application components
- consts - stores application constans

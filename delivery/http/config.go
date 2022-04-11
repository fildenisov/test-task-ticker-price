package http

import "time"

type Config struct {
	Address      string `default:":8080"`
	StartTimeout time.Duration
	StopTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

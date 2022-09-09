package main

import (
	"fmt"
	"time"

	"github.com/davidklassen/testing-workshop/pkg/greetings"
	"github.com/davidklassen/testing-workshop/pkg/retry"
)

type GreetingsRepository interface {
	FindAll() ([]*greetings.Greeting, error)
	FindRandom() (*greetings.Greeting, error)
}

type Greeter struct {
	greetingsRepo GreetingsRepository
	retryAttempts int
	retryDelay    time.Duration
	retryFactor   float64
}

func NewGreeter(greetingsRepo GreetingsRepository, retryAttempts int, retryDelay time.Duration, retryFactor float64) *Greeter {
	return &Greeter{
		greetingsRepo: greetingsRepo,
		retryAttempts: retryAttempts,
		retryDelay:    retryDelay,
		retryFactor:   retryFactor,
	}
}

func (greeter *Greeter) GetGreetings() ([]string, error) {
	var err error
	var gs []*greetings.Greeting
	err = retry.Retry(greeter.retryAttempts, greeter.retryDelay, greeter.retryFactor, func() error {
		gs, err = greeter.greetingsRepo.FindAll()
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find all greetings: %w", err)
	}

	res := make([]string, len(gs))
	for i, g := range gs {
		res[i] = g.Template()
	}

	return res, nil
}

func (greeter *Greeter) SayHello(name string) (string, error) {
	var err error
	var greeting *greetings.Greeting
	err = retry.Retry(greeter.retryAttempts, greeter.retryDelay, greeter.retryFactor, func() error {
		greeting, err = greeter.greetingsRepo.FindRandom()
		return err
	})
	if err != nil {
		return "", fmt.Errorf("failed to find random greeting: %w", err)
	}

	return greeting.Greet(name)
}

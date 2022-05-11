package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"golang.org/x/time/rate"
)

type APIConnection struct {
	rateLimiter *rate.Limiter
}

func (conn *APIConnection) ReadFile(ctx context.Context) error {
	if err := conn.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	fmt.Println("ReadFile")
	return nil
}

func (conn *APIConnection) ResolveAddress(ctx context.Context) error {
	if err := conn.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	fmt.Println("ResolveAddress")
	return nil
}

func Open() *APIConnection {
	return &APIConnection{
		rateLimiter: rate.NewLimiter(rate.Limit(100), 1),
	}
}

func main() {
	defer log.Printf("Done.")

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConn := Open()
	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			if err := apiConn.ReadFile(context.Background()); err != nil {
				log.Printf("failed to ReadFile: %w", err)
			}
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			if err := apiConn.ResolveAddress(context.Background()); err != nil {
				log.Printf("failed to ResolveAddress: %w", err)
			}
		}()
	}

	wg.Wait()
}

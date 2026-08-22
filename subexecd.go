package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	RedisServer  string
	RedisChannel string
	TriggerList  []Trigger
}

type Trigger struct {
	Message string
	Exec    []Exec
}

type Exec struct {
	Cmd     string
	Args    []string
	Timeout Duration
}

type Duration struct {
	time.Duration
}

func (duration *Duration) UnmarshalJSON(b []byte) error {
	var unmarshalledJSON interface{}

	err := json.Unmarshal(b, &unmarshalledJSON)
	if err != nil {
		return err
	}

	switch value := unmarshalledJSON.(type) {
	case float64:
		duration.Duration = time.Duration(value)
	case string:
		duration.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid duration: %#v", unmarshalledJSON)
	}

	return nil
}

func runExec(e Exec) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout.Duration)
	defer cancel()

	return exec.CommandContext(ctx, e.Cmd, e.Args...).Run()
}

func loadConfig(configfile string) {
	data, err := os.ReadFile(configfile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}
}

var cfg Config

func main() {
	configfile := flag.String("f", "subexecd.json", "path of the config file to use")
	flag.Parse()

	loadConfig(*configfile)

	client := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{cfg.RedisServer}})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}
	cancel()

	sub := client.Subscribe(context.Background(), cfg.RedisChannel)
	defer func() { _ = sub.Close() }()

	ch := sub.Channel()

	for msg := range ch {
		for _, trigger := range cfg.TriggerList {
			if trigger.Message == msg.Payload {
				for _, exec := range trigger.Exec {
					if err := runExec(exec); err != nil {
						log.Print(err)
					}
				}
				break
			}
		}
	}
}

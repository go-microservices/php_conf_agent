package main

import (
	"conf_agent/apollo"
	"conf_agent/config"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/howeyc/fsnotify"
)

var wg sync.WaitGroup
var ctx context.Context
var cancel context.CancelFunc

const LOOP = 1
const SYNC = 2

func main() {
	config.New()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel = context.WithCancel(context.Background())
	apolloConfig(ctx)

	go watchConfigFile(watcher)
	go handleSignals()

	err = watcher.Watch(config.AppConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	watcher.Close()
	log.Println("agent down")
}

func apolloConfig(ctx context.Context) {
	for _, configs := range config.Conf.Configs {
		for _, namespace := range configs.Namespace {
			wg.Add(1)
			go func(path, appId, namespace string, ctx context.Context) {
				switch config.Conf.Type {
				case LOOP:
					apollo.Loop(apollo.Configs{
						Path:      path,
						AppId:     appId,
						Namespace: namespace,
					}, &wg, ctx)
					break
				case SYNC:
					apollo.Sync(apollo.Configs{
						Path:          path,
						AppId:         appId,
						Namespace:     namespace,
						ReleaseKey:    "",
						Notifications: "",
					}, &wg, ctx)
					break

				default:
					log.Fatal("type ERR")

				}
				return

			}(configs.Path, configs.AppId, namespace, ctx)
		}

	}
}

func watchConfigFile(watcher *fsnotify.Watcher) {
	for {
		select {
		case ev := <-watcher.Event:
			log.Println("event:", ev)
			if ev.IsModify() {
				cancel()
				config.New()
				ctx, cancel = context.WithCancel(context.Background())
				apolloConfig(ctx)
			}

		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}
}

func handleSignals() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	for {
		switch <-signals {
		case syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT:
			log.Println("agent closing...")
			cancel()
		}
	}
}

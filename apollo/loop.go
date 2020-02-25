package apollo

import (
	"conf_agent/util"
	"context"
	"log"
	"sync"
	"time"
)

func Loop(configs Configs, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Printf("[appId] %v [Namespace] %v down...\n", configs.AppId, configs.Namespace)
			return
		default:
			log.Printf("[appId] %v [Namespace] %v update...\n",
				configs.AppId, configs.Namespace)

			config := ConfigCache(configs)
			util.Write(configs.Path+"/"+configs.Namespace, config)
			time.Sleep(30 * time.Second)
		}

	}
}

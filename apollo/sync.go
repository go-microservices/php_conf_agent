package apollo

import (
	"conf_agent/util"
	"context"
	"encoding/json"
	"log"
	"sync"
)

func Sync(configs Configs, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	var noId int64
	var releaseKey string
	for {
		select {
		case <-ctx.Done():
			log.Printf("[appId] %v [Namespace] %v [noId] %v [releaseKey] %s down...\n",
				configs.AppId, configs.Namespace, noId, releaseKey)
			return
		default:
			notifications, _ := json.Marshal([]map[string]interface{}{
				{
					"namespaceName":  configs.Namespace,
					"notificationId": noId,
				},
			})
			configs.Notifications = string(notifications)
			isUpdate, notificationId := Notifications(configs)
			if isUpdate {
				log.Printf("[appId] %v [Namespace] %v [noId] %v [releaseKey] %s update...\n",
					configs.AppId, configs.Namespace, noId, releaseKey)

				var config map[string]string
				configs.ReleaseKey = releaseKey
				releaseKey, config = ConfigFile(configs)
				util.Write(configs.Path+"/"+configs.Namespace, config)
				noId = notificationId
			}
		}

	}
}

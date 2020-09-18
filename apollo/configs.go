package apollo

import (
	"conf_agent/config"
	"conf_agent/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Configs struct {
	Path          string
	AppId         string
	Namespace     string
	ReleaseKey    string
	Notifications string
}

/*
通过带缓存的Http接口从Apollo读取配置
*/
func ConfigCache(configs Configs) (body map[string]string) {
	url := fmt.Sprintf("%s/configfiles/json/%s/%s/%s",
		config.Conf.Address, configs.AppId, config.Conf.ClusterName, configs.Namespace)
	if config.Conf.Ip != "" {
		url += "?ip=" + config.Conf.Ip
	} else if config.Conf.AutoIp == 1 {
		ip, err := util.ExternalIP()
		if err != nil {
			log.Println(err)
		} else {
			url += "?ip=" + ip.String()
		}
	}

	response, err := http.Get(url)
	if err != nil {
		log.Fatal("ConfigCache Get#" + err.Error())
	}
	if response.StatusCode == 200 {
		err = json.NewDecoder(response.Body).Decode(&body)
		if err != nil {
			log.Fatal("ConfigCache Decode#" + err.Error())
		}
		return
	}
	log.Fatalf("ConfigCache Get ERR %v\n", response.StatusCode)
	return
}

/**
应用感知配置更新
*/
func Notifications(configs Configs) (bool, int64) {
	query := map[string]string{
		"appId":         configs.AppId,
		"cluster":       config.Conf.ClusterName,
		"notifications": configs.Notifications,
	}
	url := fmt.Sprintf("%s/notifications/v2?%s",
		config.Conf.Address, util.FormQuery(query))
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Notifications Get#" + err.Error())
	}
	var body []struct {
		Namespace      string `json:"Namespace"`
		NotificationId int64  `json:"notificationId"`
		Messages       struct {
			Details map[string]int64 `json:"details"`
		} `json:"messages"`
	}

	if response.StatusCode == 200 {

		err = json.NewDecoder(response.Body).Decode(&body)
		if err != nil {
			log.Fatal("Notifications Decode#" + err.Error())
		}
		return true, body[0].NotificationId
	}
	return false, 0
}

/*
通过不带缓存的Http接口从Apollo读取配置
*/
func ConfigFile(configs Configs) (string, map[string]string) {
	url := fmt.Sprintf("%s/configs/%s/%s/%s?releaseKey=%s",
		config.Conf.Address, configs.AppId, config.Conf.ClusterName, configs.Namespace, configs.ReleaseKey)
	if config.Conf.Ip != "" {
		url += "&ip=" + config.Conf.Ip
	} else if config.Conf.AutoIp == 1 {
		ip, err := util.ExternalIP()
		if err != nil {
			log.Println(err)
		} else {
			url += "?ip=" + ip.String()
		}
	}
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("ConfigFile Get#" + err.Error())
	}
	var body struct {
		AppId          string            `json:"appId"`
		Cluster        string            `json:"cluster"`
		Namespace      string            `json:"Namespace"`
		Configurations map[string]string `json:"configurations"`
		ReleaseKey     string            `json:"releaseKey"`
	}
	if response.StatusCode == 200 {
		err = json.NewDecoder(response.Body).Decode(&body)
		if err != nil {
			log.Fatal("ConfigFile Decode#" + err.Error())
		}
		return body.ReleaseKey, body.Configurations
	}
	log.Fatalf("ConfigFile Get ERR# %v\n", response.StatusCode)
	return "", map[string]string{}
}

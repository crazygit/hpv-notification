package util

import (
	"encoding/json"
	"fmt"
	"github.com/crazygit/hpv-notification/config"
	"github.com/slack-go/slack"
)

func SendToSlackWebhook(webhookUrl string, msg *slack.WebhookMessage) error {
	// 仅用于调试时查看生成的消息json数据
	if config.AppConfig.Debug {
		b, err := json.MarshalIndent(msg, "", "    ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}
	return slack.PostWebhook(webhookUrl, msg)
}

package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/crazygit/hpv-notification/config"
	"github.com/crazygit/hpv-notification/internal/dal"
	"github.com/crazygit/hpv-notification/internal/dal/model"
	"github.com/crazygit/hpv-notification/internal/dal/query"
	"github.com/crazygit/hpv-notification/internal/util"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

const (
	BenDiBaoAPI = "https://wxapidg.bendibao.com/smartprogram/zhuanti.php"
	UserAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"
	HostHeader  = "wxapidg.bendibao.com"
)

// TODO: 根据自身需要修改需要监听的城市编码
var CityCodes = [2]string{"deyang", "cd"}

type (
	ResponseBody struct {
		Data struct {
			Website struct {
				Place []Place
			}
		}
	}

	Place struct {
		ID       string
		CityName string
		Name     string
		Addr     string
		// 疫苗名额
		MingE string
		// 预约条件
		Condition string
		// 预约方法
		Method   string
		Tel      string
		OrderId  string
		YYTime   string `json:"yy_time"`
		Course   string
		CityCode string
	}
)

// fetchCmd represents the sync command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch HPV info from api",
	Run: func(cmd *cobra.Command, args []string) {
		var places []Place
		for _, cityCode := range CityCodes {
			data, err := fetchData(cityCode)
			if err != nil {
				log.Fatalf("failed to fetch data, %s", err)
			}
			if data != nil {
				places = append(places, data...)
			}
		}
		ctx := cmd.Context()
		var notifications []Place
		for _, place := range places {
			exist, err := isPlaceExist(place, ctx)
			if err != nil {
				log.Fatalf("failed query place: %+v, err: %s", place, err)
			}
			if !exist {
				err := addPlace(place, ctx)
				if err != nil {
					log.Fatalf("failed add place: %+v, err: %s", place, err)
				}
				notifications = append(notifications, place)
			}
		}
		if notifications != nil {
			messages := genSlackWebhookMessage(notifications)
			for _, m := range messages {
				err := util.SendToSlackWebhook(config.AppConfig.Slack.WebhookURL, &m)
				if err != nil {
					log.Fatalf("failed send message: %+v, err: %s", m, err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func fetchData(cityCode string) ([]Place, error) {
	client := resty.New()
	var responseBody ResponseBody
	log.WithField("CityCode", cityCode).Info("fetch data")
	resp, err := client.R().
		SetHeader("User-Agent", UserAgent).
		SetHeader("Host", HostHeader).
		SetQueryParams(map[string]string{
			"platform": "wx",
			"version":  "21.12.06",
			"action":   "jiujia",
			"citycode": cityCode,
		}).
		ForceContentType("application/json").
		SetResult(&responseBody).
		Get(BenDiBaoAPI)

	if err != nil {
		return nil, fmt.Errorf("request API failed, error:%w", err)
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("bad response, response body is: %s", string(resp.Body()))
	}
	return responseBody.Data.Website.Place, nil
}

func isPlaceExist(p Place, ctx context.Context) (bool, error) {
	place := query.Use(dal.GetInstance()).Place
	v, err := place.WithContext(ctx).Where(place.ID.Eq(p.ID), place.CityCode.Eq(p.CityCode)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, err
	}
	return v != nil, nil
}

func addPlace(p Place, ctx context.Context) error {
	return query.Use(dal.GetInstance()).Place.WithContext(ctx).Create(&model.Place{
		ID:        p.ID,
		CityName:  p.CityName,
		Name:      p.Name,
		Addr:      p.Addr,
		MingE:     p.MingE,
		Condition: p.Condition,
		Method:    p.Method,
		Tel:       p.Tel,
		OrderId:   p.OrderId,
		YYTime:    p.YYTime,
		Course:    p.Course,
		CityCode:  p.CityCode,
	})
}

func genSlackWebhookMessage(places []Place) []slack.WebhookMessage {
	var messages []slack.WebhookMessage
	for _, p := range places {
		headerText := slack.NewTextBlockObject("plain_text", fmt.Sprintf("📣%s", p.Name), true, false)
		headerSection := slack.NewHeaderBlock(headerText)

		bodyText := slack.NewTextBlockObject("mrkdwn",
			fmt.Sprintf("*城市*: %s\n*地址*: %s\n*名额*: %s\n*电话*: %s\n*预约方法*: %s\n*预约条件*: %s\n*预约时间*: %s\n\n<%s|详情>\n",
				p.CityName, p.Addr, p.MingE, p.Tel, p.Method, p.Condition, p.YYTime, p.Course), false, false)
		bodySection := slack.NewSectionBlock(bodyText, nil, nil)
		messages = append(messages, slack.WebhookMessage{
			Blocks: &slack.Blocks{
				BlockSet: []slack.Block{
					headerSection,
					bodySection,
					slack.NewDividerBlock(),
				}},
		})
	}
	return messages
}

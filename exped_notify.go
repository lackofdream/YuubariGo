package yuubari_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Notifiable struct {
	*ProxyHandler
	//portDataCh to update port api data
	portDataCh chan PortAPI
	//cache to prevent duplicated notifications
	cache    *bigcache.BigCache
	tgBot    *tgbotapi.BotAPI
	tgUserId int64
}

func (n *Notifiable) notify(deckIdx int, deckName string, expedNo int, endTime int) {
	cacheKey := fmt.Sprintf("deck-%d-%d", deckIdx, endTime)
	_, err := n.cache.Get(cacheKey)
	if err == bigcache.ErrEntryNotFound {
		log.Infof("%d %s  ended exped #%d  at %d", deckIdx, deckName, expedNo, endTime)
		n.tgBot.Send(tgbotapi.NewMessage(
			n.tgUserId,
			fmt.Sprintf("Fleet %s(#%d) has finished expedition #%d", deckName, deckIdx+1, expedNo),
		))
		n.cache.Set(cacheKey, []byte{})
	}
}

func (n *Notifiable) notifyDaemon() {
	var state PortAPI
	tgUpdateConfig := tgbotapi.NewUpdate(0)
	tgUpdateConfig.Timeout = 60
	tgUpdatesCh, _ := n.tgBot.GetUpdatesChan(tgUpdateConfig)
	for {
		select {
		case <-time.After(time.Second):
			for idx, deck := range state.APIData.APIDeckPort {
				if idx == 0 || deck.APIMission[0] == 0 {
					continue
				}
				if time.Now().Unix()+59 > int64(deck.APIMission[2])/1000 {
					n.notify(idx, deck.APIName, deck.APIMission[1], deck.APIMission[2])
				}
			}
		case data := <-n.portDataCh:
			state = data
		case tgUpdate := <-tgUpdatesCh:
			log.Infof("[%s] %s", tgUpdate.Message.From.UserName, tgUpdate.Message.Text)
			switch tgUpdate.Message.Command() {
			case "list":
				sb := strings.Builder{}
				now := time.Now()
				for idx, deck := range state.APIData.APIDeckPort {
					if idx == 0 || deck.APIMission[0] == 0 {
						continue
					}
					sb.WriteString(fmt.Sprintf(
						"Fleet %s(#%d) is doing expedition #%d (%s left)\n",
						deck.APIName, idx+1, deck.APIMission[1],
						time.Unix(int64(deck.APIMission[2])/1000, 0).Sub(now)))
				}
				msg := tgbotapi.NewMessage(tgUpdate.Message.Chat.ID, sb.String())
				n.tgBot.Send(msg)
			}
		}

	}
}

func (n *Notifiable) updateTimers(req *http.Request, resp *http.Response) {
	if !strings.Contains(req.URL.Path, "/kcsapi/api_port/port") {
		return
	}
	respData := readResp(resp)
	if len(respData) < 7 {
		return
	}
	var data PortAPI
	json.NewDecoder(bytes.NewBuffer(respData[7:])).Decode(&data)
	n.portDataCh <- data
}

func MakeNotifiable(ph *ProxyHandler, tgBotToken string, tgUserID int64) *ProxyHandler {
	ret := Notifiable{
		ProxyHandler: ph,
		portDataCh:   make(chan PortAPI, 0),
		cache: func() *bigcache.BigCache {
			r, _ := bigcache.NewBigCache(bigcache.DefaultConfig(8 * time.Hour))
			return r
		}(),
		tgBot: func() *tgbotapi.BotAPI {
			r, err := tgbotapi.NewBotAPI(tgBotToken)
			if err != nil {
				log.Panic(err)
			}
			return r
		}(),
		tgUserId: tgUserID,
	}
	go ret.notifyDaemon()
	ret.RegisterPlugin(ret.updateTimers)
	return ret.ProxyHandler
}

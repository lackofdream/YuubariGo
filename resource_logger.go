package yuubari_go

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func logResources(req *http.Request, resp *http.Response) {
	if !strings.Contains(req.URL.Path, "/kcsapi/api_port/port") {
		return
	}
	respData := readResp(resp)
	if len(respData) < 7 {
		return
	}
	var data PortAPI
	json.NewDecoder(bytes.NewBuffer(respData[7:])).Decode(&data)
	if len(data.APIData.APIMaterial) < 8 {
		return
	}
	log.WithFields(log.Fields{
		"nickname": data.APIData.APIBasic.APINickname,
		"fuel": data.APIData.APIMaterial[0].APIValue,
		"ammo": data.APIData.APIMaterial[1].APIValue,
		"steel": data.APIData.APIMaterial[2].APIValue,
		"bauxite": data.APIData.APIMaterial[3].APIValue,
		"buildKit": data.APIData.APIMaterial[4].APIValue,
		"repairKit": data.APIData.APIMaterial[5].APIValue,
		"devKit": data.APIData.APIMaterial[6].APIValue,
		"revKit": data.APIData.APIMaterial[7].APIValue,
	}).Info("resources recorded")
}

func MakeResourceLogged(ph *ProxyHandler) *ProxyHandler {
	ph.RegisterPlugin(logResources)
	return ph
}

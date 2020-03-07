package Logic

import (
	"SLGGAME/AuthManager/ServiceGroup"
	"crypto/md5"
	"encoding/hex"
	jsoniter "github.com/json-iterator/go"
	"github.com/yaice-rx/yaice/utils"
	"net/http"
	"time"
)

type LoginResult struct {
	Guid      int64  `json:"guid"`
	SessionId int64  `json:"sessionId"`
	Port      int    `json:"port"`
	Host      string `json:"host"`
	Result    int    `json:"result"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var t int = 1
	serviceData := ServiceGroup.ServicesMgr.GetService(t)
	result := LoginResult{
		Guid:      time.Now().Unix(),
		SessionId: time.Now().Unix(),
		Port:      int(serviceData.Port),
		Host:      serviceData.Host,
		Result:    1,
	}
	data, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(result)
	w.Write(buildToken(data))
}

var token_ken string = "moidw908c9u8190802349"

func buildToken(data []byte) []byte {
	h := md5.New()
	h.Write([]byte(token_ken))
	md5 := h.Sum(nil)
	md5str := hex.EncodeToString(md5)
	md5len := utils.ShortToBytes(int16(len(md5str)))
	loginData := append(
		md5len,
		append([]byte(md5str),
			data...,
		)...,
	)
	return loginData
}

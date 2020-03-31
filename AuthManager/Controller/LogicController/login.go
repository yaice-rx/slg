package LogicController

import (
	"SLGGAME/AuthManager/Controller/GameController"
	"SLGGAME/Service"
	"crypto/md5"
	"encoding/hex"
	jsoniter "github.com/json-iterator/go"
	"github.com/yaice-rx/yaice/config"
	"github.com/yaice-rx/yaice/utils"
	"net/http"
)

var token_ken string = "moidw908c9u8190802349"

func Login(w http.ResponseWriter, r *http.Request) {
	result := Service.Token{
		Guid:      utils.GenSonyflake(),
		SessionId: config.ConfInstance().GetPid(),
		Port:      GameController.SnapPort,
		Host:      GameController.SnapHost,
		Result:    1,
	}
	data, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(result)
	w.Write(BuildToken(data))
}

func BuildToken(data []byte) []byte {
	h := md5.New()
	h.Write([]byte(token_ken))
	md5 := h.Sum(nil)
	md5str := hex.EncodeToString(md5)
	md5len := utils.ShortToBytes(int16(len(md5str)))
	loginData := append(md5len, append([]byte(md5str), data...)...)
	return loginData
}

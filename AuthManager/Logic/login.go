package Logic

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"github.com/yaice-rx/yaice/utils"
)

var token_ken string = "moidw908c9u8190802349"

func BuildToken(data []byte) []byte {
	h := md5.New()
	h.Write([]byte(token_ken))
	md5 := h.Sum(nil)
	md5str := hex.EncodeToString(md5)
	md5len := utils.ShortToBytes(int16(len(md5str)))
	logrus.Info("md5 :", []byte(md5str))
	loginData := append(md5len, append([]byte(md5str), data...)...)
	return loginData
}

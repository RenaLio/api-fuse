package load_balancer

import "strconv"

func GetKeyId(srv Server) string {
	key := srv.GetAPIKey()
	keyId := "不足8位，id->" + strconv.Itoa(int(srv.GetId()))
	if len(key) >= 8 {
		keyId = key[len(key)-8:]
	}
	return keyId
}

package pahttp

import "auiauto/perrorx"

func GetHello(avd string) (string, *perrorx.ErrorX) {
	ans, err := AntranceRequest("GET", avd, "hello", nil)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}

package pahttp

import "auiauto/perrorx"

func GetIsCrash(avd string) (string, *perrorx.ErrorX) {
	ans, err := AntranceRequest("GET", avd, "iscrash", nil)
	if err != nil {
		return "", perrorx.TransErrorX(err)
	}
	return ans, nil
}
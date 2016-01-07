package stormvada

import(
	sj "github.com/robsix/json"
)

func getSupportedFormats(accessToken string) (ret *sj.Json, err error) {
	req, err := newRequest("GET", "https://developer.api.autodesk.com/viewingservice/v1/supported", nil, accessToken, "")
	if err != nil {
		return nil, err
	}

	ret, err = doAdhocJsonRequest(req)
	return
}

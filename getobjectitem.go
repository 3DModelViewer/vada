package stormvada

import (
	"net/http"
)

func getObjectItem(b64UrnAndItemPath string, accessToken string) (ret *http.Response, err error) {
	req, err := newRequest("GET", "https://developer.api.autodesk.com/viewingservice/v1/items/"+b64UrnAndItemPath, nil, accessToken, "")
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	ret, err = client.Do(req)
	if ret != nil {
		err = checkResponse(ret, err)
	}
	return
}

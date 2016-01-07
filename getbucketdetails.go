package stormvada

import(
	sj "github.com/robsix/json"
)

func getBucketDetails(bucketKey string, accessToken string) (ret *sj.Json, err error) {
	req, err := newRequest("GET", "https://developer.api.autodesk.com/oss/v2/buckets/"+bucketKey+"/details", nil, accessToken, "")
	if err != nil {
		return nil, err
	}

	ret, err = doAdhocJsonRequest(req)
	return
}

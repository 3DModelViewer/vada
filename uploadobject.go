package stormvada

import (
	"git.autodesk.com/typhoon/stormfront/src/server/Godeps/_workspace/src/code.google.com/p/go-uuid/uuid"
	sj "git.autodesk.com/typhoon/stormfront/src/server/Godeps/_workspace/src/github.com/robsix/json"
	"io"
	"net/url"
)

func uploadObject(objectKey string, bucketKey string, objectReader io.Reader, accessToken string) (ret *sj.Json, err error) {
	url, err := url.Parse("https://developer.api.autodesk.com/oss/v2/buckets/"+bucketKey+"/objects/")
	if err != nil {
		return nil, err
	}

	url.Path += uuid.New() + objectKey

	req, err := newRequest("PUT", url.String(), objectReader, accessToken, "application/octet-stream")
	if err != nil {
		return nil, err
	}

	ret, err = doAdhocJsonRequest(req)
	return
}

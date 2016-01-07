package stormvada

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robsix/golog"
	"io"
	"net/http"
	"sync"
	"time"
	sj "github.com/robsix/json"
)

const (
	accessTokenExpirationBuffer = time.Duration(-10) * time.Second
)

type VadaClient interface {
	CreateBucket(bucketKey string, policyKey BucketPolicy) (*sj.Json, error)
	GetBucketDetails(bucketKey string) (*sj.Json, error)
	GetSupportedFormats() (*sj.Json, error)
	UploadObject(objectKey string, bucketKey string, objectReader io.Reader) (*sj.Json, error)
	RegisterObject(b64Urn string) (*sj.Json, error)
	GetObjectInfo(b64Urn string, guid string) (*sj.Json, error)
	GetObjectItem(b64UrnAndItemPath string) (*http.Response, error)
	GetObjectSeedFile(objectKey string, bucketKey string) (*http.Response, error)
}

func NewVadaClient(clientKey string, clientSecret string, log golog.Log) VadaClient {
	return &vadaClient{
		clientKey:    clientKey,
		clientSecret: clientSecret,
		log:          log,
	}
}

type vadaClient struct {
	clientKey          string
	clientSecret       string
	accessToken        string
	accessTokenExpires time.Time
	log                golog.Log
	mtx                sync.Mutex
}

func (v *vadaClient) getAccessToken() (string, error) {
	if time.Now().After(v.accessTokenExpires.Add(accessTokenExpirationBuffer)) {
		defer v.mtx.Unlock()
		v.mtx.Lock()
		if time.Now().After(v.accessTokenExpires.Add(accessTokenExpirationBuffer)) {
			v.log.Info("VadaClient.getAccessToken requesting new token")
			accessToken, err := getAccessToken(v.clientKey, v.clientSecret)
			if err != nil {
				v.log.Critical("VadaClient.getAccessToken error: ", err)
				return "", err
			}
			v.accessToken = accessToken.Token
			expiresDuration := time.Duration(accessToken.Expires) * time.Second
			v.accessTokenExpires = time.Now().Add(expiresDuration)
			v.log.Info("VadaClient.getAccessToken retrieved new access token: ", v.accessToken, " expires in: ", expiresDuration)
		}
	}
	return v.accessToken, nil
}

func (v *vadaClient) CreateBucket(bucketKey string, policyKey BucketPolicy) (*sj.Json, error) {
	token, err := v.getAccessToken()
	if err != nil {
		return nil, err
	}

	return createBucket(bucketKey, policyKey, token)
}

func (v *vadaClient) GetBucketDetails(bucketKey string) (*sj.Json, error) {
	token, err := v.getAccessToken()
	if err != nil {
		return nil, err
	}

	return getBucketDetails(bucketKey, token)
}

func (v *vadaClient) GetSupportedFormats() (*sj.Json, error) {
	token, err := v.getAccessToken()
	if err != nil {
		return nil, err
	}

	return getSupportedFormats(token)
}

func (v *vadaClient) UploadObject(objectKey string, bucketKey string, objectReader io.Reader) (*sj.Json, error) {
	token, err := v.getAccessToken()
	if err != nil {
		return nil, err
	}

	return uploadObject(objectKey, bucketKey, objectReader, token)
}

func (v *vadaClient) RegisterObject(b64Urn string) (*sj.Json, error) {
	token, err := v.getAccessToken()
	if err != nil {
		return nil, err
	}

	return registerObject(b64Urn, token)
}

func (v *vadaClient) GetObjectInfo(b64Urn string, guid string) (*sj.Json, error) {
	token, err := v.getAccessToken()
	if err != nil {
		return nil, err
	}

	return getObjectInfo(b64Urn, guid, token)
}

func (v *vadaClient) GetObjectItem(b64UrnAndItemPath string) (*http.Response, error) {
	token, err := v.getAccessToken()
	if err != nil {
		return nil, err
	}

	return getObjectItem(b64UrnAndItemPath, token)
}

func (v *vadaClient) GetObjectSeedFile(objectKey string, bucketKey string) (*http.Response, error) {
	token, err := v.getAccessToken()
	if err != nil {
		return nil, err
	}

	return getObjectSeedFile(objectKey, bucketKey, token)
}

/**
 * helpers
 */

func checkResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		body, _ := sj.FromReadCloser(resp.Body)
		bodyStr, _ := body.ToString()
		return errors.New(fmt.Sprintf("statusCode: %d, status: %v, body: %v", resp.StatusCode, resp.Status, bodyStr))
	}
	return nil
}

func newRequest(method string, urlStr string, body io.Reader, accessToken string, contentType string) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return req, nil
}

func doStructuredJsonRequest(req *http.Request, dst interface{}) error {
	client := http.DefaultClient
	resp, err := client.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err = checkResponse(resp, err); err != nil {
		return err
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(dst); err != nil {
		return err
	}

	return nil
}

func doAdhocJsonRequest(req *http.Request) (ret *sj.Json, err error) {
	client := http.DefaultClient
	resp, err := client.Do(req)
	if resp != nil {
		err = checkResponse(resp, err)
		if err == nil {
			ret, err = sj.FromReadCloser(resp.Body)
		}
	}
	return
}

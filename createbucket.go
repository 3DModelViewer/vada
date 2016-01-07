package stormvada

import (
	"errors"
	"regexp"
	sj "github.com/robsix/json"
)

type BucketPolicy string

const (
	Transient              = BucketPolicy("transient")
	Temporary              = BucketPolicy("temporary")
	Persistent             = BucketPolicy("persistent")
	bucketValidationRegexp = "[-_.a-z0-9]{3,128}"
)

func createBucket(bucketKey string, policyKey BucketPolicy, accessToken string) (ret *sj.Json, err error) {
	re := regexp.MustCompile(bucketValidationRegexp)
	if !re.MatchString(bucketKey) {
		return nil, errors.New("invalid bucket name: " + bucketKey + " must match regexp: " + bucketValidationRegexp)
	}

	data, err := sj.FromString(`{"bucketKey":"`+bucketKey+`","policyKey":"`+string(policyKey)+`"}`)
	if err != nil {
		return nil, err
	}

	reader, err := data.ToReader()
	if err != nil {
		return nil, err
	}

	req, err := newRequest("POST", "https://developer.api.autodesk.com/oss/v2/buckets", reader, accessToken, "application/json")
	if err != nil {
		return nil, err
	}

	ret, err = doAdhocJsonRequest(req)
	return
}

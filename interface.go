package vada

import(
	. "github.com/robsix/json"
	"io"
	"net/http"
)

type VadaClient interface {
	CreateBucket(bucketKey string, policyKey bucketPolicy) (*Json, error)
	GetBucketDetails(bucketKey string) (*Json, error)
	GetSupportedFormats() (*Json, error)
	UploadFile(objectKey string, bucketKey string, objectReader io.Reader) (*Json, error)
	RegisterFile(b64Urn string) (*Json, error)
	GetDocumentInfo(b64Urn string, guid string) (*Json, error)
	GetSheetItem(b64UrnAndItemPath string) (*http.Response, error)
	GetFile(objectKey string, bucketKey string) (*http.Response, error)
}
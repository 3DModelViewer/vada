package vada

const (
	Transient              = bucketPolicy("transient")
	Temporary              = bucketPolicy("temporary")
	Persistent             = bucketPolicy("persistent")
	bucketValidationRegexp = "[-_.a-z0-9]{3,128}"
)

type bucketPolicy string
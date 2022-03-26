package helpers

import "os"

var defaultAWSRegion = "us-east-1"

func GetAWSRegion() string {
	region := os.Getenv("AWS_REGION")
	if region != "" {
		return region
	}
	return defaultAWSRegion
}

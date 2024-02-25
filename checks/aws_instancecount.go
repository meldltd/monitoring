package checks

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
	"monitoring/spec"
	"strings"
)

func checkAWSInstanceCount(dsn string, checkSpec *spec.CheckSpec) (*map[string]string, error) {
	creds := strings.Split(dsn, "@")
	if 2 != len(creds) {
		return nil, fmt.Errorf("Invalid AWS credentials")
	}

	if nil == checkSpec.DSNParams {
		return nil, fmt.Errorf("DSNParams must be defined for AWS Check Instance Count")
	}

	regions := (*checkSpec.DSNParams)["regions"]
	if nil == regions {
		return nil, fmt.Errorf("DSNParams.regions must be defined for AWS Check Instance Count")
	}

	regionList := regions.([]interface{})
	if len(regionList) < 1 {
		return nil, fmt.Errorf("DSNParams.regions must be defined and have at least one element for AWS Check Instance Count")
	}

	result, err, m, done := gatherEC2results(regionList, creds)
	if done {
		return m, err
	}

	log.Println(result)
	return &result, err
}

func gatherEC2results(regionList []interface{}, creds []string) (map[string]string, error, *map[string]string, bool) {
	results := []string{""}
	total := 0
	filters := makeFilters()

	for _, region := range regionList {
		sdkConfig := aws.Config{
			Credentials: credentials.NewStaticCredentialsProvider(creds[0], creds[1], ""),
			Region:      region.(string),
		}

		client := ec2.NewFromConfig(sdkConfig)
		maxResults := int32(512) // Hardcoded

		out, err := client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{
			MaxResults: &maxResults,
			Filters:    filters,
		})
		if nil != err {
			return nil, errors.New(fmt.Sprintf("AWS EC2 DescribeInstances API call failed: %s", err.Error())), nil, true
		}

		results = append(results, fmt.Sprintf("%s: %d", region.(string), len(out.Reservations)))
		total += len(out.Reservations)
	}

	results = append(results, fmt.Sprintf("TOTAL COUNT: %d", total))
	msg := strings.Join(results, "\n")
	result := map[string]string{"msg": msg}
	return result, nil, nil, false
}

func makeFilters() []types.Filter {
	filterName := "instance-state-name"
	return []types.Filter{
		types.Filter{
			Name:   &filterName,
			Values: []string{"running", "pending"},
		},
	}
}

func (c *CheckHandler) CheckAWSInstanceCount(spec *spec.CheckSpec) (*map[string]string, error) {
	return checkAWSInstanceCount(spec.DSN, spec)
}

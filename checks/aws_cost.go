package checks

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"monitoring/spec"
	"strings"
	"time"
)

func checkAWSCost(checkSpec *spec.CheckSpec) (*map[string]string, error) {
	creds := strings.Split(checkSpec.DSN, "@")
	if 2 != len(creds) {
		return nil, fmt.Errorf("Invalid AWS credentials")
	}

	sdkConfig := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(creds[0], creds[1], ""),
	}
	client := costexplorer.NewFromConfig(sdkConfig)

	msg, err := formatMessage(client)
	if nil != err {
		return nil, err
	}

	result := map[string]string{"msg": msg}
	return &result, nil
}

func formatMessage(client *costexplorer.Client) (string, error) {
	now := time.Now()

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	startStr := firstOfMonth.Format("2006-01-02")
	endStr := lastOfMonth.Format("2006-01-02")
	out, err := client.GetCostAndUsage(context.Background(), &costexplorer.GetCostAndUsageInput{
		TimePeriod:  &types.DateInterval{Start: &startStr, End: &endStr},
		Granularity: "MONTHLY",
		Metrics:     []string{"BlendedCost"},
	})

	if nil != err {
		return "", errors.New("AWS CostExplorer API call failed")
	}

	msg := ""
	for _, v := range out.ResultsByTime {
		msg += fmt.Sprintf("Period: %s - %s, Amount: %s %s\n", *v.TimePeriod.Start, *v.TimePeriod.End, *v.Total["BlendedCost"].Amount, *v.Total["BlendedCost"].Unit)
	}

	startStr = time.Now().Format("2006-01-02")
	output, err := client.GetCostForecast(context.Background(), &costexplorer.GetCostForecastInput{
		TimePeriod:  &types.DateInterval{Start: &startStr, End: &endStr},
		Granularity: "MONTHLY",
		Metric:      "BLENDED_COST",
	})
	if err != nil {
		return "", errors.New("AWS CostExplorer API call failed")
	}

	for _, v := range output.ForecastResultsByTime {
		msg += fmt.Sprintf("Forecast: Period: %s - %s, Predicted amount: %s\n", *v.TimePeriod.Start, *v.TimePeriod.End, *v.MeanValue)
	}
	return msg, nil
}

func (c *CheckHandler) CheckAWSCosts(spec *spec.CheckSpec) (*map[string]string, error) {
	return checkAWSCost(spec)
}

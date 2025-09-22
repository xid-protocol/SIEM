package aws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/colin-404/logx"
	"github.com/spf13/viper"
)

type GdEvent struct {
	EventSource string             `json:"event_source"`
	InstanceId  string             `json:"instance_id"`
	Type        string             `json:"type"`
	AccountId   string             `json:"account_id"`
	Region      string             `json:"region"`
	Description string             `json:"description"`
	ID          string             `json:"id"`
	Severity    float64            `json:"severity"`
	RawData     *guardduty.Finding `json:"raw_data"`
}

// https://docs.aws.amazon.com/sdk-for-go/api/service/guardduty/
func (c *AWSCloud) GuardDuty(ctx context.Context) {
	//初始化 AWS 会话
	// sess, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String(viper.GetString("AWS.DefaultRegion")),
	// 	Credentials: credentials.NewStaticCredentials(viper.GetString("Guardduty.AwsApiKey"), viper.GetString("Guardduty.AwsSecretKey"), ""),
	// })
	// if err != nil {
	// 	panic(err)
	// }
	logx.Infof("开始获取GuardDuty事件")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(viper.GetString("AWS.DefaultRegion")),
		// 不设置 Credentials，让 AWS SDK 自动使用 IAM Role
	})
	if err != nil {
		panic(err)
	}

	// 创建 GuardDuty 客户端
	svc := guardduty.New(sess)

	updatedAt := time.Now().Add(-180 * time.Second)
	criteria := &guardduty.FindingCriteria{
		Criterion: map[string]*guardduty.Condition{
			"updatedAt": {
				GreaterThanOrEqual: aws.Int64(updatedAt.UnixNano() / int64(time.Millisecond)),
			},
		},
	}

	input := &guardduty.ListFindingsInput{
		DetectorId:      aws.String(viper.GetString("Guardduty.GuardDutyDetectorId")),
		FindingCriteria: criteria,
	}
	output, err := svc.ListFindings(input)
	if err != nil {
		log.Println("Error listing GuardDuty findings:", err)

	}

	for _, event := range output.FindingIds {
		findingInput := &guardduty.GetFindingsInput{
			DetectorId: aws.String(viper.GetString("Guardduty.GuardDutyDetectorId")),
			FindingIds: []*string{event},
		}
		findingOutput, err := svc.GetFindings(findingInput)
		if err != nil {
			log.Println(err)
		}
		// log.Println(findingOutput)

		findings := findingOutput.Findings
		for _, finding := range findings {
			fmt.Println(finding)
			c.GuardDutyChan <- GdEvent{
				EventSource: "aws_guarduty",
				InstanceId:  *finding.Resource.InstanceDetails.InstanceId,
				Region:      *finding.Region,
				Type:        *finding.Type,
				AccountId:   *finding.AccountId,
				Description: *finding.Description,
				ID:          *finding.Id,
				Severity:    *finding.Severity,
				RawData:     finding,
			}
		}
	}

}

package aws

import (
	"context"

	"github.com/spf13/viper"
	"github.com/xid-protocol/common"
)

// const (
// 	guarddutyPath   = "/protocols/SIEM/guarddutyEvent"
// 	mongoCollection = "aws_info"
// )

func (c *AWSCloud) Handler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case gdEvent := <-c.GuardDutyChan:
			process(ctx, gdEvent)
		}
	}
}

func process(ctx context.Context, gdEvent GdEvent) {

	//规范数据为XID
	// info := protocols.NewInfo(gdEvent.InstanceId, "aws-instanceID")
	// metadata := protocols.NewMetadata("create", guarddutyPath, "application/json")
	// XID := protocols.NewXID(&info, &metadata, gdEvent)

	//插入数据
	// db.MongoDB.Collection(mongoCollection).InsertOne(ctx, XID)

	//告警
	common.DoHttp("POST", viper.GetString("Syslog.Webhook"), gdEvent, nil)
}

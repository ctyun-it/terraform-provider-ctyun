package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2ConsumerStatusV3Api
/* 查询订阅组详细状态信息 */
type Mq2ConsumerStatusV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2ConsumerStatusV3Api(client *core.CtyunClient) *Mq2ConsumerStatusV3Api {
	return &Mq2ConsumerStatusV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/consumer/status",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2ConsumerStatusV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2ConsumerStatusV3Request) (*Mq2ConsumerStatusV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupName", req.GroupName)
	ctReq.AddParam("clientId", req.ClientId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2ConsumerStatusV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2ConsumerStatusV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
	ClientId   string `json:"clientId"`   /*  客户端ID  */
}

type Mq2ConsumerStatusV3Response struct {
	StatusCode *string                               `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                               `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2ConsumerStatusV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                               `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2ConsumerStatusV3ReturnObjResponse struct {
	Data *Mq2ConsumerStatusV3ReturnObjDataResponse `json:"data"` /*  消费组详情数据  */
}

type Mq2ConsumerStatusV3ReturnObjDataResponse struct {
	Properties      *Mq2ConsumerStatusV3ReturnObjDataPropertiesResponse        `json:"properties"`      /*  消费组配置属性  */
	SubscriptionSet []*Mq2ConsumerStatusV3ReturnObjDataSubscriptionSetResponse `json:"subscriptionSet"` /*  订阅关系集合  */
	MqTable         *Mq2ConsumerStatusV3ReturnObjDataMqTableResponse           `json:"mqTable"`         /*  消息队列信息表  */
	StatusTable     *Mq2ConsumerStatusV3ReturnObjDataStatusTableResponse       `json:"statusTable"`     /*  消费状态统计表  */
	Jstack          *string                                                    `json:"jstack"`          /*  线程堆栈信息  */
}

type Mq2ConsumerStatusV3ReturnObjDataPropertiesResponse struct {
	MaxReconsumeTimes             *string `json:"maxReconsumeTimes"`             /*  最大重试次数  */
	AdjustThreadPoolNumsThreshold *string `json:"adjustThreadPoolNumsThreshold"` /*  线程池数量调整阈值  */
	UnitMode                      *string `json:"unitMode"`                      /*  单元模式开关  */
	ConsumeTimeoutInSec           *string `json:"consumeTimeoutInSec"`           /*  消费超时时间（秒）  */
	TimeoutStrategy               *string `json:"timeoutStrategy"`               /*  超时处理策略  */
	ConsumerGroup                 *string `json:"consumerGroup"`                 /*  消费组名称  */
	MessageModel                  *string `json:"messageModel"`                  /*  消息模型  */
	AllocateMessageQueueStrategy  *string `json:"allocateMessageQueueStrategy"`  /*  消息队列分配策略  */
	PullThresholdSizeForTopic     *string `json:"pullThresholdSizeForTopic"`     /*  主题拉取大小阈值  */
	SuspendCurrentQueueTimeMillis *string `json:"suspendCurrentQueueTimeMillis"` /*  队列挂起时间（毫秒）  */
	PullThresholdSizeForQueue     *string `json:"pullThresholdSizeForQueue"`     /*  队列拉取大小阈值  */
	PROP_CLIENT_VERSION           *string `json:"PROP_CLIENT_VERSION"`           /*  客户端版本  */
	OffsetStore                   *string `json:"offsetStore"`                   /*  偏移量存储实现类  */
	ConsumeConcurrentlyMaxSpan    *string `json:"consumeConcurrentlyMaxSpan"`    /*  并发消费最大跨度  */
	PostSubscriptionWhenPull      *string `json:"postSubscriptionWhenPull"`      /*  拉取时提交订阅关系开关  */
	ConsumeTimestamp              *string `json:"consumeTimestamp"`              /*  消费时间戳  */
	PROP_CONSUME_TYPE             *string `json:"PROP_CONSUME_TYPE"`             /*  消费类型  */
	ConsumeMessageBatchMaxSize    *string `json:"consumeMessageBatchMaxSize"`    /*  批量消费最大消息数  */
	DefaultMQPushConsumerImpl     *string `json:"defaultMQPushConsumerImpl"`     /*  推送消费者实现类  */
	PROP_THREADPOOL_CORE_SIZE     *string `json:"PROP_THREADPOOL_CORE_SIZE"`     /*  线程池核心大小  */
	PullInterval                  *string `json:"pullInterval"`                  /*  拉取间隔（毫秒）  */
	PullThresholdForQueue         *string `json:"pullThresholdForQueue"`         /*  队列拉取阈值  */
	PullThresholdForTopic         *string `json:"pullThresholdForTopic"`         /*  主题拉取阈值  */
	ConsumeFromWhere              *string `json:"consumeFromWhere"`              /*  消费起始位置  */
	PROP_NAMESERVER_ADDR          *string `json:"PROP_NAMESERVER_ADDR"`          /*  名称服务地址  */
	PullBatchSize                 *string `json:"pullBatchSize"`                 /*  拉取批次大小  */
	ConsumeThreadMin              *string `json:"consumeThreadMin"`              /*  最小消费线程数  */
	PROP_CONSUMER_START_TIMESTAMP *string `json:"PROP_CONSUMER_START_TIMESTAMP"` /*  消费者启动时间戳  */
	ConsumeThreadMax              *string `json:"consumeThreadMax"`              /*  最大消费线程数  */
	ConsumeTimeout2               *string `json:"consumeTimeout2"`               /*  消费超时时间2（秒）  */
	Subscription                  *string `json:"subscription"`                  /*  订阅关系  */
	PROP_CONSUMEORDERLY           *string `json:"PROP_CONSUMEORDERLY"`           /*  顺序消费开关  */
	MessageListener               *string `json:"messageListener"`               /*  消息监听器实现类  */
}

type Mq2ConsumerStatusV3ReturnObjDataSubscriptionSetResponse struct {
	ClassFilterMode   *bool     `json:"classFilterMode"`   /*  类过滤模式开关  */
	Topic             *string   `json:"topic"`             /*  订阅主题  */
	SubString         *string   `json:"subString"`         /*  订阅表达式  */
	TagsSet           []*string `json:"tagsSet"`           /*  标签集合  */
	CodeSet           []*string `json:"codeSet"`           /*  代码集合  */
	SubVersion        *int64    `json:"subVersion"`        /*  订阅版本号  */
	ExpressionType    *string   `json:"expressionType"`    /*  表达式类型  */
	FilterClassSource *string   `json:"filterClassSource"` /*  过滤类源码  */
}

type Mq2ConsumerStatusV3ReturnObjDataMqTableResponse struct {
	CommitOffset            *int64 `json:"commitOffset"`            /*  提交偏移量  */
	CachedMsgMinOffset      *int64 `json:"cachedMsgMinOffset"`      /*  缓存消息最小偏移量  */
	CachedMsgMaxOffset      *int64 `json:"cachedMsgMaxOffset"`      /*  缓存消息最大偏移量  */
	CachedMsgCount          *int32 `json:"cachedMsgCount"`          /*  缓存消息数量  */
	CachedMsgSizeInMiB      *int32 `json:"cachedMsgSizeInMiB"`      /*  缓存消息大小（MiB）  */
	TransactionMsgMinOffset *int64 `json:"transactionMsgMinOffset"` /*  事务消息最小偏移量  */
	TransactionMsgMaxOffset *int64 `json:"transactionMsgMaxOffset"` /*  事务消息最大偏移量  */
	TransactionMsgCount     *int32 `json:"transactionMsgCount"`     /*  事务消息数量  */
	Locked                  *bool  `json:"locked"`                  /*  是否锁定  */
	TryUnlockTimes          *int32 `json:"tryUnlockTimes"`          /*  解锁尝试次数  */
	LastLockTimestamp       *int64 `json:"lastLockTimestamp"`       /*  最后锁定时间戳  */
	Droped                  *bool  `json:"droped"`                  /*  是否删除  */
	LastPullTimestamp       *int64 `json:"lastPullTimestamp"`       /*  最后拉取时间戳  */
	LastConsumeTimestamp    *int64 `json:"lastConsumeTimestamp"`    /*  最后消费时间戳  */
}

type Mq2ConsumerStatusV3ReturnObjDataStatusTableResponse struct {
	PullRT            *int32 `json:"pullRT"`            /*  拉取响应时间（毫秒）  */
	PullTPS           *int32 `json:"pullTPS"`           /*  拉取TPS  */
	ConsumeRT         *int32 `json:"consumeRT"`         /*  消费响应时间（毫秒）  */
	ConsumeOKTPS      *int32 `json:"consumeOKTPS"`      /*  消费成功TPS  */
	ConsumeFailedTPS  *int32 `json:"consumeFailedTPS"`  /*  消费失败TPS  */
	ConsumeFailedMsgs *int32 `json:"consumeFailedMsgs"` /*  消费失败消息数  */
}

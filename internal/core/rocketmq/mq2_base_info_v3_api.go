package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2BaseInfoV3Api
/* 查询实例基本信息 */
type Mq2BaseInfoV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2BaseInfoV3Api(client *core.CtyunClient) *Mq2BaseInfoV3Api {
	return &Mq2BaseInfoV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instance/baseInfo",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2BaseInfoV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2BaseInfoV3Request) (*Mq2BaseInfoV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2BaseInfoV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2BaseInfoV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例id  */
}

type Mq2BaseInfoV3Response struct {
	StatusCode *string                         `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    *string                         `json:"message"`    /*  接口调用状态描述。成功时为"success"，失败时为具体失败信息  */
	ReturnObj  *Mq2BaseInfoV3ReturnObjResponse `json:"returnObj"`  /*  核心返回对象  */
	Error      *string                         `json:"error"`      /*  错误码。仅失败时返回，描述具体错误信息  */
}

type Mq2BaseInfoV3ReturnObjResponse struct {
	InstanceBaseInfo *Mq2BaseInfoV3ReturnObjInstanceBaseInfoResponse `json:"instanceBaseInfo"` /*  MQ实例基础信息详情。仅接口调用成功时返回  */
}

type Mq2BaseInfoV3ReturnObjInstanceBaseInfoResponse struct {
	ProdInstId    *string `json:"prodInstId"`    /*  实例id  */
	BillMode      *string `json:"billMode"`      /*  计费模式 1：包周期；2：按需  */
	ProdInstName  *string `json:"prodInstName"`  /*  MQ实例名称  */
	RunningState  *string `json:"runningState"`  /*  实例运行状态  */
	ProdType      *string `json:"prodType"`      /*  产品类型  */
	MachineSpec   *string `json:"machineSpec"`   /*  机器规格  */
	TpsLimit      *string `json:"tpsLimit"`      /*  TPS限制值  */
	TopicLimit    *string `json:"topicLimit"`    /*  Topic数量限制。按需返回，可能为null  */
	DiskSpace     *string `json:"diskSpace"`     /*  磁盘空间大小  */
	NetName       *string `json:"netName"`       /*  网络名称（VPC名称）  */
	Subnet        *string `json:"subnet"`        /*  子网名称  */
	SecurityGroup *string `json:"securityGroup"` /*  安全组ID  */
	ExpTime       *string `json:"expTime"`       /*  过期时间。按需返回，可能为null  */
	NameServer    *string `json:"nameServer"`    /*  NameServer地址列表，多个地址以分号分隔  */
	CrtTime       *string `json:"crtTime"`       /*  创建时间（UTC格式）  */
}

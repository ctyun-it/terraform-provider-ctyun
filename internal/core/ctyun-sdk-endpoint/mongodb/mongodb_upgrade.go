package mongodb

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type MongodbUpgradeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbUpgradeApi(client *ctyunsdk.CtyunClient) *MongodbUpgradeApi {
	return &MongodbUpgradeApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/upgrade",
		},
	}
}

func (this *MongodbUpgradeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbUpgradeRequest, header *MongodbUpgradeRequestHeader) (upgradeResp *MongodbUpgradeResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if err != nil {
		return
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return
	}
	upgradeResp = &MongodbUpgradeResponse{}
	err = resp.Parse(upgradeResp)
	if err != nil {
		return
	}
	return upgradeResp, nil
}

type MongodbUpgradeRequest struct {
	InstId              string                 `json:"instId"`                        // 实例ID，必填
	NodeType            *string                `json:"nodeType,omitempty"`            // 节点类型，主节点或备份节点，非必填
	ProdId              *int64                 `json:"prodId,omitempty"`              // 产品ID，非必填
	DiskVolume          *int32                 `json:"diskVolume,omitempty"`          // 升级到的磁盘容量，单位G，非必填
	ProdPerformanceSpec *string                `json:"prodPerformanceSpec,omitempty"` // 产品规格，非必填
	IsUpgradeBackup     *bool                  `json:"isUpgradeBackup,omitempty"`     // DDS模块磁盘扩容时候会使用 是否主磁盘与备磁盘一起扩容
	AzList              []AvailabilityZoneInfo `json:"azList,omitempty"`              // 可用区节点相关信息，非必填
}

type MongodbUpgradeRequestHeader struct {
	ProjectID *string `json:"projectID"`
}

type AvailabilityZoneInfo struct {
	NodeType              *string `json:"nodeType"`              // 节点类型
	AvailabilityZoneName  string  `json:"availabilityZoneName"`  // 可用区名称
	AvailabilityZoneCount int32   `json:"availabilityZoneCount"` // 可用区数量
}

type MongodbUpgradeResponse struct {
	StatusCode int32                            `json:"statusCode"` // 接口状态码
	Error      string                           `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string                           `json:"message"`    // 描述信息
	ReturnObj  *MongodbUpgradeResponseReturnObj `json:"returnObj"`  // 返回对象，类型为 DataObject
}

type MongodbUpgradeResponseReturnObj struct {
	NewOrderId string `json:"newOrderId"` // 订单ID
}

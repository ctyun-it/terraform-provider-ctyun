package pgsql

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlUpgradeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlUpgradeApi(client *ctyunsdk.CtyunClient) *PgsqlUpgradeApi {
	return &PgsqlUpgradeApi{client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/upgrade",
		},
	}
}

type PgsqlUpgradeRequest struct {
	InstId              string                      `json:"instId"`                        // 实例id，必填
	NodeType            *string                     `json:"nodeType,omitempty"`            // 节点类型，optional, master=主节点, backup=备份节点
	ProdId              *int64                      `json:"prodId,omitempty"`              // 产品ID（升级目标的prodId），optional
	DiskVolume          *int32                      `json:"diskVolume,omitempty"`          // 升级到的磁盘容量，单位G，范围为[100，32768]，optional
	ProdPerformanceSpec *string                     `json:"prodPerformanceSpec,omitempty"` // 产品规格名，类似"4C8G"这种，optional
	AzList              []PgsqlUpgradeRequestAzList `json:"azList,omitempty"`              // 可用区节点相关信息，optional
}

type PgsqlUpgradeRequestHeader struct {
	ProjectID *string `json:"projectId,omitempty"` //项目id
}
type PgsqlUpgradeRequestAzList struct {
	AvailabilityZoneName  string `json:"availabilityZoneName"`  // 可用区名称
	AvailabilityZoneCount int32  `json:"availabilityZoneCount"` // 可用区数量
}

type PgsqlUpgradeResponse struct {
	StatusCode int32                          `json:"statusCode"` // 接口状态码
	Error      string                         `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string                         `json:"message"`    // 描述信息
	ReturnObj  *PgsqlUpgradeResponseReturnObj `json:"returnObj"`  // 返回对象，类型为 DataObject
}

type PgsqlUpgradeResponseReturnObj struct {
	NewOrderId string `json:"newOrderId"`
}

func (this *PgsqlUpgradeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlUpgradeRequest, header *PgsqlUpgradeRequestHeader) (upgrade *PgsqlUpgradeResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	response := PgsqlUpgradeResponse{}
	err = resp.Parse(&response)
	if err != nil {
		return
	}
	return &response, nil
}

package mysql

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// 关系数据库MySQL版实例开通
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=63&api=9672&data=75&isNormal=1&vid=69

type TeledbCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbCreateApi(client *ctyunsdk.CtyunClient) *TeledbCreateApi {
	return &TeledbCreateApi{client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/accept",
		},
	}
}

func (this *TeledbCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCreateRequest, header *TeledbCreateRequestHeader) (*TeledbCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(req)
	if err != nil {
		return nil, err
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return nil, err
	}
	response := TeledbCreateResponse{}
	err = resp.Parse(&response)
	//err = resp.ParseByStandardModelWithCheck(response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type TeledbCreateRequest struct {
	BillMode              string                     `json:"billMode"`
	RegionId              string                     `json:"regionId"`
	ProdVersion           string                     `json:"prodVersion"`
	ProdSpecName          string                     `json:"prodSpecName,omitempty"`
	AvailabilityZone      []string                   `json:"availabilityZone,omitempty"`
	VpcId                 string                     `json:"vpcId"`
	HostType              string                     `json:"hostType"`
	SubnetId              string                     `json:"subnetId"`
	SecurityGroupId       string                     `json:"securityGroupId"`
	Name                  string                     `json:"name"`
	Password              string                     `json:"password,omitempty"`
	Period                int32                      `json:"period"`
	Count                 int32                      `json:"count"`
	AutoRenewStatus       int32                      `json:"autoRenewStatus"`
	ProdId                int64                      `json:"prodId"`
	ProdPerformanceSpeces []string                   `json:"prodPerformanceSpeces,omitempty"`
	MysqlNodeInfoList     []MysqlNodeInfoListRequest `json:"mysqlNodeInfoList,omitempty"`
	CpuType               string                     `json:"cpuType"`
	OsType                string                     `json:"osType"`
}
type TeledbCreateRequestHeader struct {
	ProjectID *string `json:"project_id"`
}

type MysqlNodeInfoListRequest struct {
	NodeType             string                        `json:"nodeType"`             // master:实例规格(单机，一主一备，一主两备), readNode: 高级设置: 只读实例
	InstSpec             string                        `json:"instSpec"`             // 实例规格（默认：通用型=1）
	StorageType          string                        `json:"storageType"`          // 存储类型: SSD=超高IO、SATA=普通IO、SAS=高IO、SSD-genric=通用型SSD、FAST-SSD=极速型SSD
	StorageSpace         int32                         `json:"storageSpace"`         // 存储空间(单位:G，范围100,32768)
	ProdPerformanceSpec  string                        `json:"prodPerformanceSpec"`  // 规格(例: 4C8G)
	Disks                int32                         `json:"disks"`                // 磁盘（默认为1）
	BackupStorageType    string                        `json:"backupStorageType"`    // 存储类型: SSD=超高IO、SATA=普通IO、SAS=高IO、SSD-genric=通用型SSD、FAST-SSD=极速型SSD
	BackupStorageSpace   int32                         `json:"backupStorageSpace"`   // 备份空间大小（类型SATA）
	AvailabilityZoneInfo []AvailabilityZoneInfoRequest `json:"availabilityZoneInfo"` // 可用区信息
}
type AvailabilityZoneInfoRequest struct {
	AvailabilityZoneName  string `json:"availabilityZoneName"`  // 资源池可用区名称
	AvailabilityZoneCount int32  `json:"availabilityZoneCount"` // 资源池可用区总数
	NodeType              string `json:"nodeType"`              // 表示分布AZ的节点类型，master/slave/readNode
}

type returnObj struct {
	Data *returnObjData `json:"data,omitempty"`
}

type returnObjData struct {
	ErrorMessage *string `json:"errorMessage,omitempty"`
	Submitted    bool    `json:"submitted"`
	NewOrderId   *string `json:"newOrderId,omitempty"`
	NewOrderNo   *string `json:"newOrderNo,omitempty"`
	TotalPrice   float64 `json:"totalPrice"`
}
type TeledbCreateResponse struct {
	StatusCode int32      `json:"statusCode"`        // 接口状态码
	Error      *string    `json:"error,omitempty"`   // 错误码，失败时返回，成功时为空
	Message    *string    `json:"message,omitempty"` // 描述信息
	ReturnObj  *returnObj `json:"returnObj"`         // 返回对象
}

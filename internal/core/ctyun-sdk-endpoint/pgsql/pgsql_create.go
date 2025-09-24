package pgsql

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlCreateApi(client *ctyunsdk.CtyunClient) *PgsqlCreateApi {
	return &PgsqlCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-acceptor/v2/openapi/accept-order-info/accept",
		},
	}
}

func (this *PgsqlCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlCreateRequest, header *PgsqlCreateRequestHeader) (createResponse *PgsqlCreateResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectId != nil {
		builder.AddHeader("Project-Id", *header.ProjectId)
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	response := PgsqlCreateResponse{}
	err = resp.Parse(&response)
	if err != nil {
		return
	}
	return &response, nil
}

type PgsqlCreateRequest struct {
	BillMode              string                                `json:"billMode"`                        // 计费模式：包周期或按需
	RegionId              string                                `json:"regionId"`                        // 目标资源池ID
	ResPoolCode           *string                               `json:"resPoolCode,omitempty"`           // 资源池名称
	HostType              string                                `json:"hostType"`                        // 主机类型
	ProdVersion           string                                `json:"prodVersion"`                     // 产品版本
	MysqlNodeInfoList     []PgsqlCreateRequestMysqlNodeInfoList `json:"mysqlNodeInfoList"`               // Mysql节点信息
	ProdId                int64                                 `json:"prodId"`                          // 产品ID
	ProjectName           *string                               `json:"projectName,omitempty"`           // 企业项目名称，默认为default
	AutoScaleParam        *PgsqlCreateRequestAutoScaleParam     `json:"autoScaleParam,omitempty"`        // 自动扩容参数
	BackupStorageType     *string                               `json:"backupStorageType,omitempty"`     // 备份存储类型
	VpcId                 string                                `json:"vpcId,omitempty"`                 // 虚拟私有云ID
	SubnetId              string                                `json:"subnetId"`                        // 子网ID
	SecurityGroupId       string                                `json:"securityGroupId"`                 // 安全组ID
	AppointVip            *string                               `json:"appointVip,omitempty"`            // 指定VIP
	VpcName               *string                               `json:"vpcName,omitempty"`               // VPC名称
	SubnetName            *string                               `json:"subnetName,omitempty"`            // 子网名称
	SecurityGroupName     *string                               `json:"securityGroupName,omitempty"`     // 安全组名称
	Name                  string                                `json:"name"`                            // 集群名称
	Password              string                                `json:"password,omitempty"`              // 管理员密码（加密）
	ParamTemplateId       string                                `json:"paramTemplateId"`                 // 参数模板ID
	Period                int32                                 `json:"period"`                          // 购买时长（单位：月）
	Count                 int32                                 `json:"count"`                           // 购买数量
	AutoRenewStatus       int32                                 `json:"autoRenewStatus"`                 // 自动续订状态
	CaseSensitive         string                                `json:"caseSensitive"`                   // 是否区分大小写
	TimeZone              string                                `json:"timeZone"`                        // 时区
	OsType                *string                               `json:"osType,omitempty"`                // 操作系统类型
	CpuType               *string                               `json:"cpuType,omitempty"`               // CPU类型
	ProjectId             *string                               `json:"projectId,omitempty"`             // 企业项目ID
	ProdSpecName          *string                               `json:"prodSpecName,omitempty"`          // 产品规格名称
	IsMGR                 *string                               `json:"isMGR,omitempty"`                 // 是否开启MRG
	VpcCIDR               *string                               `json:"vpcCIDR,omitempty"`               // VPC网段
	ServerCollation       *string                               `json:"serverCollation,omitempty"`       // 数据库字符集
	CrossInstanceBackup   bool                                  `json:"crossInstanceBackup,omitempty"`   // 是否为恢复到新实例工单
	SourceInstId          *string                               `json:"sourceInstId,omitempty"`          // 源实例ID
	BackupId              *string                               `json:"backupId,omitempty"`              // 备份集ID
	BackupTimePoint       *string                               `json:"backupTimePoint,omitempty"`       // 恢复时间点
	IsCrossRegionRecovery bool                                  `json:"isCrossRegionRecovery,omitempty"` // 是否跨域恢复
}

type PgsqlCreateRequestAutoScaleParam struct {
	AutoScale       string `json:"autoScale"`       // 自动扩容标志，1为自动扩容，0为不自动扩容，必填
	MaxScale        int64  `json:"maxScale"`        // 存储扩容上限，单位为GB，必填
	ActiveScaleRate string `json:"activeScaleRate"` // 触发扩容的百分比，取值范围为1-100，必填
}

type PgsqlCreateRequestMysqlNodeInfoList struct {
	NodeType             string                                   `json:"nodeType"`                     // 节点类型，主实例或只读实例
	InstSpec             string                                   `json:"instSpec"`                     // 实例规格
	StorageType          string                                   `json:"storageType"`                  // 存储类型
	StorageSpace         int32                                    `json:"storageSpace"`                 // 存储空间，单位：GB
	ProdPerformanceSpec  string                                   `json:"prodPerformanceSpec"`          // 产品性能规格
	Disks                int32                                    `json:"disks"`                        // 磁盘数量，默认为1
	AvailabilityZoneInfo []PgsqlCreateRequestAvailabilityZoneInfo `json:"availabilityZoneInfo"`         // 可用区信息
	BackupStorageType    *string                                  `json:"backupStorageType,omitempty"`  // 备份存储类型
	BackupStorageSpace   *string                                  `json:"backupStorageSpace,omitempty"` // 备份存储空间大小
}
type PgsqlCreateRequestAvailabilityZoneInfo struct {
	AvailabilityZoneName  string  `json:"availabilityZoneName"`  // 可用区名称
	AvailabilityZoneCount int32   `json:"availabilityZoneCount"` // 资源池可用区总数
	NodeType              string  `json:"nodeType"`              // 分布AZ的节点类型，主/从/只读
	DisplayName           *string `json:"displayName,omitempty"` // 可用区显示名，非必传
	SpecId                string  `json:"specId"`                // 规格ID
}
type PgsqlCreateRequestHeader struct {
	ProjectId *string `json:"projectId,omitempty"`
}

type PgsqlCreateResponseReturnObjData struct {
	NewOrderId   *string `json:"newOrderId"` //订单id
	ErrorMessage *string `json:"errorMessage"`
	Submitted    bool    `json:"submitted"`
	NewOrderNo   string  `json:"newOrderNo"`
	TotalPrice   float32 `json:"totalPrice"`
}

type PgsqlCreateResponseReturnObj struct {
	Data PgsqlCreateResponseReturnObjData `json:"data"`
}

type PgsqlCreateResponse struct {
	StatusCode int32                         `json:"statusCode"`
	Error      string                        `json:"error"`     //错误码。当接口失败时才返回具体错误编码，成功不返回或者为空
	Message    string                        `json:"message"`   //描述信息
	ReturnObj  *PgsqlCreateResponseReturnObj `json:"returnObj"` // 返回对象
}

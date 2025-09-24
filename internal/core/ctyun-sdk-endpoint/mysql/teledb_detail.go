package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbQueryDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbQueryDetailApi(client *ctyunsdk.CtyunClient) *TeledbQueryDetailApi {
	return &TeledbQueryDetailApi{client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/instance",
		},
	}
}

func (this *TeledbQueryDetailApi) Do(ctx context.Context, credentials ctyunsdk.Credential, req *TeledbQueryDetailRequest, headers *TeledbQueryDetailRequestHeaders) (detailResp *TeledbQueryDetailResponse, err error) {
	builder := this.WithCredential(&credentials)
	_, err = builder.WriteJson(req)
	if headers.ProjectID != nil {
		builder.AddHeader("project-id", *headers.ProjectID)
	}
	if headers.InstID != "" {
		builder.AddHeader("inst-id", headers.InstID)
	}
	if headers.RegionID == "" {
		err = errors.New("regionId is empty")
		return
	}
	builder.AddHeader("regionId", headers.RegionID)
	if req.OuterProdInstId == "" {
		err = errors.New("outerProdInstId is empty")
		return
	}
	builder.AddParam("outerProdInstId", req.OuterProdInstId)
	if err != nil {
		return
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	detailResp = &TeledbQueryDetailResponse{}
	err = resp.Parse(detailResp)
	if err != nil {
		return
	}
	return detailResp, nil
}

type TeledbQueryDetailRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` //实例ID，必填
}
type TeledbQueryDetailRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"`
	InstID    string  `json:"instanceId"`
	RegionID  string  `json:"regionId"`
}

type TeledbQueryDetailResponse struct {
	StatusCode int32                `json:"statusCode"`
	Message    string               `json:"message"`
	ReturnObj  *DetailRespReturnObj `json:"returnObj"`
	Error      string               `json:"error"`
}

type DetailRespReturnObj struct {
	ProdBillType                int32    `json:"prodBillType"`
	ProdType                    int32    `json:"prodType"`
	ProdRunningStatus           int32    `json:"prodRunningStatus"`
	ProdOrderStatus             int32    `json:"prodOrderStatus"`
	Alive                       int32    `json:"alive"`
	Vip                         string   `json:"vip"`
	Vip6                        string   `json:"vip6"`
	WritePort                   string   `json:"writePort"`
	ReadPort                    string   `json:"readPort"`
	CreateTime                  int64    `json:"createTime"`
	ExpireTime                  int64    `json:"expireTime"`
	MachineSpec                 string   `json:"machineSpec"`
	ProdDbEngine                string   `json:"prodDbEngine"`
	DiskTotal                   string   `json:"diskTotal"`
	DiskRated                   int32    `json:"diskRated"`
	DiskType                    string   `json:"diskType"`
	DiskSize                    int32    `json:"diskSize"`
	BackupDiskUsedRated         int32    `json:"backupDiskUsedRated"`
	BackupUsedDiskSize          float32  `json:"backupUsedDiskSize"`
	BackupStoreType             int32    `json:"backupStoreType"`
	ProdInstSetName             string   `json:"prodInstSetName"`
	NetName                     string   `json:"netName"`
	EIP                         string   `json:"eIP"`
	EIPStatus                   int32    `json:"eipStatus"`
	Subnet                      string   `json:"subnet"`
	SSlStatus                   int32    `json:"sslStatus"`
	AzInfoList                  []AzInfo `json:"azInfoList"`
	NewMysqlVersion             string   `json:"newMysqlVersion"`
	SubnetId                    string   `json:"subnetId"`
	OrderId                     int64    `json:"orderId"`
	SecurityGroupId             string   `json:"securityGroupId"`
	TplName                     string   `json:"tplName"`
	AuditLogStatus              int32    `json:"auditLogStatus"`
	SecurityGroup               string   `json:"securityGroup"`
	OuterProdInstId             string   `json:"outerProdInstId"`
	ProdInstId                  int32    `json:"prodInstId"`
	ProdInstFlag                string   `json:"prodInstFlag"`
	DiskUsed                    string   `json:"diskUsed"`
	ParameterGroupUsed          string   `json:"parameterGroupUsed"`
	VpcId                       string   `json:"vpcId"`
	Usezos                      int32    `json:"usezos"`
	ProdInstName                string   `json:"prodInstName"`
	ProdId                      int64    `json:"prodId"`
	Resources                   string   `json:"resources"`
	UserId                      int32    `json:"userId"`
	ProdBillTime                int32    `json:"prodBillTime"`
	TenantId                    string   `json:"tenantId"`
	BackupDiskSize              int32    `json:"backupDiskSize"`
	TplCode                     string   `json:"tplCode"`
	DbMysqlVersion              string   `json:"dbMysqlVersion"`
	ProjectId                   string   `json:"projectId"`
	ProjectName                 string   `json:"projectName"`
	DiskDataUsed                string   `json:"diskDataUsed"`
	PauseEnable                 bool     `json:"pauseEnable"`
	LowerCaseTableNames         int32    `json:"lowerCaseTableNames"`
	InstReleaseProtectionStatus int32    `json:"instReleaseProtectionStatus"`
	MysqlPort                   string   `json:"mysqlPort"`
	DiskLogUsed                 string   `json:"diskLogUsed"`
	DiskDataUsedRated           int32    `json:"diskDataUsedRated"`
	Timezone                    string   `json:"timezone"`
	LastManualBackUp            int32    `json:"lastManualBackUp"`
	HostSeries                  string   `json:"hostSeries"`
	SecurityGroupStatus         int32    `json:"securityGroupStatus"`
	RenewalEnable               bool     `json:"renewalEnable"`
	IsMGR                       int32    `json:"isMGR"`
	DiskLogUsedRated            int32    `json:"diskLogUsedRated"`
}

type AzInfo struct {
	ResId  int64  `json:"resId"`
	Role   string `json:"role"`
	AzName string `json:"azName"`
	AzId   string `json:"azId"`
}

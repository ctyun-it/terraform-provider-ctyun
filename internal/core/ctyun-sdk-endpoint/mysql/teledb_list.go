package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetListApi(client *ctyunsdk.CtyunClient) *TeledbGetListApi {
	return &TeledbGetListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/instance/instance-list",
		},
	}
}

type TeledbGetListHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` //项目id
	RegionID  string  `json:"regionId"`            //资源区regionId，比如实例在资源区A，则需要填写A资源区的regionId
}

func (this *TeledbGetListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetListRequest, header *TeledbGetListHeaders) (*TeledbGetListResponse, error) {
	var err error
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if header.RegionID == "" {
		err = errors.New("regionId is empty")
		return nil, err
	}
	builder.AddHeader("regionId", header.RegionID)
	if err != nil {
		return nil, err
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return nil, err
	}
	var listResponse TeledbGetListResponse
	err = resp.Parse(&listResponse)
	if err != nil {
		return nil, err
	}
	return &listResponse, nil
}

type TeledbGetListRequest struct {
	PageNow      int32   `json:"pageNow"`                // 当前页
	PageSize     int32   `json:"pageSize"`               // 单页记录条数
	TagVOList    []TagVO `json:"tagVOList,omitempty"`    // 标签列表
	ResDbEngine  *string `json:"resDbEngine,omitempty"`  // 数据库引擎
	ProdInstName *string `json:"prodInstName,omitempty"` // 实例名称
	Vip          *string `json:"vip,omitempty"`          // 连接IP
}
type TagVO struct {
	Value   *string `json:"value,omitempty"`   // 标签值
	Key     *string `json:"key,omitempty"`     // 标签key
	LabelId *string `json:"labelId,omitempty"` // k-v的唯一标识
}

type TeledbGetListResponse struct {
	StatusCode int32                           `json:"statusCode"`        // 返回码
	Message    *string                         `json:"message,omitempty"` // 返回消息
	ReturnObj  *TeledbGetListResponseReturnObj `json:"returnObj"`         // 分页信息
	Error      *string                         `json:"error,omitempty"`   // 错误码（失败时才返回）
}

type TeledbGetListResponseReturnObj struct {
	PageNum          int32                                   `json:"pageNum"`          // 当前页
	PageSize         int32                                   `json:"pageSize"`         // 每页的数量
	Size             int32                                   `json:"size"`             // 当前页的数量
	StartRow         int32                                   `json:"startRow"`         // 当前页面第一个元素在数据库中的行号
	EndRow           int32                                   `json:"endRow"`           // 当前页面最后一个元素在数据库中的行号
	Total            int32                                   `json:"total"`            // 总记录数
	Pages            int32                                   `json:"pages"`            // 总页数
	FirstPage        int32                                   `json:"firstPage"`        // 第一页
	PrePage          int32                                   `json:"prePage"`          // 前一页
	IsFirstPage      bool                                    `json:"isFirstPage"`      // 是否为第一页
	IsLastPage       bool                                    `json:"isLastPage"`       // 是否为最后一页
	HasPreviousPage  bool                                    `json:"hasPreviousPage"`  // 是否有前一页
	HasNextPage      bool                                    `json:"hasNextPage"`      // 是否有下一页
	NavigatePages    int32                                   `json:"navigatePages"`    // 导航页码数
	NavigatePageNums []int32                                 `json:"navigatepageNums"` // 所有导航页号
	List             []TeledbGetListResponseReturnDetailList `json:"list"`             // 结果集(每页显示的数据)
}

type TeledbGetListResponseReturnDetailList struct {
	ProdInstName                string `json:"prodInstName"`
	OuterProdInstId             string `json:"outerProdInstId"`
	ProdBillType                int32  `json:"prodBillType"`
	ProdType                    int32  `json:"prodType"`
	ProdRunningStatus           int32  `json:"prodRunningStatus"`
	ProdOrderStatus             int32  `json:"prodOrderStatus"`
	Alive                       int32  `json:"alive"`
	Vip                         string `json:"vip"`
	Vip6                        string `json:"vip6"`
	WritePort                   string `json:"writePort"`
	ReadPort                    string `json:"readPort"`
	CreateTime                  int64  `json:"createTime"`
	ExpireTime                  int64  `json:"expireTime"`
	MachineSpec                 string `json:"machineSpec"`
	ProdDbEngine                string `json:"prodDbEngine"`
	DiskSize                    int32  `json:"diskSize"`
	NewMysqlVersion             string `json:"newMysqlVersion"`
	DiskType                    string `json:"diskType"`
	BackupDiskUsedRated         int32  `json:"backupDiskUsedRated"`
	NetName                     string `json:"netName"`
	Subnet                      string `json:"subnet"`
	ProdInstFlag                string `json:"prodInstFlag"`
	OrderId                     int64  `json:"orderId"`
	CanOperate                  int32  `json:"canOperate"`
	TplName                     string `json:"tplName"`
	AuditLogStatus              int32  `json:"auditLogStatus"`
	ParameterGroupUsed          string `json:"parameterGroupUsed"`
	VpcId                       string `json:"vpcId"`
	Usezos                      int32  `json:"usezos"`
	ProdInstId                  int64  `json:"prodInstId"`
	DbMysqlVersion              string `json:"dbMysqlVersion"`
	Resources                   string `json:"resources"`
	UserId                      int64  `json:"userId"`
	ProdBillTime                int32  `json:"prodBillTime"`
	TenantId                    string `json:"tenantId"`
	TplCode                     string `json:"tplCode"`
	ProdId                      int64  `json:"prodId"`
	ProdInstSetName             string `json:"prodInstSetName"`
	SecurityGroup               string `json:"securityGroup"`
	ProjectId                   string `json:"projectId"`
	ProjectName                 string `json:"projectName"`
	PauseEnable                 bool   `json:"pauseEnable"`
	InstReleaseProtectionStatus int32  `json:"instReleaseProtectionStatus"`
	LastManualBackUp            int64  `json:"lastManualBackUp"`
	RenewalEnable               bool   `json:"renewalEnable"`
	IsMGR                       int32  `json:"isMGR"`
	//ReadNode                    []Node `json:"readNode"`
}

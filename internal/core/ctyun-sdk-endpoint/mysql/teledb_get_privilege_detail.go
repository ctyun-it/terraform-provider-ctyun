package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetPrivilegeDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetPrivilegeDetailApi(client *ctyunsdk.CtyunClient) *TeledbGetPrivilegeDetailApi {
	return &TeledbGetPrivilegeDetailApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/account/account-detail",
		},
	}
}

func (this *TeledbGetPrivilegeDetailApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetPrivilegeDetailRequest, header *TeledbGetPrivilegeDetailRequestHeader) (CGetPrivilegeDetailResp *TeledbGetPrivilegeDetailResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if req.OuterProdInstId == "" || header.InstID == "" {
		err = errors.New("instId 为空")
		return
	}
	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	CGetPrivilegeDetailResp = &TeledbGetPrivilegeDetailResponse{}
	err = resp.Parse(CGetPrivilegeDetailResp)
	if err != nil {
		return
	}
	return CGetPrivilegeDetailResp, nil
}

type TeledbGetPrivilegeDetailRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 实例ID，必填
	AccountName     string `json:"accountName"`     // 数据库账户名， 必填
}
type TeledbGetPrivilegeDetailRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbGetPrivilegeDetailResponse struct {
	StatusCode int32                                       `json:"statusCode"`      // 接口状态码
	Error      *string                                     `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                      `json:"message"`         // 描述信息
	ReturnObj  []TeledbGetPrivilegeDetailResponseReturnObj `json:"returnObj"`
}

type TeledbGetPrivilegeDetailResponseReturnObj struct {
	AccountName   string `json:"accountName"`   // 数据库账户名
	GrantSchema   string `json:"grantSchema"`   // 授权数据库名
	SelectPriv    string `json:"selectPriv"`    // 用户是否可以通过SELECT命令选择数据 (Y/N)
	InsertPriv    string `json:"insertPriv"`    // 用户是否可以通过INSERT命令插入数据 (Y/N)
	UpdatePriv    string `json:"updatePriv"`    // 用户是否可以通过UPDATE命令修改现有数据 (Y/N)
	DeletePriv    string `json:"deletePriv"`    // 用户是否可以通过DELETE命令删除现有数据 (Y/N)
	CreatePriv    string `json:"createPriv"`    // 用户是否可以创建新的数据库和表 (Y/N)
	DropPriv      string `json:"dropPriv"`      // 用户是否可以删除现有数据库和表 (Y/N)
	ShowViewPriv  string `json:"show_viewPriv"` // 用户是否可以查看视图或了解视图如何执行 (Y/N)，此权限只在MySQL 5.0及更高版本中有意义
	AlterPriv     string `json:"alterPriv"`     // 用户是否可以修改或删除存储函数及函数 (Y/N)，此权限是在MySQL 5.0中引入的
	IndexPriv     string `json:"indexPriv"`     // 用户是否可以创建和删除表索引 (Y/N)
	LockTablePriv string `json:"lockTablePriv"` // 用户是否可以使用LOCK (Y/N)
}

package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbUpdateAccessWhiteList struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbUpdateAccessWhiteList(client *ctyunsdk.CtyunClient) *TeledbUpdateAccessWhiteList {
	return &TeledbUpdateAccessWhiteList{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/RDS2/v1/open-api/rds-manage/rds-access-white-list",
		},
	}
}

type TeledbUpdateAccessWhiteListRequest struct {
	OuterProdInstID string   `json:"outerProdInstId"`
	GroupName       string   `json:"groupName"`
	GroupWhiteList  []string `json:"groupWhiteList"`
}

type TeledbUpdateAccessWhiteListRequestHeader struct {
	ProjectID *string `json:"project-id"`
	InstID    *string `json:"inst-id"`
	RegionID  string  `json:"regionId"`
}

type TeledbUpdateAccessWhiteListResponse struct {
	StatusCode int32  `json:"statusCode"` // 接口状态码
	Error      string `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string `json:"message"`    // 描述信息
}

func (this *TeledbUpdateAccessWhiteList) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUpdateAccessWhiteListRequest, header *TeledbUpdateAccessWhiteListRequestHeader) (bindResponse *TeledbUpdateAccessWhiteListResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if header.InstID != nil {
		builder.AddHeader("inst-id", *header.InstID)
	}
	if header.RegionID == "" {
		err = errors.New("更新Mysql白名单，region_id必填")
	}
	if req.OuterProdInstID == "" {
		err = errors.New("更新Mysql白名单，prod_inst_id 必填")
		return
	}
	if req.GroupName == "" {
		err = errors.New("更新Mysql白名单，group_name 必填")
		return
	}
	if req.GroupWhiteList == nil {
		err = errors.New("更新Mysql白名单，group_white_list 必填")
		return
	}
	//builder.AddParam("outerProdInstId", req.OuterProdInstID)
	//builder.AddParam("groupName", req.GroupName)
	//builder.AddParam("groupWhiteList", req.GroupWhiteList)
	builder.AddHeader("regionId", header.RegionID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	bindResponse = &TeledbUpdateAccessWhiteListResponse{}
	err = resp.Parse(bindResponse)
	if err != nil {
		return
	}
	return bindResponse, nil
}

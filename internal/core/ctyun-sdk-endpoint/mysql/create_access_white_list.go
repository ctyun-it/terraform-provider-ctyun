package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbCreateAccessWhiteList struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbCreateAccessWhiteList(client *ctyunsdk.CtyunClient) *TeledbCreateAccessWhiteList {
	return &TeledbCreateAccessWhiteList{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/rds-manage/rds-access-white-list",
		},
	}
}

type TeledbCreateAccessWhiteListRequest struct {
	OuterProdInstID string   `json:"outerProdInstId"`
	GroupName       string   `json:"groupName"`
	GroupWhiteList  []string `json:"groupWhiteList"`
}

type TeledbCreateAccessWhiteListRequestHeader struct {
	ProjectID *string `json:"project-id"`
	InstID    *string `json:"inst-id"`
	RegionID  string  `json:"regionId"`
}

type TeledbCreateAccessWhiteListResponse struct {
	StatusCode int32  `json:"statusCode"` // 接口状态码
	Error      string `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string `json:"message"`    // 描述信息
}

func (this *TeledbCreateAccessWhiteList) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCreateAccessWhiteListRequest, header *TeledbCreateAccessWhiteListRequestHeader) (bindResponse *TeledbCreateAccessWhiteListResponse, err error) {
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
		err = errors.New("创建Mysql白名单，region_id必填")
		return
	}
	builder.AddHeader("regionId", header.RegionID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	bindResponse = &TeledbCreateAccessWhiteListResponse{}
	err = resp.Parse(bindResponse)
	if err != nil {
		return
	}
	return bindResponse, nil
}

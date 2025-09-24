package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// PgsqlGetNodeListApi 查询Pgsql实例节点列表
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=65&api=8917&data=72&isNormal=1&vid=67
type PgsqlGetNodeListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetNodeListApi(client *ctyunsdk.CtyunClient) *PgsqlGetNodeListApi {
	return &PgsqlGetNodeListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/node/list",
		},
	}
}

type PgsqlGetNodeListHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` //项目id
	RegionID  string  `json:"regionId"`            //资源区regionId，比如实例在资源区A，则需要填写A资源区的regionId
}

func (this *PgsqlGetNodeListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetNodeListRequest, header *PgsqlGetNodeListHeaders) (*PgsqlGetNodeListResponse, error) {
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
	if req.ProdInstId == "" {
		err = errors.New("实例id is empty")
		return nil, err
	}
	builder.AddParam("prodInstId", req.ProdInstId)

	if err != nil {
		return nil, err
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return nil, err
	}
	var listResponse PgsqlGetNodeListResponse
	err = resp.Parse(&listResponse)
	if err != nil {
		return nil, err
	}
	return &listResponse, nil
}

type PgsqlGetNodeListRequest struct {
	ProdInstId string `json:"prodInstId"` //实例id
}

type PgsqlGetNodeListResponse struct {
	StatusCode int32                               `json:"statusCode"`        // 返回码
	Message    *string                             `json:"message,omitempty"` // 返回消息
	ReturnObj  []PgsqlGetNodeListResponseReturnObj `json:"returnObj"`         // 分页信息
	Error      *string                             `json:"error,omitempty"`   // 错误码（失败时才返回）
}

type PgsqlGetNodeListResponseReturnObj struct {
	ProdInstId   string `json:"prodInstId"`
	NodeId       int64  `json:"nodeId"`
	Host         string `json:"host"`
	Port         int32  `json:"port"`
	Primary      int32  `json:"primary"`      // 是否为主节点，1:是，0:不是
	ExpireTime   string `json:"expireTime"`   //过期时间，格式为：yyyy:MM:dd HH:mm:ss
	ProdInstFlag string `json:"prodInstFlag"` //数据库类型
	AzId         string `json:"azId"`         //可用区id
	AzName       string `json:"azName"`       //可用区名称
}

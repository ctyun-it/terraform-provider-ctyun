package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutObjectRetentionApi
/* 设置对象的保留期限配置，在保留期限内的对象不可被彻底删除和篡改。
 */type ZosPutObjectRetentionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutObjectRetentionApi(client *core.CtyunClient) *ZosPutObjectRetentionApi {
	return &ZosPutObjectRetentionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-object-retention",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutObjectRetentionApi) Do(ctx context.Context, credential core.Credential, req *ZosPutObjectRetentionRequest) (*ZosPutObjectRetentionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosPutObjectRetentionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutObjectRetentionRequest struct {
	Bucket                    string `json:"bucket,omitempty"`          /*  桶名  */
	RegionID                  string `json:"regionID,omitempty"`        /*  区域 ID  */
	Key                       string `json:"key,omitempty"`             /*  对象名称  */
	VersionID                 string `json:"versionID,omitempty"`       /*  版本ID，在开启多版本时可使用  */
	BypassGovernanceRetention bool   `json:"bypassGovernanceRetention"` /*  指示此操作是否应绕过Governance模式限制  */
	RetentionMode             string `json:"retentionMode,omitempty"`   /*  保留模式，必须为 COMPLIANCE 或 GOVERNANCE  */
	RetainUntilDate           int64  `json:"retainUntilDate,omitempty"` /*  保留截止日期, utc 时间戳，单位秒，距当前时刻不超过 70 年（按1年365天计）  */
}

type ZosPutObjectRetentionResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

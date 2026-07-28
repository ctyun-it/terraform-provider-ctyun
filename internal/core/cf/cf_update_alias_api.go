package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfUpdateAliasApi
/* 修改别名 */
type CfUpdateAliasApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfUpdateAliasApi(client *core.CtyunClient) *CfUpdateAliasApi {
	return &CfUpdateAliasApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/openapi/v1/functions/{functionName}/aliases/{aliasName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfUpdateAliasApi) Do(ctx context.Context, credential core.Credential, req *CfUpdateAliasRequest) (*CfUpdateAliasResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder = builder.ReplaceUrl("aliasName", req.AliasName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfUpdateAliasRequest
		RegionId     interface{} `json:"regionId,omitempty"`
		FunctionName interface{} `json:"functionName,omitempty"`
		AliasName    interface{} `json:"aliasName,omitempty"`
	}{
		req, nil, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfUpdateAliasResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfUpdateAliasRequest struct {
	FunctionName string                    `json:"functionName"`          /*  函数名称，函数必须存在  */
	AliasName    string                    `json:"aliasName"`             /*  别名名称，别名必须存在  */
	RegionId     string                    `json:"regionId"`              /*  资源池id，标识不同的地区，如：华东1、西南1  */
	Gray         *CfUpdateAliasGrayRequest `json:"gray,omitempty"`        /*  别名灰度版本相关信息  */
	VersionId    string                    `json:"versionId"`             /*  主版本  */
	Description  *string                   `json:"description,omitempty"` /*  关于别名的描述  */
}

type CfUpdateAliasGrayRequest struct {
	RawType   int32                           `json:"type"`             /*  灰度类型，当前支持：1、按百分比随机灰度  */
	Config    *CfUpdateAliasGrayConfigRequest `json:"config,omitempty"` /*  该灰度类型配置  */
	VersionId string                          `json:"versionId"`        /*  灰度版本  */
}

type CfUpdateAliasGrayConfigRequest struct {
	Weight int32 `json:"weight"` /*  切流的比例。假设值为 5%，函数计算会将 5% 的流量到打到灰度版本，95% 的流量打到主版本。范围是[0-100]  */
}

type CfUpdateAliasResponse struct {
	StatusCode *int32                          `json:"statusCode"` /*  状态码。0表示成功，非0表示不成功  */
	Error      *string                         `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败。  */
	Message    *string                         `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfUpdateAliasReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfUpdateAliasReturnObjResponse struct {
	AliasName       *string                                        `json:"aliasName"`       /*  别名名称  */
	VersionId       *string                                        `json:"versionId"`       /*  主版本ID  */
	GrayVersionId   *string                                        `json:"grayVersionId"`   /*  灰度版本ID  */
	Description     *string                                        `json:"description"`     /*  关于别名的描述  */
	CreateTime      *string                                        `json:"createTime"`      /*  创建时间  */
	UpdateTime      *string                                        `json:"updateTime"`      /*  更新时间  */
	GrayType        *int32                                         `json:"grayType"`        /*  灰度类型，当前支持：1、按百分比随机灰度  */
	AliasGrayConfig *CfUpdateAliasReturnObjAliasGrayConfigResponse `json:"aliasGrayConfig"` /*  灰度配置  */
}

type CfUpdateAliasReturnObjAliasGrayConfigResponse struct {
	Weight *int32 `json:"weight"` /*  切流的比例。假设值为 5%，函数计算会将 5% 的流量到打到灰度版本，95% 的流量打到主版本。范围是[0-100]  */
}

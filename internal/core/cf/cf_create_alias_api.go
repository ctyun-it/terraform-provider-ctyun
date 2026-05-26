package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfCreateAliasApi
/* 创建别名 */
type CfCreateAliasApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfCreateAliasApi(client *core.CtyunClient) *CfCreateAliasApi {
	return &CfCreateAliasApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/functions/{functionName}/aliases",
			ContentType:  "application/json",
		},
	}
}

func (a *CfCreateAliasApi) Do(ctx context.Context, credential core.Credential, req *CfCreateAliasRequest) (*CfCreateAliasResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfCreateAliasRequest
		RegionId     interface{} `json:"regionId,omitempty"`
		FunctionName interface{} `json:"functionName,omitempty"`
	}{
		req, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfCreateAliasResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfCreateAliasRequest struct {
	FunctionName string                    `json:"functionName"`          /*  函数名称，函数必须存在  */
	RegionId     string                    `json:"regionId"`              /*  资源池id，标识不同的地区，如：华东1、西南1  */
	AliasName    string                    `json:"aliasName"`             /*  别名。只能包含字母、数字和中划线。只能字母开头，字母数字结尾。长度在 2~44之间。  */
	Description  *string                   `json:"description,omitempty"` /*  关于别名的描述  */
	VersionId    string                    `json:"versionId"`             /*  主版本ID  */
	Gray         *CfCreateAliasGrayRequest `json:"gray,omitempty"`        /*  灰度版本的配置  */
}

type CfCreateAliasGrayRequest struct {
	VersionId string                          `json:"versionId"`        /*  灰度版本ID  */
	RawType   int32                           `json:"type"`             /*  灰度类型，当前支持：1、按百分比随机灰度  */
	Config    *CfCreateAliasGrayConfigRequest `json:"config,omitempty"` /*  对应类型的配置  */
}

type CfCreateAliasGrayConfigRequest struct {
	Weight int32 `json:"weight"` /*  切流的比例。假设值为 5%，函数计算会将 5% 的流量到打到灰度版本，95% 的流量打到主版本。范围是[0-100]  */
}

type CfCreateAliasResponse struct {
	StatusCode *int32                          `json:"statusCode"` /*  状态码。0表示成功，非0表示不成功  */
	Error      *string                         `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败。  */
	Message    *string                         `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfCreateAliasReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfCreateAliasReturnObjResponse struct {
	AliasName       *string                                        `json:"aliasName"`       /*  别名  */
	VersionId       *string                                        `json:"versionId"`       /*  主版本ID  */
	GrayVersionId   *string                                        `json:"grayVersionId"`   /*  灰度版本ID  */
	Description     *string                                        `json:"description"`     /*  关于别名的描述  */
	CreateTime      *string                                        `json:"createTime"`      /*  创建时间  */
	UpdateTime      *string                                        `json:"updateTime"`      /*  更新时间  */
	GrayType        *int32                                         `json:"grayType"`        /*  灰度类型，当前支持：1、按百分比随机灰度  */
	AliasGrayConfig *CfCreateAliasReturnObjAliasGrayConfigResponse `json:"aliasGrayConfig"` /*  灰度配置  */
}

type CfCreateAliasReturnObjAliasGrayConfigResponse struct {
	Weight *int32 `json:"weight"` /*  切流的比例。假设值为 5%，函数计算会将 5% 的流量到打到灰度版本，95% 的流量打到主版本。范围是[0-100]  */
}

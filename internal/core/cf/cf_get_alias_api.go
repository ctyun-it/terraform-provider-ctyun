package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetAliasApi
/* 查询别名详情 */
type CfGetAliasApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetAliasApi(client *core.CtyunClient) *CfGetAliasApi {
	return &CfGetAliasApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/functions/{functionName}/aliases/{aliasName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetAliasApi) Do(ctx context.Context, credential core.Credential, req *CfGetAliasRequest) (*CfGetAliasResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder = builder.ReplaceUrl("aliasName", req.AliasName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetAliasResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetAliasRequest struct {
	FunctionName string `json:"functionName"` /*  函数名称，函数必须存在。  */
	AliasName    string `json:"aliasName"`    /*  别名名称，别名必须存在。  */
	RegionId     string `json:"regionId"`     /*  资源池id，标识不同的地区，如：华东1、西南1  */
}

type CfGetAliasResponse struct {
	StatusCode *int32                       `json:"statusCode"` /*  状态码。0表示成功，非0表示不成功。  */
	Code       *string                      `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败。  */
	Message    *string                      `json:"message"`    /*  错误描述信息信息  */
	ReturnObj  *CfGetAliasReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfGetAliasReturnObjResponse struct {
	AliasName       *string                                     `json:"aliasName"`       /*  别名  */
	VersionId       *string                                     `json:"versionId"`       /*  主版本ID  */
	GrayVersionId   *string                                     `json:"grayVersionId"`   /*  灰度版本ID  */
	Description     *string                                     `json:"description"`     /*  关于别名的描述  */
	CreateTime      *string                                     `json:"createTime"`      /*  创建时间  */
	UpdateTime      *string                                     `json:"updateTime"`      /*  更新时间  */
	GrayType        *int32                                      `json:"grayType"`        /*  灰度类型，当前支持：1、按百分比随机灰度  */
	AliasGrayConfig *CfGetAliasReturnObjAliasGrayConfigResponse `json:"aliasGrayConfig"` /*  灰度配置  */
}

type CfGetAliasReturnObjAliasGrayConfigResponse struct {
	Weight *int32 `json:"weight"` /*  切流的比例。假设值为 5%，函数计算会将 5% 的流量到打到灰度版本，95% 的流量打到主版本。范围是[0-100]  */
}

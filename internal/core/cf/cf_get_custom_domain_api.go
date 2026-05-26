package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfGetCustomDomainApi
/* 获取自定义域名配置 */
type CfGetCustomDomainApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetCustomDomainApi(client *core.CtyunClient) *CfGetCustomDomainApi {
	return &CfGetCustomDomainApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/domains/customdomains/{domainName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetCustomDomainApi) Do(ctx context.Context, credential core.Credential, req *CfGetCustomDomainRequest) (*CfGetCustomDomainResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("domainName", req.DomainName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("cnameCheck", strconv.FormatBool(req.CnameCheck))
	if req.Region != nil && *req.Region != "" {
		ctReq.AddParam("region", *req.Region)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetCustomDomainResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetCustomDomainRequest struct {
	DomainName string  `json:"domainName"`       /*  域名  */
	RegionId   string  `json:"regionId"`         /*  资源池id  */
	CnameCheck bool    `json:"cnameCheck"`       /*  可选，是否检查cname配置，true：需要检查，false：不需要检查。检查结果在响应体的cnameValid字段（仅需要检查时该响应字段才有意义）  */
	Region     *string `json:"region,omitempty"` /*  指定自定义域名所在资源池，默认为当前资源池  */
}

type CfGetCustomDomainResponse struct {
	StatusCode *int32                              `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string                             `json:"error"`      /*  错误码  */
	Message    *string                             `json:"message"`    /*  信息  */
	ReturnObj  *CfGetCustomDomainReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfGetCustomDomainReturnObjResponse struct {
	CreatedAt    *string                                         `json:"createdAt"`    /*  创建时间  */
	UpdatedAt    *string                                         `json:"updatedAt"`    /*  更新时间  */
	Creator      *int32                                          `json:"creator"`      /*  创建者ID  */
	Editor       *int32                                          `json:"editor"`       /*  编辑者ID  */
	CertConfig   *CfGetCustomDomainReturnObjCertConfigResponse   `json:"certConfig"`   /*  HTTPS证书信息  */
	DomainName   *string                                         `json:"domainName"`   /*  自定义域名  */
	Protocol     *string                                         `json:"protocol"`     /*  协议类型  */
	RouteConfig  *CfGetCustomDomainReturnObjRouteConfigResponse  `json:"routeConfig"`  /*  路由表  */
	AuthConfig   *CfGetCustomDomainReturnObjAuthConfigResponse   `json:"authConfig"`   /*  认证配置  */
	DomainStatus *string                                         `json:"domainStatus"` /*  域名备案状态  */
	FilingStatus *CfGetCustomDomainReturnObjFilingStatusResponse `json:"filingStatus"` /*  域名备案详细信息  */
	Description  *string                                         `json:"description"`  /*  描述  */
	CnameValid   *bool                                           `json:"cnameValid"`   /*  CNAME是否有效，仅在应用场景下有意义  */
}

type CfGetCustomDomainReturnObjCertConfigResponse struct {
	Certificate *string `json:"certificate"` /*  证书  */
	CertName    *string `json:"certName"`    /*  证书名称  */
	PrivateKey  *string `json:"privateKey"`  /*  私钥  */
}

type CfGetCustomDomainReturnObjRouteConfigResponse struct {
	Routes []*CfGetCustomDomainReturnObjRouteConfigRoutesResponse `json:"routes"` /*  路由映射  */
}

type CfGetCustomDomainReturnObjAuthConfigResponse struct {
	AuthType  *string                                                `json:"authType"`  /*  认证类型  */
	JwtConfig *CfGetCustomDomainReturnObjAuthConfigJwtConfigResponse `json:"jwtConfig"` /*  JWT配置  */
}

type CfGetCustomDomainReturnObjFilingStatusResponse struct {
	Domain       *string `json:"domain"`       /*  一级域名  */
	RecordStatus *string `json:"recordStatus"` /*  备案状态,0：未备案 1：已备案  */
	RecordNumber *string `json:"recordNumber"` /*  备案号  */
	ErrMessage   *string `json:"errMessage"`   /*  异常信息  */
}

type CfGetCustomDomainReturnObjRouteConfigRoutesResponse struct {
	EnableJwt          *int32    `json:"enableJwt"`          /*  是否启用JWT  */
	FunctionId         *int32    `json:"functionId"`         /*  函数ID  */
	FunctionName       *string   `json:"functionName"`       /*  函数名称  */
	FunctionUniqueName *string   `json:"functionUniqueName"` /*  函数唯一名称  */
	Methods            []*string `json:"methods"`            /*  请求方法  */
	Path               *string   `json:"path"`               /*  请求路径  */
	Qualifier          *string   `json:"qualifier"`          /*  函数版本  */
}

type CfGetCustomDomainReturnObjAuthConfigJwtConfigResponse struct {
	Jwks        *string                                                             `json:"jwks"`        /*  JWK字符串  */
	TokenConfig []*CfGetCustomDomainReturnObjAuthConfigJwtConfigTokenConfigResponse `json:"tokenConfig"` /*  Token配置  */
	MatchMode   *CfGetCustomDomainReturnObjAuthConfigJwtConfigMatchModeResponse     `json:"matchMode"`   /*  JWT匹配模式，固定为All  */
}

type CfGetCustomDomainReturnObjAuthConfigJwtConfigTokenConfigResponse struct {
	Location *string `json:"location"` /*  Token读取位置  */
	Name     *string `json:"name"`     /*  字段名  */
}

type CfGetCustomDomainReturnObjAuthConfigJwtConfigMatchModeResponse struct {
	Mode *string   `json:"mode"` /*  匹配模式  */
	Path []*string `json:"path"` /*  匹配路径  */
}

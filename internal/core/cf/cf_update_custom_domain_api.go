package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfUpdateCustomDomainApi
/* 更新自定义域名 */
type CfUpdateCustomDomainApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfUpdateCustomDomainApi(client *core.CtyunClient) *CfUpdateCustomDomainApi {
	return &CfUpdateCustomDomainApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/openapi/v1/domains/customdomains/{domainName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfUpdateCustomDomainApi) Do(ctx context.Context, credential core.Credential, req *CfUpdateCustomDomainRequest) (*CfUpdateCustomDomainResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("domainName", req.DomainName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfUpdateCustomDomainRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfUpdateCustomDomainResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfUpdateCustomDomainRequest struct {
	DomainName  string               `json:"domainName"`            /*  域名 */
	RegionId    string               `json:"regionId"`              /*  资源池 id */
	Protocol    *string              `json:"protocol,omitempty"`    /*  协议类型：HTTP, HTTPS, HTTP,HTTPS */
	Description *string              `json:"description,omitempty"` /*  自定义描述 */
	CertConfig  *CfUpdateCertConfig  `json:"certConfig,omitempty"`  /*  HTTPS 证书的信息 */
	RouteConfig *CfUpdateRouteConfig `json:"routeConfig,omitempty"` /*  自定义域名访问时的 PATH 到 Function 映射 */
	AuthConfig  *CfUpdateAuthConfig  `json:"authConfig,omitempty"`  /*  认证配置信息 */
}

type CfUpdateCustomDomainResponse struct {
	StatusCode *int32                                 `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string                                `json:"error"`      /*  错误码  */
	Message    *string                                `json:"message"`    /*  信息  */
	ReturnObj  *CfUpdateCustomDomainReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfUpdateCustomDomainReturnObjResponse struct {
	CreatedAt    *string                                            `json:"createdAt"`    /*  创建时间 */
	UpdatedAt    *string                                            `json:"updatedAt"`    /*  更新时间 */
	Creator      *int32                                             `json:"creator"`      /*  创建者 ID */
	Editor       *int32                                             `json:"editor"`       /*  编辑者 ID */
	CertConfig   *CfUpdateCustomDomainReturnObjCertConfigResponse   `json:"certConfig"`   /*  HTTPS 证书信息 */
	DomainName   *string                                            `json:"domainName"`   /*  自定义域名 */
	Protocol     *string                                            `json:"protocol"`     /*  协议类型 */
	RouteConfig  *CfUpdateCustomDomainReturnObjRouteConfigResponse  `json:"routeConfig"`  /*  路由表 */
	AuthConfig   *CfUpdateCustomDomainReturnObjAuthConfigResponse   `json:"authConfig"`   /*  认证配置 */
	DomainStatus *string                                            `json:"domainStatus"` /*  域名备案状态 */
	FilingStatus *CfUpdateCustomDomainReturnObjFilingStatusResponse `json:"filingStatus"` /*  域名备案详细信息 */
	Description  *string                                            `json:"description"`  /*  描述 */
	CnameValid   *bool                                              `json:"cnameValid"`   /*  CNAME 是否有效，仅在应用场景下有意义 */
}

type CfUpdateCustomDomainReturnObjCertConfigResponse struct {
	Certificate *string `json:"certificate"` /*  证书  */
	CertName    *string `json:"certName"`    /*  证书名称  */
	PrivateKey  *string `json:"privateKey"`  /*  私钥  */
}

type CfUpdateCustomDomainReturnObjRouteConfigResponse struct {
	Routes []*CfUpdateCustomDomainReturnObjRouteConfigRoutesResponse `json:"routes"` /*  路由映射  */
}

type CfUpdateCustomDomainReturnObjAuthConfigResponse struct {
	AuthType  *string                                         `json:"authType"`  /*  认证类型 */
	JwtConfig *CfUpdateCustomDomainReturnObjJwtConfigResponse `json:"jwtConfig"` /*  JWT 配置 */
}

type CfUpdateCustomDomainReturnObjFilingStatusResponse struct {
	Domain       *string `json:"domain"`       /*  一级域名  */
	RecordStatus *string `json:"recordStatus"` /*  备案状态,0：未备案 1：已备案  */
	RecordNumber *string `json:"recordNumber"` /*  备案号  */
	ErrMessage   *string `json:"errMessage"`   /*  异常信息  */
}

type CfUpdateCustomDomainReturnObjRouteConfigRoutesResponse struct {
	EnableJwt          *int32    `json:"enableJwt"`          /*  是否启用JWT  */
	FunctionId         *int32    `json:"functionId"`         /*  函数ID  */
	FunctionName       *string   `json:"functionName"`       /*  函数名称  */
	FunctionUniqueName *string   `json:"functionUniqueName"` /*  函数唯一名称  */
	Methods            []*string `json:"methods"`            /*  请求方法  */
	Path               *string   `json:"path"`               /*  请求路径  */
	Qualifier          *string   `json:"qualifier"`          /*  函数版本  */
}

type CfUpdateCustomDomainReturnObjJwtConfigResponse struct {
	Jwks        *string                                                      `json:"jwks"`        /*  JWK 字符串 */
	TokenConfig []*CfUpdateCustomDomainReturnObjJwtConfigTokenConfigResponse `json:"tokenConfig"` /*  Token 配置 */
	MatchMode   *CfUpdateCustomDomainReturnObjJwtConfigMatchModeResponse     `json:"matchMode"`   /*  JWT 匹配模式，固定为 All */
}

type CfUpdateCustomDomainReturnObjJwtConfigTokenConfigResponse struct {
	Location *string `json:"location"` /*  Token 读取位置 */
	Name     *string `json:"name"`     /*  字段名 */
}

type CfUpdateCustomDomainReturnObjJwtConfigMatchModeResponse struct {
	Mode *string   `json:"mode"` /*  匹配模式 */
	Path []*string `json:"path"` /*  匹配路径 */
}

type CfUpdateCertConfig struct {
	Certificate string `json:"certificate"` /*  HTTPS 证书内容 */
	CertName    string `json:"certName"`    /*  证书的名称 */
	PrivateKey  string `json:"privateKey"`  /*  证书私钥内容 */
}

type CfUpdateRouteConfig struct {
	Routes []CfUpdatePathConfig `json:"routes,omitempty"` /*  路由规则列表 */
}

type CfUpdatePathConfig struct {
	FunctionName       string  `json:"functionName"`        /*  路由规则对应的函数名称 */
	Path               string  `json:"path"`                /*  路由规则对应的请求路径 */
	Qualifier          *string `json:"qualifier,omitempty"` /*  路由规则对应的函数版本或别名或 LATEST */
	EnableJwt          int32   `json:"enableJwt"`           /*  是否开启 JWT 认证：0：不开；1:开启 */
	FunctionId         int32   `json:"functionId"`          /*  路由规则对应的函数 ID */
	FunctionUniqueName string  `json:"functionUniqueName"`  /*  函数的唯一名称 */
}

type CfUpdateAuthConfig struct {
	AuthType  string             `json:"authType"`            /*  认证类型：anonymousjwt */
	JwtConfig *CfUpdateJwtConfig `json:"jwtConfig,omitempty"` /*  jwt 认证配置 */
}

type CfUpdateJwtConfig struct {
	Jwks        string                `json:"jwks"`        /*  包含一个或多个 JWK 的 JSON 字符串 */
	TokenConfig []CfUpdateTokenConfig `json:"tokenConfig"` /*  token 解析的相关配置 */
	ClaimTrans  []CfUpdateClaimTran   `json:"claimTrans"`  /*  JWT 的 claim 映射到请求的某个位置 */
}

type CfUpdateTokenConfig struct {
	Location     string  `json:"location"`               /*  token 读取的位置：Cookie, Header, Query */
	Name         string  `json:"name"`                   /*  token 在读取位置对应的字段名 */
	RemovePrefix *string `json:"removePrefix,omitempty"` /*  token 需要去除的前缀（仅读取位置为 Header 时使用） */
}

type CfUpdateClaimTran struct {
	ClaimName     string `json:"claimName"`     /*  要进行映射的 claim 字段 */
	TargetName    string `json:"targetName"`    /*  映射后的字段名 */
	TransLocation string `json:"transLocation"` /*  映射到请求的位置：Header */
}

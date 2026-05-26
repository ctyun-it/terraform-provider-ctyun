package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfCreateCustomDomainApi
/* 创建自定义域名 */
type CfCreateCustomDomainApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfCreateCustomDomainApi(client *core.CtyunClient) *CfCreateCustomDomainApi {
	return &CfCreateCustomDomainApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/domains/customdomains",
			ContentType:  "application/json",
		},
	}
}

func (a *CfCreateCustomDomainApi) Do(ctx context.Context, credential core.Credential, req *CfCreateCustomDomainRequest) (*CfCreateCustomDomainResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfCreateCustomDomainRequest
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
	var resp CfCreateCustomDomainResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfCreateCustomDomainRequest struct {
	RegionId    string               `json:"regionId"`              /*  资源池 id  */
	DomainName  string               `json:"domainName"`            /*  自定义域名名称（已备案或接入备案） */
	Protocol    *string              `json:"protocol,omitempty"`    /*  域名支持的协议类型：HTTP, HTTPS */
	Description *string              `json:"description,omitempty"` /*  自定义描述 */
	CertConfig  *CfCreateCertConfig  `json:"certConfig,omitempty"`  /*  HTTPS 证书的信息 */
	RouteConfig *CfCreateRouteConfig `json:"routeConfig,omitempty"` /*  自定义域名访问时的 PATH 到 Function 映射 */
	AuthConfig  *CfCreateAuthConfig  `json:"authConfig"`            /*  认证配置信息 */
}

type CfCreateCertConfig struct {
	CertName    string `json:"certName"`    /*  证书的名称 */
	Certificate string `json:"certificate"` /*  HTTPS 证书内容 */
	PrivateKey  string `json:"privateKey"`  /*  证书私钥内容 */
}

type CfCreateRouteConfig struct {
	Routes []CfCreatePathConfig `json:"routes,omitempty"` /*  路由规则列表 */
}

type CfCreatePathConfig struct {
	EnableJwt          int32  `json:"enableJwt"`           /*  是否开启 JWT 认证：0：不开；1:开启 */
	FunctionId         int32  `json:"functionId"`          /*  路由规则对应的函数 ID */
	FunctionName       string `json:"functionName"`        /*  路由规则对应的函数名称 */
	FunctionUniqueName string `json:"functionUniqueName"`  /*  函数的唯一名称 */
	Path               string `json:"path"`                /*  路由规则对应的请求路径 */
	Qualifier          string `json:"qualifier,omitempty"` /*  路由规则对应的函数版本或别名或 LATEST */
}

type CfCreateAuthConfig struct {
	AuthType  string             `json:"authType"`            /*  认证类型：anonymousjwt */
	JwtConfig *CfCreateJwtConfig `json:"jwtConfig,omitempty"` /*  jwt 认证配置 */
}

type CfCreateJwtConfig struct {
	ClaimTrans  []CfCreateClaimTran   `json:"claimTrans"`  /*  JWT 的 claim 映射到请求的某个位置 */
	Jwks        string                `json:"jwks"`        /*  包含一个或多个 JWK 的 JSON 字符串 */
	TokenConfig []CfCreateTokenConfig `json:"tokenConfig"` /*  token 解析的相关配置 */
}

type CfCreateClaimTran struct {
	ClaimName     string `json:"claimName"`     /*  要进行映射的 claim 字段 */
	TargetName    string `json:"targetName"`    /*  映射后的字段名 */
	TransLocation string `json:"transLocation"` /*  映射到请求的位置：Header */
}

type CfCreateTokenConfig struct {
	Location     string  `json:"location"`               /*  token 读取的位置：Cookie, Header, Query */
	Name         string  `json:"name"`                   /*  token 在读取位置对应的字段名 */
	RemovePrefix *string `json:"removePrefix,omitempty"` /*  token 需要去除的前缀（仅读取位置为 Header 时使用） */
}

type CfCreateCustomDomainResponse struct {
	StatusCode *int32                                 `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Code       *string                                `json:"code"`       /*  错误码  */
	Message    *string                                `json:"message"`    /*  信息  */
	ReturnObj  *CfCreateCustomDomainReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfCreateCustomDomainReturnObjResponse struct {
	CreatedAt    *string                                           `json:"createdAt"`    /*  创建时间  */
	UpdatedAt    *string                                           `json:"updatedAt"`    /*  更新时间  */
	Creator      *int32                                            `json:"creator"`      /*  创建者 ID  */
	Editor       *int32                                            `json:"editor"`       /*  编辑者 ID  */
	CertConfig   *CfCreateCustomDomainReturnObjCertConfigResponse  `json:"certConfig"`   /*  HTTPS 证书信息  */
	DomainName   *string                                           `json:"domainName"`   /*  自定义域名  */
	Protocol     *string                                           `json:"protocol"`     /*  协议类型  */
	RouteConfig  *CfCreateCustomDomainReturnObjRouteConfigResponse `json:"routeConfig"`  /*  路由表  */
	AuthConfig   *CfCreateCustomDomainReturnObjAuthConfigResponse  `json:"authConfig"`   /*  认证配置  */
	DomainStatus *string                                           `json:"domainStatus"` /*  域名备案状态  */
	Description  *string                                           `json:"description"`  /*  描述  */
	CnameValid   *bool                                             `json:"cnameValid"`   /*  CNAME 是否有效，仅在应用场景下有意义  */
}

type CfCreateCustomDomainReturnObjCertConfigResponse struct {
	Certificate *string `json:"certificate"` /*  证书  */
	CertName    *string `json:"certName"`    /*  证书名称  */
	PrivateKey  *string `json:"privateKey"`  /*  私钥  */
}

type CfCreateCustomDomainReturnObjRouteConfigResponse struct {
	Routes []*CfCreateCustomDomainReturnObjRouteConfigRoutesResponse `json:"routes"` /*  路由映射  */
}

type CfCreateCustomDomainReturnObjRouteConfigRoutesResponse struct {
	EnableJwt          *int32  `json:"enableJwt"`          /*  是否启用 JWT  */
	FunctionId         *int32  `json:"functionId"`         /*  函数 ID  */
	FunctionName       *string `json:"functionName"`       /*  函数名称  */
	FunctionUniqueName *string `json:"functionUniqueName"` /*  函数唯一名称  */
	Path               *string `json:"path"`               /*  请求路径  */
	Qualifier          *string `json:"qualifier"`          /*  函数版本  */
}

type CfCreateCustomDomainReturnObjAuthConfigResponse struct {
	AuthType  *string                                         `json:"authType"`  /*  认证类型  */
	JwtConfig *CfCreateCustomDomainReturnObjJwtConfigResponse `json:"jwtConfig"` /*  JWT 配置  */
}

type CfCreateCustomDomainReturnObjJwtConfigResponse struct {
	Jwks        *string                                                      `json:"jwks"`        /*  JWK 字符串  */
	TokenConfig []*CfCreateCustomDomainReturnObjJwtConfigTokenConfigResponse `json:"tokenConfig"` /*  Token 配置  */
	MatchMode   *CfCreateCustomDomainReturnObjJwtConfigMatchModeResponse     `json:"matchMode"`   /*  JWT 匹配模式，固定为 All  */
	ClaimTrans  []*CfCreateCustomDomainReturnObjJwtConfigClaimTranResponse   `json:"claimTrans"`  /*  JWT 的 claim 映射  */
}

type CfCreateCustomDomainReturnObjJwtConfigTokenConfigResponse struct {
	Location     *string `json:"location"`     /*  Token 读取位置  */
	Name         *string `json:"name"`         /*  字段名  */
	RemovePrefix *string `json:"removePrefix"` /*  token 需要去除的前缀 */
}

type CfCreateCustomDomainReturnObjJwtConfigClaimTranResponse struct {
	ClaimName     *string `json:"claimName"`     /*  要进行映射的 claim 字段 */
	TargetName    *string `json:"targetName"`    /*  映射后的字段名 */
	TransLocation *string `json:"transLocation"` /*  映射到请求的位置 */
}

type CfCreateCustomDomainReturnObjJwtConfigMatchModeResponse struct {
	Mode *string   `json:"mode"` /*  匹配模式  */
	Path []*string `json:"path"` /*  匹配路径  */
}

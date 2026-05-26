package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfListCustomDomainsApi
/* 获取自定义域名列表 */
type CfListCustomDomainsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListCustomDomainsApi(client *core.CtyunClient) *CfListCustomDomainsApi {
	return &CfListCustomDomainsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/domains/customdomains",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListCustomDomainsApi) Do(ctx context.Context, credential core.Credential, req *CfListCustomDomainsRequest) (*CfListCustomDomainsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.PageIndex != nil && *req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(*req.PageIndex), 10))
	}
	if req.PageSize != nil && *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	if req.SearchKey != nil && *req.SearchKey != "" {
		ctReq.AddParam("searchKey", *req.SearchKey)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListCustomDomainsResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListCustomDomainsRequest struct {
	RegionId  string  `json:"regionId"`            /*  资源池id  */
	PageIndex *int32  `json:"pageIndex,omitempty"` /*  页码，默认为-1  */
	PageSize  *int32  `json:"pageSize,omitempty"`  /*  每页大小，默认为-1  */
	SearchKey *string `json:"searchKey,omitempty"` /*  模糊查询的关键字，默认为空  */
}

type CfListCustomDomainsResponse struct {
	StatusCode *int32                                `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string                               `json:"error"`      /*  错误码  */
	Message    *string                               `json:"message"`    /*  信息  */
	ReturnObj  *CfListCustomDomainsReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfListCustomDomainsReturnObjResponse struct {
	Data       []*CfListCustomDomainsReturnObjDataResponse     `json:"data"`       /*  分页数据  */
	Pagination *CfListCustomDomainsReturnObjPaginationResponse `json:"pagination"` /*  分页信息  */
}

type CfListCustomDomainsReturnObjDataResponse struct {
	CreatedAt    *string                                               `json:"createdAt"`    /*  创建时间  */
	UpdatedAt    *string                                               `json:"updatedAt"`    /*  更新时间  */
	Creator      *int32                                                `json:"creator"`      /*  创建者ID  */
	Editor       *int32                                                `json:"editor"`       /*  编辑者ID  */
	CertConfig   *CfListCustomDomainsReturnObjDataCertConfigResponse   `json:"certConfig"`   /*  HTTPS证书信息  */
	DomainName   *string                                               `json:"domainName"`   /*  自定义域名  */
	Protocol     *string                                               `json:"protocol"`     /*  协议类型  */
	RouteConfig  *CfListCustomDomainsReturnObjDataRouteConfigResponse  `json:"routeConfig"`  /*  路由表  */
	AuthConfig   *CfListCustomDomainsReturnObjDataAuthConfigResponse   `json:"authConfig"`   /*  认证配置  */
	DomainStatus *string                                               `json:"domainStatus"` /*  域名备案状态  */
	FilingStatus *CfListCustomDomainsReturnObjDataFilingStatusResponse `json:"filingStatus"` /*  域名备案详细信息  */
	Description  *string                                               `json:"description"`  /*  描述  */
	CnameValid   *bool                                                 `json:"cnameValid"`   /*  CNAME是否有效，仅在应用场景下有意义  */
}

type CfListCustomDomainsReturnObjPaginationResponse struct {
	PageIndex *int32 `json:"pageIndex"` /*  页码  */
	PageSize  *int32 `json:"pageSize"`  /*  每页大小  */
	Total     *int32 `json:"total"`     /*  总记录数  */
}

type CfListCustomDomainsReturnObjDataCertConfigResponse struct {
	Certificate *string `json:"certificate"` /*  证书  */
	CertName    *string `json:"certName"`    /*  证书名称  */
	PrivateKey  *string `json:"privateKey"`  /*  私钥  */
}

type CfListCustomDomainsReturnObjDataRouteConfigResponse struct {
	Routes []*CfListCustomDomainsReturnObjDataRouteConfigRoutesResponse `json:"routes"` /*  路由映射  */
}

type CfListCustomDomainsReturnObjDataAuthConfigResponse struct {
	AuthType  *string                                                      `json:"authType"`  /*  认证类型  */
	JwtConfig *CfListCustomDomainsReturnObjDataAuthConfigJwtConfigResponse `json:"jwtConfig"` /*  JWT配置  */
}

type CfListCustomDomainsReturnObjDataFilingStatusResponse struct {
	Domain       *string `json:"domain"`       /*  一级域名  */
	RecordStatus *string `json:"recordStatus"` /*  备案状态,0：未备案 1：已备案  */
	RecordNumber *string `json:"recordNumber"` /*  备案号  */
	ErrMessage   *string `json:"errMessage"`   /*  异常信息  */
}

type CfListCustomDomainsReturnObjDataRouteConfigRoutesResponse struct {
	EnableJwt          *int32    `json:"enableJwt"`          /*  是否启用JWT  */
	FunctionId         *int32    `json:"functionId"`         /*  函数ID  */
	FunctionName       *string   `json:"functionName"`       /*  函数名称  */
	FunctionUniqueName *string   `json:"functionUniqueName"` /*  函数唯一名称  */
	Methods            []*string `json:"methods"`            /*  请求方法  */
	Path               *string   `json:"path"`               /*  请求路径  */
	Qualifier          *string   `json:"qualifier"`          /*  函数版本  */
}

type CfListCustomDomainsReturnObjDataAuthConfigJwtConfigResponse struct {
	Jwks        *string                                                                   `json:"jwks"`        /*  JWK字符串  */
	TokenConfig []*CfListCustomDomainsReturnObjDataAuthConfigJwtConfigTokenConfigResponse `json:"tokenConfig"` /*  Token配置  */
	MatchMode   *CfListCustomDomainsReturnObjDataAuthConfigJwtConfigMatchModeResponse     `json:"matchMode"`   /*  JWT匹配模式，固定为All  */
}

type CfListCustomDomainsReturnObjDataAuthConfigJwtConfigTokenConfigResponse struct {
	Location *string `json:"location"` /*  Token读取位置  */
	Name     *string `json:"name"`     /*  字段名  */
}

type CfListCustomDomainsReturnObjDataAuthConfigJwtConfigMatchModeResponse struct {
	Mode *string   `json:"mode"` /*  匹配模式  */
	Path []*string `json:"path"` /*  匹配路径  */
}

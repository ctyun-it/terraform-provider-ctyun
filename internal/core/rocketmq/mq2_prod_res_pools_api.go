package rocketmq

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2ProdResPoolsApi
/* 查上线资源池 */
type Mq2ProdResPoolsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2ProdResPoolsApi(client *core.CtyunClient) *Mq2ProdResPoolsApi {
	return &Mq2ProdResPoolsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/prodResPools",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *Mq2ProdResPoolsApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *Mq2ProdResPoolsRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.Mq2ProdResPoolsApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type Mq2ProdResPoolsRequest struct {
	ProdCode      *string `json:"prodCode,omitempty"`      /*  产品编码  */
	OuterProdCode *string `json:"outerProdCode,omitempty"` /*  外部产品编码  */
	ResPoolCode   *string `json:"resPoolCode,omitempty"`   /*  资源池编码  */
}

type Mq2ProdResPoolsResponse struct {
	StatusCode    *int32                             `json:"statusCode"`    /*  返回编码。成功：800，失败：900  */
	Message       *string                            `json:"message"`       /*  提示信息  */
	ReturnObj     *string                            `json:"returnObj"`     /*  返回对象。  */
	ResPools      []*Mq2ProdResPoolsResPoolsResponse `json:"resPools"`      /*  资源池列表  */
	ResPoolCode   *string                            `json:"resPoolCode"`   /*  资源池编码  */
	ResPoolName   *string                            `json:"resPoolName"`   /*  资源池名称  */
	Products      []*Mq2ProdResPoolsProductsResponse `json:"products"`      /*  资源池产品列表  */
	ProdName      *string                            `json:"prodName"`      /*  产品名称  */
	ProdCode      *string                            `json:"prodCode"`      /*  产品编码  */
	OuterProdCode *string                            `json:"outerProdCode"` /*  外部产品编码  */
	ProdStatus    *string                            `json:"prodStatus"`    /*  产品状态： 按资源池维度 2：已上架  */
}

type Mq2ProdResPoolsResPoolsResponse struct{}

type Mq2ProdResPoolsProductsResponse struct{}

package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// HpfsInfoProtocolServiceApi
/* 查询协议服务详情
 */type HpfsInfoProtocolServiceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsInfoProtocolServiceApi(client *core.CtyunClient) *HpfsInfoProtocolServiceApi {
	return &HpfsInfoProtocolServiceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/info-protocol-service",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsInfoProtocolServiceApi) Do(ctx context.Context, credential core.Credential, req *HpfsInfoProtocolServiceRequest) (*HpfsInfoProtocolServiceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("protocolServiceID", req.ProtocolServiceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp HpfsInfoProtocolServiceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsInfoProtocolServiceRequest struct {
	RegionID          string `json:"regionID,omitempty"`          /*  资源池 ID，例：100054c0416811e9a6690242ac110002  */
	ProtocolServiceID string `json:"protocolServiceID,omitempty"` /*  协议服务唯一ID  */
}

type HpfsInfoProtocolServiceResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码(800为成功，900为处理中/失败，详见errorCode)  */
	Message     string                                    `json:"message"`     /*  响应描述  */
	Description string                                    `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsInfoProtocolServiceReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                    `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                    `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsInfoProtocolServiceReturnObjResponse struct {
	RegionID              string `json:"regionID"`              /*  资源池ID  */
	SfsUID                string `json:"sfsUID"`                /*  并行文件唯一 ID  */
	AzName                string `json:"azName"`                /*  多可用区下可用区的名字  */
	ProtocolServiceID     string `json:"protocolServiceID"`     /*  协议服务唯一ID  */
	ProtocolServiceStatus string `json:"protocolServiceStatus"` /*  协议服务的状态，creating/available//deleting/create_fail/agent_err。creating：协议服务创建中；available：协议服务可用；deleting：协议服务删除中；create_fail：协议服务创建失败；agent_err：底层协议服务组件异常（该异常状态可恢复）  */
	ProtocolSpec          string `json:"protocolSpec"`          /*  协议规格  */
	ProtocolType          string `json:"protocolType"`          /*  协议类型  */
	VpcSharePath          string `json:"vpcSharePath"`          /*  vpc挂载地址(ipv4)  */
	VpcSharePathV6        string `json:"vpcSharePathV6"`        /*  vpc挂载地址(ipv6)  */
	VpceSharePath         string `json:"vpceSharePath"`         /*  vpce挂载地址（ipv4）  */
	VpceSharePathV6       string `json:"vpceSharePathV6"`       /*  vpce挂载地址（ipv6）  */
	VpcID                 string `json:"vpcID"`                 /*  虚拟网 ID  */
	VpcName               string `json:"vpcName"`               /*  vpc名称  */
	SubnetID              string `json:"subnetID"`              /*  子网ID  */
	CreateTime            int64  `json:"createTime"`            /*  协议服务的创建时间  */
	FailMsg               string `json:"failMsg"`               /*  协议服务的异常原因  */
	ProtocolDescription   string `json:"protocolDescription"`   /*  协议服务的描述  */
}

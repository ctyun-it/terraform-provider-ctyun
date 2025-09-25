package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcL2gwConnectionCreateApi
/* 创建l2gw_connection
 */type CtvpcL2gwConnectionCreateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcL2gwConnectionCreateApi(client *core.CtyunClient) *CtvpcL2gwConnectionCreateApi {
	return &CtvpcL2gwConnectionCreateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/l2gw_connection/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcL2gwConnectionCreateApi) Do(ctx context.Context, credential core.Credential, req *CtvpcL2gwConnectionCreateRequest) (*CtvpcL2gwConnectionCreateResponse, error) {
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
	var resp CtvpcL2gwConnectionCreateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcL2gwConnectionCreateRequest struct {
	RegionID    string  `json:"regionID,omitempty"`    /*  资源池 ID  */
	Name        string  `json:"name,omitempty"`        /*  二层网关名称，二层网关名称，支持拉丁字母、中文、数字，下划线，连字符，必须以中文 / 英文字母开头，不能以数字、_和-、 http: / https: 开头，长度 2 - 32  */
	Description *string `json:"description,omitempty"` /*  描述，支持拉丁字母、中文、数字, 特殊字符：~!@#$%^& ***\*()_-+= <>?:"{},./;'[\****\]***\*·~！@#￥%……&\****（） —— -+={}\《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
	L2gwID      string  `json:"l2gwID,omitempty"`      /*  二层网关id,标准型l2gw可创建l2gw_connection三个，增强型6个  */
	SubnetID    string  `json:"subnetID,omitempty"`    /*  二层连接子网,子网和l2gw必须属于同一个vpc，一个子网只能创建1个l2gw_connection  */
	L2conIp     string  `json:"l2conIp,omitempty"`     /*  接口ip，必须是子网可用ip  */
	TunnelID    int32   `json:"tunnelID"`              /*  隧道号  */
	TunnelIp    string  `json:"tunnelIp,omitempty"`    /*  隧道ip  */
}

type CtvpcL2gwConnectionCreateResponse struct {
	StatusCode  int32                                       `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                     `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                     `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                     `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcL2gwConnectionCreateReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcL2gwConnectionCreateReturnObjResponse struct {
	ID                 *string `json:"ID,omitempty"`                 /*  l2gw_connection ID  */
	Name               *string `json:"name,omitempty"`               /*  名字  */
	Description        *string `json:"description,omitempty"`        /*  描述  */
	L2gwID             *string `json:"l2gwID,omitempty"`             /*  l2gw ID  */
	ConnectionSubnetID *string `json:"connectionSubnetID,omitempty"` /*  连接子网ID  */
	ConnectionPortIP   *string `json:"connectionPortIP,omitempty"`   /*  port ip  */
	ConnectionPortID   *string `json:"connectionPortID,omitempty"`   /*  portID  */
	RemoteVtepIp       *string `json:"remoteVtepIp,omitempty"`       /*  远端vtepip  */
	RemoteVtepPort     int32   `json:"remoteVtepPort"`               /*  远端vtepport  */
	TunnelID           int32   `json:"tunnelID"`                     /*  隧道ID  */
	CreatedAt          *string `json:"createdAt,omitempty"`          /*  创建时间  */
	UpdatedAt          *string `json:"updatedAt,omitempty"`          /*  更新时间  */
}

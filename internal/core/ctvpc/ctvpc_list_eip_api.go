package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcListEipApi
/* 调用此接口可查询指定地域已创建的弹性公网IP（Elastic IP Address，简称EIP）。
 */type CtvpcListEipApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListEipApi(client *core.CtyunClient) *CtvpcListEipApi {
	return &CtvpcListEipApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/eip/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListEipApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListEipRequest) (*CtvpcListEipResponse, error) {
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
	var resp CtvpcListEipResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListEipRequest struct {
	ClientToken string    `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string    `json:"regionID,omitempty"`    /*  资源池 ID  */
	ProjectID   *string   `json:"projectID,omitempty"`   /*  企业项目 ID，默认为"0"  */
	Page        int32     `json:"page"`                  /*  分页参数  */
	PageNo      int32     `json:"pageNo"`                /*  列表的页码，默认值为 1, 推荐使用该字段, page 后续会废弃  */
	PageSize    int32     `json:"pageSize"`              /*  每页数据量大小，取值 1-50  */
	Ids         []*string `json:"ids"`                   /*  是 Array 类型，里面的内容是 String  */
	Status      *string   `json:"status,omitempty"`      /*  eip状态 ACTIVE（已绑定）/ DOWN（未绑定）/ FREEZING（已冻结）/ EXPIRED（已过期），不传是查询所有状态的 EIP  */
	IpType      *string   `json:"ipType,omitempty"`      /*  ip类型 ipv4 / ipv6  */
	EipType     *string   `json:"eipType,omitempty"`     /*  eip类型 normal / cn2  */
	Ip          *string   `json:"ip,omitempty"`          /*  弹性 IP 的 ip 地址  */
}

type CtvpcListEipResponse struct {
	StatusCode   int32                            `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                          `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                          `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                          `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtvpcListEipReturnObjResponse `json:"returnObj"`             /*  object  */
	TotalCount   int32                            `json:"totalCount"`            /*  列表条目数  */
	CurrentCount int32                            `json:"currentCount"`          /*  分页查询时每页的行数。  */
	TotalPage    int32                            `json:"totalPage"`             /*  总页数  */
}

type CtvpcListEipReturnObjResponse struct {
	Eips []*CtvpcListEipReturnObjEipsResponse `json:"eips"` /*  弹性 IP 列表  */
}

type CtvpcListEipReturnObjEipsResponse struct {
	ID               *string `json:"ID,omitempty"`               /*  eip ID  */
	Name             *string `json:"name,omitempty"`             /*  eip 名称  */
	Description      *string `json:"description,omitempty"`      /*  描述  */
	EipAddress       *string `json:"eipAddress,omitempty"`       /*  eip 地址  */
	AssociationID    *string `json:"associationID,omitempty"`    /*  当前绑定的实例的 ID  */
	AssociationType  *string `json:"associationType,omitempty"`  /*  当前绑定的实例类型  */
	PrivateIpAddress *string `json:"privateIpAddress,omitempty"` /*  交换机网段内的一个 IP 地址  */
	Bandwidth        int32   `json:"bandwidth"`                  /*  带宽峰值大小，单位 Mb  */
	Status           *string `json:"status,omitempty"`           /*  1.ACTIVE 2.DOWN 3.ERROR 4.UPDATING 5.BANDING_OR_UNBANGDING 6.DELETING 7.DELETED 8.EXPIRED  */
	Tags             *string `json:"tags,omitempty"`             /*  EIP 的标签集合  */
	CreatedAt        *string `json:"createdAt,omitempty"`        /*  创建时间  */
	UpdatedAt        *string `json:"updatedAt,omitempty"`        /*  更新时间  */
	BandwidthID      *string `json:"bandwidthID,omitempty"`      /*  绑定的共享带宽 ID  */
	BandwidthType    *string `json:"bandwidthType,omitempty"`    /*  eip带宽规格：standalone / upflowc  */
	ExpiredAt        *string `json:"expiredAt,omitempty"`        /*  到期时间  */
	LineType         *string `json:"lineType,omitempty"`         /*  线路类型  */
	ProjectID        *string `json:"projectID,omitempty"`        /*  项目ID  */
	PortID           *string `json:"portID,omitempty"`           /*  绑定的网卡 id  */
	IsPackaged       *bool   `json:"isPackaged"`                 /*  表示是否与 vm 一起订购  */
	BillingMethod    *string `json:"billingMethod,omitempty"`    /*  计费类型：periodic 包周期，on_demand 按需  */
}

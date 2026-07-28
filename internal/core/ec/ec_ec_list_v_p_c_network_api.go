package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcListVPCNetworkApi
/* 查看VPC网络实例 */
type EcEcListVPCNetworkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcListVPCNetworkApi(client *core.CtyunClient) *EcEcListVPCNetworkApi {
	return &EcEcListVPCNetworkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/vpc-instance/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcListVPCNetworkApi) Do(ctx context.Context, credential core.Credential, req *EcEcListVPCNetworkRequest) (*EcEcListVPCNetworkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcListVPCNetworkRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcListVPCNetworkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcListVPCNetworkRequest struct {
	EcID         string  `json:"ecID"`                   /*  82c40584-fe8f-49b0-9e09-ef915caeee10  */
	CgwID        *string `json:"cgwID,omitempty"`        /*  云网关ID   */
	InstanceID   *string `json:"instanceID,omitempty"`   /*  指定网络实例ID   */
	QueryContent *string `json:"queryContent,omitempty"` /*  模糊匹配，支持VPCID，VPCName，dcName三个属性  */
	Status       *string `json:"status,omitempty"`       /*  指定状态描述<br/>取值包括：<br/>'creating'：加载中<br/>'running'：已连接<br/>'removing'：卸载中<br/>'flushing'：路由待更新<br>'error'：失败  */
	IsAuth       *int32  `json:"isAuth,omitempty"`       /*  筛选是否是跨账号实例<br/>取值包括：<br/>0:本账号<br/>1:跨账号  */
	IsExclusive  *int32  `json:"isExclusive,omitempty"`  /*  筛选是否是专属云实例<br/>取值包括：<br/>0:公有云<br/>1:专属云  */
	PageNo       *int32  `json:"pageNo,omitempty"`       /*  当前页数，从1开始  */
	PageSize     *int32  `json:"pageSize,omitempty"`     /*  当前页大小，默认10  */
}

type EcEcListVPCNetworkResponse struct {
	StatusCode  *int32                               `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                              `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                              `json:"message"`     /*  失败时的错误描述，一般为英文描述   */
	Description *string                              `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                              `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcListVPCNetworkReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcListVPCNetworkReturnObjResponse struct {
	CurrentCount *int32                                        `json:"currentCount"` /*  当前页记录数  */
	TotalPage    *int32                                        `json:"totalPage"`    /*  总页数  */
	TotalCount   *int32                                        `json:"totalCount"`   /*  查询的总记录数  */
	Results      []*EcEcListVPCNetworkReturnObjResultsResponse `json:"results"`      /*  返回查询结果，Json数组   */
}

type EcEcListVPCNetworkReturnObjResultsResponse struct {
	EcID        *string                                                 `json:"ecID"`        /*  云间高速实例ID  */
	CgwID       *string                                                 `json:"cgwID"`       /*  云网关ID  */
	CgwName     *string                                                 `json:"cgwName"`     /*  云网关名称  */
	Status      *string                                                 `json:"status"`      /*  指定状态描述<br/>取值包括：<br/>'creating'：加载中<br/>'running'：已连接<br/>'removing'：卸载中<br/>'flushing'：路由待更新<br/>'error'：失败  */
	VpcName     *string                                                 `json:"vpcName"`     /*  网络实例名称  */
	VpcID       *string                                                 `json:"vpcID"`       /*  vpc ID  */
	VpcCIDR     *string                                                 `json:"vpcCIDR"`     /*  网络实例CIDR  */
	DcName      *string                                                 `json:"dcName"`      /*  资源池名称  */
	DcID        *string                                                 `json:"dcID"`        /*  资源池ID信息  */
	DcType      *string                                                 `json:"dcType"`      /*  资源池类型，<br/>取值范围:<br/>'CNP':CNP资源池<br/>'MAZ':MAZ资源池<br/>'OS':OS资源池<br/>'CS':CS资源池<br/>'PRVT':私有云资源池  */
	RtbName     *string                                                 `json:"rtbName"`     /*  路由表名称  */
	RtbID       *string                                                 `json:"rtbID"`       /*  路由表ID  */
	ExclusiveID *string                                                 `json:"exclusiveID"` /*  专属云资源池ID  */
	InstanceID  *string                                                 `json:"instanceID"`  /*  vpc网络实例UUID  */
	IsAuth      *int32                                                  `json:"isAuth"`      /*  是否是跨账号实例<br/>取值包括：<br/>0:本账号<br/>1:跨账号  */
	IsExclusive *int32                                                  `json:"isExclusive"` /*  是否是专属云实例<br/>取值包括：<br/>0:公有云<br/>1:专属云  */
	RouteLearn  *int32                                                  `json:"routeLearn"`  /*  路由学习开关，开启后云网关自动学习网络实例路由<br/>取值范围:<br/>1:学习<br/>0:不学习<br/>默认学习  */
	RouteSync   *int32                                                  `json:"routeSync"`   /*  路由同步开关，开启后云网关路由自动同步到网络实例<br/>取值范围:<br/>1:同步<br/>0:不同步<br/>默认同步  */
	CreateDate  *string                                                 `json:"createDate"`  /*  创建时间  */
	SubnetList  []*EcEcListVPCNetworkReturnObjResultsSubnetListResponse `json:"subnetList"`  /*  子网信息  */
}

type EcEcListVPCNetworkReturnObjResultsSubnetListResponse struct {
	Gateway    *string `json:"gateway"`    /*  IPv4网关  */
	CIDR       *string `json:"CIDR"`       /*  子网IPv4 CIDR  */
	SubnetID   *string `json:"subnetID"`   /*  子网IPv4 ID   */
	IPVersion  *string `json:"IPVersion"`  /*  IP类型，取值范围<br/>IPv6：IPv6类型<br/>IPv4：IPv4类型  */
	SubnetName *string `json:"subnetName"` /*  子网名称  */
}

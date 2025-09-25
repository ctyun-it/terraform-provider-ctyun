package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtecsListPortsV41Api
/* 查询虚拟网卡列表<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权
 */type CtecsListPortsV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListPortsV41Api(client *core.CtyunClient) *CtecsListPortsV41Api {
	return &CtecsListPortsV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/ports/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListPortsV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListPortsV41Request) (*CtecsListPortsV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.VpcID != "" {
		ctReq.AddParam("vpcID", req.VpcID)
	}
	if req.DeviceID != "" {
		ctReq.AddParam("deviceID", req.DeviceID)
	}
	if req.SubnetID != "" {
		ctReq.AddParam("subnetID", req.SubnetID)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListPortsV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListPortsV41Request struct {
	RegionID   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	VpcID      string /*  虚拟私有云ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-虚拟私有云</a>来了解虚拟私有云<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4814&data=94">查询VPC列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4811&data=94">创建VPC</a><br />注：在多可用区类型资源池下，vpcID通常以“vpc-”开头，非多可用区类型资源池vpcID为uuid格式  */
	DeviceID   string /*  关联设备id，即云主机id，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a>  */
	SubnetID   string /*  所属子网ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-子网</a>来了解子网<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94">查询子网列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4812&data=94">创建子网</a>  */
	PageNumber int32  /*  页码，取值范围：正整数（≥1），注：默认值为1。建议使用pageNo，该字段未来将会下线。  */
	PageSize   int32  /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
	PageNo     int32  /*  页码，取值范围：正整数（≥1），注：默认值为1  */
}

type CtecsListPortsV41Response struct {
	StatusCode   int32                                 `json:"statusCode,omitempty"`   /*  返回状态码（800为成功，900为失败）  */
	Message      string                                `json:"message,omitempty"`      /*  英文描述信息  */
	Description  string                                `json:"description,omitempty"`  /*  中文描述信息  */
	ReturnObj    []*CtecsListPortsV41ReturnObjResponse `json:"returnObj"`              /*  接口业务数据  */
	ErrorCode    string                                `json:"errorCode,omitempty"`    /*  错误码，为product.module.code三段式码  */
	Error        string                                `json:"error,omitempty"`        /*  错误码，为product.module.code三段式码  */
	CurrentCount int32                                 `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                 `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                 `json:"totalPage,omitempty"`    /*  总页数  */
}

type CtecsListPortsV41ReturnObjResponse struct {
	NetworkInterfaceName string   `json:"networkInterfaceName,omitempty"` /*  虚拟网名称  */
	NetworkInterfaceID   string   `json:"networkInterfaceID,omitempty"`   /*  虚拟网id  */
	VpcID                string   `json:"vpcID,omitempty"`                /*  所属vpc  */
	SubnetID             string   `json:"subnetID,omitempty"`             /*  所属子网id  */
	Role                 int32    `json:"role,omitempty"`                 /*  网卡类型: 0 主网卡， 1 弹性网卡  */
	MacAddress           string   `json:"macAddress,omitempty"`           /*  mac地址  */
	PrimaryPrivateIp     string   `json:"primaryPrivateIp,omitempty"`     /*  主ip  */
	Ipv6Addresses        []string `json:"ipv6Addresses"`                  /*  ipv6地址  */
	InstanceID           string   `json:"instanceID,omitempty"`           /*  关联的设备id  */
	InstanceType         string   `json:"instanceType,omitempty"`         /*  设备类型 VM, BM, Other  */
	Description          string   `json:"description,omitempty"`          /*  描述  */
	SecurityGroupIds     []string `json:"securityGroupIds"`               /*  安全组ID列表  */
	SecondaryPrivateIps  []string `json:"secondaryPrivateIps"`            /*  辅助私网IP  */
	AdminStatus          string   `json:"adminStatus,omitempty"`          /*  是否启用DOWN, UP  */
	AssociatedEip        string   `json:"associatedEip,omitempty"`        /*  绑定弹性IP的信息  */
}

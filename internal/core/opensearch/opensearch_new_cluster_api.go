package opensearch

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

// OpensearchNewClusterApi
/* 创建 OpenSearch 集群实例<br />
<b>准备工作：</b><br />
&emsp;&emsp;构造请求：在调用前需要了解如何构造请求<br />
&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用<br />
<b>注意事项：</b><br />
&emsp;&emsp;实例名称由大小写字母、数字、下划线（_）或连字符（-）组成，且不以下划线（_）或连字符（-）开头，长度是1-32位<br />
&emsp;&emsp;密码应为数字、大写字母、小写字母、特殊符号（@$!%*#_~?）的组合，长度在12－26位<br />
*/type OpensearchNewClusterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewOpensearchNewClusterApi(client *core.CtyunClient) *OpensearchNewClusterApi {
	return &OpensearchNewClusterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/os/openapi/v1/order/new",
			ContentType:  "application/json",
		},
	}
}

func (a *OpensearchNewClusterApi) Do(ctx context.Context, credential core.Credential, req *OrderNewRequest) (*OrderNewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*OrderNewRequest
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
	var resp OrderNewResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type OrderNewRequest struct {
	AvailableZoneID string       `json:"available_zone_id"`      /*  可用区ID，选择多可用区部署时，每个可用区ID用英文逗号隔开  */
	ClusterName     string       `json:"cluster_name"`           /*  实例名称，由大小写字母、数字、下划线（_）或连字符（-）组成，且不以下划线（_）或连字符（-）开头，长度是1-32位  */
	CycleCnt        *int32       `json:"cycle_cnt,omitempty"`    /*  订购周期，当cycle_type为2时，取值范围是1-11；cycle_type为3时，取值固定为1  */
	PayType         int32        `json:"pay_type"`               /*  付费类型 ，1：包年包月，2：按需付费  */
	CycleType       *int32       `json:"cycle_type,omitempty"`   /*  订购周期，2代表按月购买，3代表按年购买  */
	RegionID        string       `json:"region_id"`              /*  资源池ID，云搜索服务当前支持的资源池参考文档  */
	VPCID           string       `json:"vpc_id"`                 /*  vpcId  */
	SubnetID        string       `json:"subnet_id"`              /*  子网id  */
	SecurityGroupID string       `json:"security_groups_id"`     /*  安全组id  */
	EnableIPv6      string       `json:"enable_ipv6"`            /*  开启IPv6：开启:OPEN 关闭:CLOSE  */
	ComponentPwd    string       `json:"component_pwd"`          /*  组件密码 ，密码应为数字、大写字母、小写字母、特殊符号（@$!%*#_~?）的组合，长度在12－26位  */
	ClusterType     int32        `json:"cluster_type"`           /*  集群类型：1：OpenSearch，2：Elasticsearch  */
	OSType          string       `json:"os_type"`                /*  操作系统类型，ctyun操作系统：CTyun、麒麟操作系统：Kylin  */
	EnableHTTPS     *string      `json:"enable_https,omitempty"` /*  不开启https: CLOSE,开启：OPEN；默认CLOSE  */
	AutoPay         bool         `json:"auto_pay,omitempty"`     /*  true:自动支付，false:不自动支付；默认false  */
	NodeDetails     []NodeDetail `json:"node_details"`           /*  节点组详情  */
}

type NodeDetail struct {
	HostNum        int32  `json:"host_num"`          /*  节点数：MASTER 最小为3，最大为50；EXCLUSIVE_MASTER 最大为3；COORDINATE 最大为32；COLD 最大为50；  */
	IOType         string `json:"io_type"`           /*  IO类型：SSD-genric（通用型SSD）、SAS（高IO）、SSD（超高IO）、XSSD-0、XSSD-1  */
	Volume         int32  `json:"volume"`            /*  存储容量：MASTER 节点可选：40 - 6144GB；EXCLUSIVE_MASTER 节点可选：固定 40 GB；COORDINATE 节点可选：固定 40 GB；COLD 节点可选：40 - 6144GB  */
	IaasVMSpecCode string `json:"iaas_vm_spec_code"` /*  实例规格code，每个资源池可用区下可选择的机型，参考订购页面展示的机型信息  */
	NodeGroupType  string `json:"node_group_type"`   /*  节点组类型：MASTER（数据节点）/EXCLUSIVE_MASTER（专属master节点）/COORDINATE（专属协调节点）/COLD（冷数据节点）  */
}

type OrderNewResponse struct {
	StatusCode int32           `json:"statusCode"` /*  状态码，成功："200"，失败："500"  */
	Error      string          `json:"error"`      /*  错误码，请求成功时，不返回该字段  */
	Message    string          `json:"message"`    /*  用来简述当前接口调用状态以及必要提示信息  */
	ReturnObj  json.RawMessage `json:"returnObj"`  /*  返回结果  */
}

//type ReturnObj struct {
//	OrderNo string `json:"orderNo"` /*  订单号  */
//}

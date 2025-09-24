package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingConfigCreateApi
/* 创建一个弹性伸缩配置<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingConfigCreateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingConfigCreateApi(client *core.CtyunClient) *ScalingConfigCreateApi {
	return &ScalingConfigCreateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/config-create",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingConfigCreateApi) Do(ctx context.Context, credential core.Credential, req *ScalingConfigCreateRequest) (*ScalingConfigCreateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingConfigCreateRequest
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
	var resp ScalingConfigCreateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingConfigCreateRequest struct {
	RegionID            string                               `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	Name                string                               `json:"name,omitempty"`         /*  伸缩配置名称  */
	ImageID             string                               `json:"imageID,omitempty"`      /*  镜像ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10030151">镜像概述</a>来了解云主机镜像<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89">查询可以使用的镜像资源</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4765&data=89">创建私有镜像（云主机系统盘）</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5230&data=89">创建私有镜像（云主机数据盘）</a>  */
	SecurityGroupIDList []string                             `json:"securityGroupIDList"`    /*  安全组ID列表，非多可用区资源池该参数为必填，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>了解安全组相关信息 <br />获取： <br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a>  */
	SpecName            string                               `json:"specName,omitempty"`     /*  规格名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10118193">规格说明</a>了解弹性云主机的选型基本信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8327&data=87">查询一个或多个云主机规格资源</a>  */
	Volumes             []*ScalingConfigCreateVolumesRequest `json:"volumes"`                /*  磁盘类型和大小列表，元素为volume  */
	UseFloatings        int32                                `json:"useFloatings,omitempty"` /*  是否使用弹性IP。<br> 取值范围：<br>1：不使用。<br>2：自动分配。  */
	BandWidth           int32                                `json:"bandWidth,omitempty"`    /*  弹性IP带宽，单位：Mbps，useFloatings为2时必填，范围1-3000  */
	LoginMode           int32                                `json:"loginMode,omitempty"`    /*  登录方式。<br> 取值范围：<br>1：密码。<br>2：密钥对。  */
	Username            string                               `json:"username,omitempty"`     /*  用户名，loginMode为1时，必填  */
	Password            string                               `json:"password,omitempty"`     /*  密码，loginMode为1时，必填  */
	KeyPairID           string                               `json:"keyPairID,omitempty"`    /*  密钥对ID，loginMode为2时，必填，您可以查看<a href="https://www.ctyun.cn/document/10026730/10230540">密钥对</a>来了解密钥对相关内容 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8342&data=87">查询一个或多个密钥对</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8344&data=87">创建一对SSH密钥对</a>  */
	Tags                []*ScalingConfigCreateTagsRequest    `json:"tags"`                   /*  标签集  */
	AzNames             []string                             `json:"azNames"`                /*  可用区列表，仅多可用区资源池支持，多可用区资源池该参数为必填，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a>  */
	MonitorService      *bool                                `json:"monitorService"`         /*  是否开启详细监控，默认开启。<br> 取值范围：<br>true：开启。<br>false：关闭  */
}

type ScalingConfigCreateVolumesRequest struct {
	VolumeType string `json:"volumeType,omitempty"` /*  磁盘类型： SATA/SAS/SSD/SATA-KUNPENG/SATA-HAIGUANG/SAS-KUNPENG/SAS-HAIGUANG/SSD-genric。不同资源池可配置的volumeType有差异，详细请参考云硬盘  */
	VolumeSize int32  `json:"volumeSize,omitempty"` /*  磁盘大小,单位G  */
	DiskMode   string `json:"diskMode,omitempty"`   /*  数据盘磁盘模式，默认为VBD <br> 取值范围：<br>VBD（虚拟块存储设备）。<br>ISCSI（小型计算机系统接口）  */
	Flag       int32  `json:"flag,omitempty"`       /*  磁盘类型。<br/>取值范围：<br/>1：系统盘。<br/>2：数据盘。<br/>系统盘限制为1块。  */
}

type ScalingConfigCreateTagsRequest struct {
	Key   string `json:"key,omitempty"`   /*  标签键  */
	Value string `json:"value,omitempty"` /*  标签值  */
}

type ScalingConfigCreateResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingConfigCreateReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                                `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingConfigCreateReturnObjResponse struct {
	Id int64 `json:"id"` /*  伸缩配置ID  */
}

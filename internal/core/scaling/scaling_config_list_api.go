package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingConfigListApi
/* 查询弹性伸缩配置<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingConfigListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingConfigListApi(client *core.CtyunClient) *ScalingConfigListApi {
	return &ScalingConfigListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/config-list",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingConfigListApi) Do(ctx context.Context, credential core.Credential, req *ScalingConfigListRequest) (*ScalingConfigListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingConfigListRequest
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
	var resp ScalingConfigListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingConfigListRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	ConfigID int64  `json:"configID,omitempty"` /*  伸缩配置ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5068&data=93">查询弹性伸缩配置</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4995&data=93">创建一个弹性伸缩配置  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  页码  */
	Page     int32  `json:"page,omitempty"`     /*  【Deprecated】页码  */
	PageSize int32  `json:"pageSize,omitempty"` /*  分页查询时设置的每页行数，取值范围:[1~100]，默认值为10  */
}

type ScalingConfigListResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                `json:"description"` /*  失败时的错误描述,一般为中文描述  */
	ReturnObj   []*ScalingConfigListReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                                `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingConfigListReturnObjResponse struct {
	ConfigID          int32                                        `json:"configID"`          /*  伸缩配置ID  */
	Name              string                                       `json:"name"`              /*  伸缩配置名称  */
	RegionID          string                                       `json:"regionID"`          /*  资源池ID  */
	Visibility        int32                                        `json:"visibility"`        /*  镜像类型。<br>取值范围：<br> 1：公有镜像<br>0：私有镜像  */
	ImageName         string                                       `json:"imageName"`         /*  镜像名称  */
	ImageID           string                                       `json:"imageID"`           /*  镜像ID  */
	SecurityGroupList []*string                                    `json:"securityGroupList"` /*  安全组ID  */
	Cpu               int32                                        `json:"cpu"`               /*  CPU核数  */
	Memory            int32                                        `json:"memory"`            /*  内存，单位：G  */
	SpecName          string                                       `json:"specName"`          /*  规格名称  */
	OsType            int32                                        `json:"osType"`            /*  镜像系统类型。<br>取值范围：<br>1：Linux<br>2：Windows  */
	Volumes           []*ScalingConfigListReturnObjVolumesResponse `json:"volumes"`           /*  磁盘类型和大小列表，元素为volume  */
	UseFloatings      int32                                        `json:"useFloatings"`      /*  是否使用弹性IP。<br> 取值范围：<br>1：不使用。<br>2：自动分配。  */
	Bandwidth         int32                                        `json:"bandwidth"`         /*  带宽，单位：Mbps  */
	LoginMode         int32                                        `json:"loginMode"`         /*  登录方式。<br>取值范围：<br>1：密码。<br>2：密钥对  */
	Username          string                                       `json:"username"`          /*  用户名，loginMode为1时，必填  */
	GroupCount        int32                                        `json:"groupCount"`        /*  绑定的伸缩组个数  */
	Tags              []*ScalingConfigListReturnObjTagsResponse    `json:"tags"`              /*  标签集  */
	AzNames           string                                       `json:"azNames"`           /*  可用区名称  */
	MonitorService    *bool                                        `json:"monitorService"`    /*  是否开启详细监控，默认开启。<br> 取值范围：<br>true：开启。<br>false：关闭  */
}

type ScalingConfigListReturnObjSecurityGroupListResponse struct{}

type ScalingConfigListReturnObjVolumesResponse struct {
	VolumeType string `json:"volumeType"` /*  磁盘类型： SATA/SAS/SSD/SATA-KUNPENG/SATA-HAIGUANG/SAS-KUNPENG/SAS-HAIGUANG/SSD-genric。不同资源池可配置的volumeType有差异，详细请参考云硬盘  */
	VolumeSize int32  `json:"volumeSize"` /*  磁盘大小  */
	DiskMode   string `json:"diskMode"`   /*  数据盘磁盘模式，默认为VBD <br> 取值范围：<br>VBD（虚拟块存储设备）。<br>ISCSI（小型计算机系统接口）  */
	Flag       int32  `json:"flag"`       /*  磁盘类型。<br>取值范围：<br>1：系统盘。<br>2：数据盘。<br>系统盘限制为1块。  */
}

type ScalingConfigListReturnObjTagsResponse struct {
	Key   string `json:"key"`   /*  标签键  */
	Value string `json:"value"` /*  标签值  */
}

package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsTotalInstanceBackupVolumeSizeV41Api
/* 云主机备份查询虚机磁盘大小。<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsTotalInstanceBackupVolumeSizeV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsTotalInstanceBackupVolumeSizeV41Api(client *core.CtyunClient) *CtecsTotalInstanceBackupVolumeSizeV41Api {
	return &CtecsTotalInstanceBackupVolumeSizeV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/backup/statistics",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsTotalInstanceBackupVolumeSizeV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsTotalInstanceBackupVolumeSizeV41Request) (*CtecsTotalInstanceBackupVolumeSizeV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("instanceID", req.InstanceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsTotalInstanceBackupVolumeSizeV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsTotalInstanceBackupVolumeSizeV41Request struct {
	RegionID   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a>  */
}

type CtecsTotalInstanceBackupVolumeSizeV41Response struct {
	StatusCode  int32                                                   `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                                  `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                                  `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                                  `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                                  `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsTotalInstanceBackupVolumeSizeV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsTotalInstanceBackupVolumeSizeV41ReturnObjResponse struct {
	TotalDiskSize int32 `json:"totalDiskSize,omitempty"` /*  云主机磁盘占用大小，单位GB  */
}

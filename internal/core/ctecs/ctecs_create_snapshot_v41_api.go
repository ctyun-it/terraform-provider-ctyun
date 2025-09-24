package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreateSnapshotV41Api
/* 创建云主机快照<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项</b>：<br />&emsp;&emsp;云主机快照的配额限制：单台云主机可创建的快照配额数量。确认个人在不同资源池下资源配额，可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9714&data=87">用户配额查询</a>接口进行查询<br />
 */type CtecsCreateSnapshotV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreateSnapshotV41Api(client *core.CtyunClient) *CtecsCreateSnapshotV41Api {
	return &CtecsCreateSnapshotV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/snapshot/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreateSnapshotV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsCreateSnapshotV41Request) (*CtecsCreateSnapshotV41Response, error) {
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
	var resp CtecsCreateSnapshotV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreateSnapshotV41Request struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID   string `json:"instanceID,omitempty"`   /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	SnapshotName string `json:"snapshotName,omitempty"` /*  云主机快照名称。满足以下规则：长度为2-63字符，头尾不支持输入空格  */
}

type CtecsCreateSnapshotV41Response struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                   `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                   `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsCreateSnapshotV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsCreateSnapshotV41ReturnObjResponse struct {
	JobID          string `json:"jobID,omitempty"`          /*  任务ID  */
	SnapshotStatus string `json:"snapshotStatus,omitempty"` /*  云主机快照状态：<br />pending：创建中, <br />available：可用， <br />restoring：恢复中，<br />error：错误  */
	InstanceID     string `json:"instanceID,omitempty"`     /*  云主机ID  */
	InstanceName   string `json:"instanceName,omitempty"`   /*  云主机名称  */
	SnapshotID     string `json:"snapshotID,omitempty"`     /*  云主机快照ID  */
	ProjectID      string `json:"projectID,omitempty"`      /*  企业项目ID  */
	SnapshotName   string `json:"snapshotName,omitempty"`   /*  云主机快照名称  */
}

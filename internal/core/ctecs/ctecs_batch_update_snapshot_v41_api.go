package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsBatchUpdateSnapshotV41Api
/* 批量更改云主机快照名称和描述<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsBatchUpdateSnapshotV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsBatchUpdateSnapshotV41Api(client *core.CtyunClient) *CtecsBatchUpdateSnapshotV41Api {
	return &CtecsBatchUpdateSnapshotV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/snapshot/batch-update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsBatchUpdateSnapshotV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsBatchUpdateSnapshotV41Request) (*CtecsBatchUpdateSnapshotV41Response, error) {
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
	var resp CtecsBatchUpdateSnapshotV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsBatchUpdateSnapshotV41Request struct {
	RegionID   string                                          `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	UpdateInfo []*CtecsBatchUpdateSnapshotV41UpdateInfoRequest `json:"updateInfo"`         /*  批量更新信息列表  */
}

type CtecsBatchUpdateSnapshotV41UpdateInfoRequest struct {
	SnapshotID          string `json:"snapshotID,omitempty"`          /*  云主机快照ID，<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8349&data=87">查询云主机快照列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8352&data=87">创建云主机快照</a>  */
	SnapshotName        string `json:"snapshotName,omitempty"`        /*  云主机快照名称。满足以下规则：长度为2-63字符，头尾不支持输入空格  */
	SnapshotDescription string `json:"snapshotDescription,omitempty"` /*  云主机快照描述，字符长度不超过256字符  */
}

type CtecsBatchUpdateSnapshotV41Response struct {
	StatusCode  int32                                         `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                        `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                        `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                        `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsBatchUpdateSnapshotV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsBatchUpdateSnapshotV41ReturnObjResponse struct {
	SnapshotIDList []string `json:"snapshotIDList"` /*  云主机快照ID  */
}

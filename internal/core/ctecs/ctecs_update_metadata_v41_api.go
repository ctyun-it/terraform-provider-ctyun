package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsUpdateMetadataV41Api
/* 为云主机更新元数据：<br />&emsp;&emsp;如果元数据中没有待更新字段，则自动添加该字段<br />&emsp;&emsp;如果元数据中已存在待更新字段，则直接更新字段值<br />&emsp;&emsp;如果元数据中的字段不在请求参数中，则保持不变<br />&emsp;&emsp;如果isForce为true则覆盖原有的元数据信息<br /> 元数据的使用参考<a href="https://www.ctyun.cn/document/10026730/10387544">云主机自定义元数据</a><br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br /><b>注意事项</b><br />&emsp;&emsp;在云主机绑定委托情况下，注意委托对应的元数据键，不要覆盖更新，尤其慎用isForce参数
 */type CtecsUpdateMetadataV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsUpdateMetadataV41Api(client *core.CtyunClient) *CtecsUpdateMetadataV41Api {
	return &CtecsUpdateMetadataV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/metadata/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsUpdateMetadataV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsUpdateMetadataV41Request) (*CtecsUpdateMetadataV41Response, error) {
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
	var resp CtecsUpdateMetadataV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsUpdateMetadataV41Request struct {
	RegionID   string                                 `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br />  <span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string                                 `json:"instanceID,omitempty"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	Metadata   *CtecsUpdateMetadataV41MetadataRequest `json:"metadata"`             /*  元数据信息，用户自定义metadata键值对  */
	IsForce    *bool                                  `json:"isForce"`              /*  是否覆盖原有的元数据信息。缺省为不覆盖，新增元数据信息  */
}

type CtecsUpdateMetadataV41MetadataRequest struct{}

type CtecsUpdateMetadataV41Response struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                   `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                   `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsUpdateMetadataV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsUpdateMetadataV41ReturnObjResponse struct {
	InstanceID string                                           `json:"instanceID,omitempty"` /*  云主机ID  */
	Metadata   *CtecsUpdateMetadataV41ReturnObjMetadataResponse `json:"metadata"`             /*  元数据信息  */
}

type CtecsUpdateMetadataV41ReturnObjMetadataResponse struct{}

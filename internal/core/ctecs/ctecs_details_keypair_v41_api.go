package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDetailsKeypairV41Api
/* 此接口提供用户查询SSH密钥对功能。系统会接收用户输入的查询条件，并返回符合条件的密钥对详细信息。用户可根据此接口的返回值了解对应条件下的密钥对信息<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsDetailsKeypairV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDetailsKeypairV41Api(client *core.CtyunClient) *CtecsDetailsKeypairV41Api {
	return &CtecsDetailsKeypairV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/keypair/details",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDetailsKeypairV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsDetailsKeypairV41Request) (*CtecsDetailsKeypairV41Response, error) {
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
	var resp CtecsDetailsKeypairV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDetailsKeypairV41Request struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	ProjectID    string `json:"projectID,omitempty"`    /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
	KeyPairName  string `json:"keyPairName,omitempty"`  /*  密钥对名称。满足以下规则：只能由数字、字母、-组成，不能以数字和-开头、以-结尾，且长度为2-63字符.  */
	QueryContent string `json:"queryContent,omitempty"` /*  模糊匹配查询内容（匹配字段：keyPairName、keyPairID）  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsDetailsKeypairV41Response struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中或失败)  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                   `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                   `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsDetailsKeypairV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsDetailsKeypairV41ReturnObjResponse struct {
	CurrentCount int32                                             `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                             `json:"totalCount,omitempty"`   /*  总记录数  */
	Results      []*CtecsDetailsKeypairV41ReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsDetailsKeypairV41ReturnObjResultsResponse struct {
	PublicKey   string `json:"publicKey,omitempty"`   /*  密钥对的公钥  */
	KeyPairName string `json:"keyPairName,omitempty"` /*  密钥对名称  */
	FingerPrint string `json:"fingerPrint,omitempty"` /*  密钥对的指纹，采用MD5信息摘要算法  */
	KeyPairID   string `json:"keyPairID,omitempty"`   /*  密钥对的ID  */
	ProjectID   string `json:"projectID,omitempty"`   /*  企业项目ID  */
}

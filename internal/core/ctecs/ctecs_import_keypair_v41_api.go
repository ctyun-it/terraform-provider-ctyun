package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsImportKeypairV41Api
/* 导入由其他工具产生的RSA密钥对的公钥部分，密钥对的类型必须是SSH或x509<br /><b>准备工作</b>：&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsImportKeypairV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsImportKeypairV41Api(client *core.CtyunClient) *CtecsImportKeypairV41Api {
	return &CtecsImportKeypairV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/keypair/import-keypair",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsImportKeypairV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsImportKeypairV41Request) (*CtecsImportKeypairV41Response, error) {
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
	var resp CtecsImportKeypairV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsImportKeypairV41Request struct {
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	KeyPairName string `json:"keyPairName,omitempty"` /*  密钥对名称。满足以下规则：只能由数字、字母、-组成，不能以数字和-开头、以-结尾，且长度为2-63字符  */
	PublicKey   string `json:"publicKey,omitempty"`   /*  导入的公钥信息。最多支持1024字符长度（包括1024字符）的公钥导入；仅支持RSA类型的密钥  */
	ProjectID   string `json:"projectID,omitempty"`   /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"  */
}

type CtecsImportKeypairV41Response struct {
	StatusCode  int32                                   `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                  `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                  `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                  `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsImportKeypairV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsImportKeypairV41ReturnObjResponse struct {
	PublicKey   string `json:"publicKey,omitempty"`   /*  密钥对的公钥  */
	KeyPairName string `json:"keyPairName,omitempty"` /*  密钥对名称  */
	FingerPrint string `json:"fingerPrint,omitempty"` /*  密钥对的指纹，采用MD5信息摘要算法  */
}

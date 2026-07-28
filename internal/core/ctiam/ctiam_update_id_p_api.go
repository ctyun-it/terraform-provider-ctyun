package ctiam

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamUpdateIdPApi
/* 编辑身份供应商 */
type CtiamUpdateIdPApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamUpdateIdPApi(client *core.CtyunClient) *CtiamUpdateIdPApi {
	return &CtiamUpdateIdPApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/identityProvider/updateIdP",
			ContentType:  "multipart/form-data",
		},
	}
}

func (api *CtiamUpdateIdPApi) Do(credential *core.Credential, apis *Apis, yourEndpoint string, request *CtiamUpdateIdPRequest) {
	header := core.StructToHeader(request)
	headerMap := core.String2Map(header)
	var fileMap map[string]string
	fileMap = core.StructToFileMap(request)
	var dataMap map[string]string
	dataMap = make(map[string]string, 3)
	core.PostHttpForFormData("multipart/form-data", yourEndpoint+apis.CtiamUpdateIdPApi.template.UrlPath, credential.GetAccessKey(), credential.GetSecretKey(), headerMap, fileMap, dataMap)
}

type CtiamUpdateIdPRequest struct {
	Id       string  `json:"id"`                 /*  ID  */
	File     []*int8 `json:"file,omitempty"`     /*  文件  */
	FileName *string `json:"fileName,omitempty"` /*  文件名称（需携带后缀）  */
	Remark   *string `json:"remark,omitempty"`   /*  描述  */
}

type CtiamUpdateIdPResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}

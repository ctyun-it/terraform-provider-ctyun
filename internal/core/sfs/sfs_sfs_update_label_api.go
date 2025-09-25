package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsUpdateLabelApi
/* 为指定文件系统实例添加标签，支持添加单个和多个标签。
 */type SfsSfsUpdateLabelApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsUpdateLabelApi(client *core.CtyunClient) *SfsSfsUpdateLabelApi {
	return &SfsSfsUpdateLabelApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/update-label",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsUpdateLabelApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsUpdateLabelRequest) (*SfsSfsUpdateLabelResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsUpdateLabelRequest
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
	var resp SfsSfsUpdateLabelResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsUpdateLabelRequest struct {
	SfsUID    string                               `json:"sfsUID,omitempty"`   /*  弹性文件唯一资源id  */
	RegionID  string                               `json:"regionID,omitempty"` /*  资源池（区域）ID  */
	LabelList []*SfsSfsUpdateLabelLabelListRequest `json:"labelList"`          /*  标签和相应的操作类型  */
}

type SfsSfsUpdateLabelLabelListRequest struct {
	Key         string `json:"key,omitempty"`         /*  标签键。不能为空首字符不能为空格，长度不超过32个字符  */
	Value       string `json:"value,omitempty"`       /*  标签值。不能为空首字符不能为空格，长度不超过32个字符  */
	OperateType string `json:"operateType,omitempty"` /*  操作类型   绑定 BIND   解绑 UNBIND  */
}

type SfsSfsUpdateLabelResponse struct {
	ReturnObj   *SfsSfsUpdateLabelReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	StatusCode  int32                               `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     string                              `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                              `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string                              `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]，出错时才返回  */
	Error       string                              `json:"error"`       /*  业务细分码，为Product.Module.Code三段式码大驼峰形式，出错时才返回  */
}

type SfsSfsUpdateLabelReturnObjResponse struct {
	Code    string `json:"code"`    /*  返回状态码（800为成功，其他为失败）  */
	Success *bool  `json:"success"` /*  true为成功，false为失败  */
}

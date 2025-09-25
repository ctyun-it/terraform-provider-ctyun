package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsAddMountpointSfsApi
/* 根据资源池ID和文件系统ID，为指定文件系统添加挂载点
 */type SfsAddMountpointSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsAddMountpointSfsApi(client *core.CtyunClient) *SfsAddMountpointSfsApi {
	return &SfsAddMountpointSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/add-mountpoint-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsAddMountpointSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsAddMountpointSfsRequest) (*SfsAddMountpointSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsAddMountpointSfsRequest
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
	var resp SfsAddMountpointSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsAddMountpointSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一ID  */
	VpcID    string `json:"vpcID,omitempty"`    /*  vpcID。可以通过<a href="https://www.ctyun.cn/document/10026755/10040788" target="_blank">查询VPC列表</a>获取，如需新增可以<a href="https://www.ctyun.cn/document/10026755/10040800" target="_blank">创建VPC</a>。也可以登录官网，在“控制中心-网络-虚拟私有云”控制台查询vpcID，具体请参考<a href="https://www.ctyun.cn/document/10026755" target="_blank">虚拟私有云</a>。  */
	SubnetID string `json:"subnetID,omitempty"` /*  子网ID。可以通过<a href="https://www.ctyun.cn/document/10026755/10040797" target="_blank">查询子网列表</a>获取，如需新增可以<a href="https://www.ctyun.cn/document/10026755/10040804" target="_blank">创建子网</a>。也可以登录官网，在“控制中心-网络-虚拟私有云”控制台查询subnetID，具体请参考<a href="https://www.ctyun.cn/document/10026755" target="_blank">虚拟私有云</a>。  */
}

type SfsAddMountpointSfsResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsAddMountpointSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsAddMountpointSfsReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  添加挂载点的操作ID  */
}

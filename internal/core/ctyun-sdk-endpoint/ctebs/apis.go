package ctebs

import (
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// Apis api的接口
type Apis struct {
	EbsCreateApi       *EbsCreateApi
	EbsDeleteApi       *EbsDeleteApi
	EbsChangeNameApi   *EbsChangeNameApi
	EbsChangeSizeApi   *EbsChangeSizeApi
	EbsAssociateApi    *EbsAssociateApi
	EbsDisassociateApi *EbsDisassociateApi
	EbsShowApi         *EbsShowApi
}

// NewApis 构建
func NewApis(client *ctyunsdk.CtyunClient) *Apis {
	builder := ctyunsdk.NewApiHookBuilder()
	for _, hook := range client.Config.ApiHooks {
		builder.AddHooks(hook)
	}
	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointCtebsTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentTest, EndpointCtebsTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentProd, EndpointCtebsProd)
	return &Apis{
		EbsCreateApi:       NewEbsCreateApi(client),
		EbsDeleteApi:       NewEbsDeleteApi(client),
		EbsChangeNameApi:   NewEbsChangeNameApi(client),
		EbsChangeSizeApi:   NewEbsChangeSizeApi(client),
		EbsAssociateApi:    NewEbsAssociateApi(client),
		EbsDisassociateApi: NewEbsDisassociateApi(client),
		EbsShowApi:         NewEbsShowApi(client),
	}
}

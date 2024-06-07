package ctvpc

import (
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// Apis api的接口
type Apis struct {
	VpcCreateApi                      *VpcCreateApi
	VpcUpdateApi                      *VpcUpdateApi
	VpcDeleteApi                      *VpcDeleteApi
	VpcListApi                        *VpcListApi
	VpcQueryApi                       *VpcQueryApi
	SubnetCreateApi                   *SubnetCreateApi
	SubnetUpdateApi                   *SubnetUpdateApi
	SubnetDeleteApi                   *SubnetDeleteApi
	SubnetListApi                     *SubnetListApi
	SubnetQueryApi                    *SubnetQueryApi
	EipCreateApi                      *EipCreateApi
	EipDeleteApi                      *EipDeleteApi
	EipModifySpecApi                  *EipModifySpecApi
	EipChangeNameApi                  *EipChangeNameApi
	EipAssociateApi                   *EipAssociateApi
	EipDisassociateApi                *EipDisassociateApi
	EipShowApi                        *EipShowApi
	SecurityGroupCreateApi            *SecurityGroupCreateApi
	SecurityGroupModifyAttributionApi *SecurityGroupModifyAttributionApi
	SecurityGroupDeleteApi            *SecurityGroupDeleteApi
	SecurityGroupDescribeAttributeApi *SecurityGroupDescribeAttributeApi
	SecurityGroupRuleEgressRevokeApi  *SecurityGroupRuleEgressRevokeApi
	SecurityGroupRuleEgressModifyApi  *SecurityGroupRuleEgressModifyApi
	SecurityGroupRuleIngressRevokeApi *SecurityGroupRuleIngressRevokeApi
	SecurityGroupRuleIngressModifyApi *SecurityGroupRuleIngressModifyApi
	SecurityGroupRuleDescribeApi      *SecurityGroupRuleDescribeApi
	BandwidthDescribeApi              *BandwidthDescribeApi
	BandwidthCreateApi                *BandwidthCreateApi
	BandwidthChangeNameApi            *BandwidthChangeNameApi
	BandwidthChangeSpecApi            *BandwidthChangeSpecApi
	BandwidthDeleteApi                *BandwidthDeleteApi
	BandwidthAssociateEipApi          *BandwidthAssociateEipApi
	BandwidthDisassociateEipApi       *BandwidthDisassociateEipApi
}

// NewApis 构建
func NewApis(client *ctyunsdk.CtyunClient) *Apis {
	builder := ctyunsdk.NewApiHookBuilder()
	for _, hook := range client.Config.ApiHooks {
		builder.AddHooks(hook)
	}
	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointCtvpcTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentTest, EndpointCtvpcTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentProd, EndpointCtvpcProd)
	return &Apis{
		VpcCreateApi:                      NewVpcCreateApi(client),
		VpcUpdateApi:                      NewVpcUpdateApi(client),
		VpcDeleteApi:                      NewVpcDeleteApi(client),
		VpcListApi:                        NewVpcListApi(client),
		VpcQueryApi:                       NewVpcQueryApi(client),
		SubnetCreateApi:                   NewSubnetCreateApi(client),
		SubnetUpdateApi:                   NewSubnetUpdateApi(client),
		SubnetDeleteApi:                   NewSubnetDeleteApi(client),
		SubnetListApi:                     NewSubnetListApi(client),
		SubnetQueryApi:                    NewSubnetQueryApi(client),
		EipCreateApi:                      NewEipCreateApi(client),
		EipDeleteApi:                      NewEipDeleteApi(client),
		EipModifySpecApi:                  NewEipModifySpecApi(client),
		EipChangeNameApi:                  NewEipChangeNameApi(client),
		EipAssociateApi:                   NewEipAssociateApi(client),
		EipDisassociateApi:                NewEipDisassociateApi(client),
		EipShowApi:                        NewEipShowApi(client),
		SecurityGroupCreateApi:            NewSecurityGroupCreateApi(client),
		SecurityGroupModifyAttributionApi: NewSecurityGroupModifyAttributionApi(client),
		SecurityGroupDeleteApi:            NewSecurityGroupDeleteApi(client),
		SecurityGroupDescribeAttributeApi: NewSecurityGroupDescribeAttributeApi(client),
		SecurityGroupRuleEgressRevokeApi:  NewSecurityGroupRuleEgressRevokeApi(client),
		SecurityGroupRuleEgressModifyApi:  NewSecurityGroupRuleEgressModifyApi(client),
		SecurityGroupRuleIngressRevokeApi: NewSecurityGroupRuleIngressRevokeApi(client),
		SecurityGroupRuleIngressModifyApi: NewSecurityGroupRuleIngressModifyApi(client),
		SecurityGroupRuleDescribeApi:      NewSecurityGroupRuleDescribeApi(client),
		BandwidthDescribeApi:              NewBandwidthDescribeApi(client),
		BandwidthCreateApi:                NewBandwidthCreateApi(client),
		BandwidthChangeNameApi:            NewBandwidthChangeNameApi(client),
		BandwidthChangeSpecApi:            NewBandwidthChangeSpecApi(client),
		BandwidthDeleteApi:                NewBandwidthDeleteApi(client),
		BandwidthAssociateEipApi:          NewBandwidthAssociateEipApi(client),
		BandwidthDisassociateEipApi:       NewBandwidthDisassociateEipApi(client),
	}
}

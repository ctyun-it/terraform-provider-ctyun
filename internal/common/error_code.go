package common

import (
	"errors"
)

const (
	OpenapiSecurityGroupRuleNotFound = "Openapi.SecurityGroupRule.NotFound"
	OpenapiOrderInprogress           = "Openapi.Order.Inprogress"
	EcsInstanceNotFound              = "Ecs.Instance.NotFound"
	EcsInstanceStatusNotRunning      = "Ecs.Instance.StatusNotRunning"
	EcsInstanceStatusNotStopped      = "Ecs.Instance.StatusNotStopped"
	ImageImageCheckNotFound          = "Image.ImageCheck.NotFound"
	OpenapiEipNotFound               = "Openapi.Eip.NotFound"
	EbsEbsInfoDataDamaged            = "ebs.ebsInfo.dataDamaged"
	OpenapiSecurityGroupNotFound     = "Openapi.SecurityGroup.NotFound"
	OpenapiSharedbandwidthNotFound   = "Openapi.Sharedbandwidth.NotFound"
	EbsOrderInProgress               = "ebs.order.inProgress"
	OpenapiVpcNotFound               = "Openapi.Vpc.NotFound"
	OpenapiSubnetNotFound            = "Openapi.Subnet.NotFound"
	EcsAffinityGroupNotBound         = "Ecs.AffinityGroup.NotBound"
	OpenapiRouterTableAccessFailed   = "Openapi.RouterTable.AccessFailed"
	OpenapiVpceEndpointNotFound      = "Openapi.VpceEndpoint.NotFound"
	CtiamNoPermission                = "CTIAM_0005"
	CtiamNoPrivilege                 = "CTIAM_1044"

	ErrorStatusCode        = 900
	NormalStatusCode       = 800
	NormalStatusCodeString = "800"
)

var InvalidReturnObjError = errors.New("invalid return object")
var InvalidReturnObjResultsError = errors.New("invalid result object results")

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
	EbsEbsInfoNotExists              = "ebs.ebsInfo.notExists"
	OpenapiSecurityGroupNotFound     = "Openapi.SecurityGroup.NotFound"
	OpenapiSharedbandwidthNotFound   = "Openapi.Sharedbandwidth.NotFound"
	EbsOrderInProgress               = "ebs.order.inProgress"
	OpenapiVpcNotFound               = "Openapi.Vpc.NotFound"
	OpenapiSubnetNotFound            = "Openapi.Subnet.NotFound"
	OpenapiPrivateZoneNotFound       = "Openapi.PrivateZone.NotFound"
	OpenapiLoadBalancerNotFound      = "Openapi.LoadBalancer.NotFound"
	OpenapiHealthCheckNotFound       = "Openapi.HealthCheck.NotFound"
	OpenapiTargetGroupNotFound       = "Openapi.TargetGroup.NotFound"
	OpenapiListenerNotFound          = "Openapi.Listener.NotFound"
	EcsAffinityGroupNotBound         = "Ecs.AffinityGroup.NotBound"
	OpenapiRouterTableAccessFailed   = "Openapi.RouterTable.AccessFailed"
	OpenapiAclNotFoundMsg            = "resource not found"
	OpenapiVpceEndpointNotFound      = "Openapi.VpceEndpoint.NotFound"
	OpenapiVpceServiceNotFound       = "Openapi.EndpointService.NotFound"
	OpenapiOssNoSuchBucket           = "oss.resource.noSuchBucket"
	OpenapiSfsNotExists              = "sfs.sfsInfo.resourceNotExists"
	OpenapiPrefixListAccessFailed    = "Openapi.PrefixList.AccessFailed"
	OpenapiRouteTableAccessFailed    = "Openapi.RouteTable.AccessFailed"
	OpenapiVpceServiceAccessFailed   = "Openapi.Endpoint_service.AccessFailed"
	OpenapiDhcpoptionsetsNotFound    = "Openapi.Dhcpoptionsets.NotFound"
	OpenapiEbmNotFound               = "Ebm.Instance.NotFound"
	OpenapiEbsBackupNotFound         = "EbsBackup.BackupInfo.NotFound"
	OpenapiEbsBackupPolicyNotFound   = "EbsBackup.PolicyInfo.NotFound"
	OpenapiEcsBackupPolicyNotFound   = "EcsBackup.Backup.NotFound"
	OpenapiAccessControlNotFound     = "Openapi.AccessControl.NotFound"
	OpenapiCertificateAccessFailed   = "Openapi.Certificate.AccessFailed"
	OpenapiElbPolicyNotFound         = "Openapi.ElbPolicy.NotFound"
	OpenapiTargetNotFound            = "Openapi.Target.NotFound"
	OpenapiDnatNotFound              = "Openapi.Dnat.NotFound"
	OpenapiSnatNotFound              = "Openapi.Snat.NotFound"
	OpenapiParameterError            = "Openapi.Parameter.Error"
	OpenapiPrivateZoneRecordNotFound = "Openapi.PrivateZoneRecord.NotFound"
	OpenapiHavipNotFound             = "Openapi.Havip.NotFound"
	OpenapiVpcPeeringNotFound        = "Openapi.VpcPeering.NotFound"

	OpenapiVpcPortNotFound = "Openapi.Parameter.Error"
	OpenapiCCSENotExist    = "CCE_2024"
	CtiamNoPermission      = "CTIAM_0005"
	CtiamNoPrivilege       = "CTIAM_1044"

	ErrorStatusCode        = 900
	NormalStatusCode       = 800
	NormalStatusCodeString = "800"
)

var InvalidReturnObjError = errors.New("invalid return object")
var InvalidReturnObjResultsError = errors.New("invalid result object results")
var ResourceNotExistError = errors.New("resource not exist")

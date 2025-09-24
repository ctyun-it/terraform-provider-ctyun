package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryCustomerQuotasInRegionV41Api
/* 根据regionID查询用户配额
 */type CtecsQueryCustomerQuotasInRegionV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryCustomerQuotasInRegionV41Api(client *core.CtyunClient) *CtecsQueryCustomerQuotasInRegionV41Api {
	return &CtecsQueryCustomerQuotasInRegionV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/region/customer-quotas",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryCustomerQuotasInRegionV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryCustomerQuotasInRegionV41Request) (*CtecsQueryCustomerQuotasInRegionV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsQueryCustomerQuotasInRegionV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryCustomerQuotasInRegionV41Request struct {
	RegionID string /*  资源池ID  */
}

type CtecsQueryCustomerQuotasInRegionV41Response struct {
	StatusCode  int32                                                 `json:"statusCode,omitempty"`  /*  返回状态码('800为成功，900为失败)  ，默认值:800  */
	ErrorCode   string                                                `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码。为空表示成功。  */
	Message     string                                                `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                                `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsQueryCustomerQuotasInRegionV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsQueryCustomerQuotasInRegionV41ReturnObjResponse struct {
	Quotas       *CtecsQueryCustomerQuotasInRegionV41ReturnObjQuotasResponse       `json:"quotas"`       /*  本资源池配额信息  */
	Global_quota *CtecsQueryCustomerQuotasInRegionV41ReturnObjGlobal_quotaResponse `json:"global_quota"` /*  全局配额信息  */
}

type CtecsQueryCustomerQuotasInRegionV41ReturnObjQuotasResponse struct {
	Network_acl_limit                        int32  `json:"network_acl_limit,omitempty"`                        /*  ACL规则个数  */
	Max_capacity_of_disk_creation_cs         int32  `json:"max_capacity_of_disk_creation_cs,omitempty"`         /*  单块磁盘创建时的最大容量-CS(GB)  */
	Disk_backup_capacity_limit               int32  `json:"disk_backup_capacity_limit,omitempty"`               /*  云硬盘备份容量上限-OS(GB)  */
	Storage_limit                            int32  `json:"storage_limit,omitempty"`                            /*  存储总容量上限(GB)  */
	Network_limit_each_vpc                   int32  `json:"network_limit_each_vpc,omitempty"`                   /*  单个VPC下子网个数上限-CS  */
	Load_balancer_limit_each_ip_os           int32  `json:"load_balancer_limit_each_ip_os,omitempty"`           /*  单个负载均衡下的监听器个数上限-OS  */
	Monitoring_item_limit                    int32  `json:"monitoring_item_limit,omitempty"`                    /*  单个监控视图下的监控项个数上限  */
	Monitor_alerm_rules_limit                int32  `json:"monitor_alerm_rules_limit,omitempty"`                /*  告警规则个数上限  */
	Vm_limit_each_load_balancer_os           int32  `json:"vm_limit_each_load_balancer_os,omitempty"`           /*  单个监听器下可绑定的主机个数上限  */
	Network_limit_each_vpc_os                int32  `json:"network_limit_each_vpc_os,omitempty"`                /*  单个VPC下的子网个数上限-OS  */
	Pm_limit_per_platform                    int32  `json:"pm_limit_per_platform,omitempty"`                    /*  单资源池下物理机个数上限  */
	Snapshot_limit_per_cloud_server_os       int32  `json:"snapshot_limit_per_cloud_server_os,omitempty"`       /*  单台云服务器快照上限-OS  */
	Max_duration_of_elastic_ip_creation      string `json:"max_duration_of_elastic_ip_creation,omitempty"`      /*  创建弹性IP可选的最大时长(年)  */
	Vpc_limit_os                             int32  `json:"vpc_limit_os,omitempty"`                             /*  VPC上限-OS  */
	Memory_limit                             int32  `json:"memory_limit,omitempty"`                             /*  内存上限(GB)  */
	Max_bandwidth_of_elastic_ip_creation     int32  `json:"max_bandwidth_of_elastic_ip_creation,omitempty"`     /*  创建弹性IP时的带宽上限  */
	Network_cards_limit                      int32  `json:"network_cards_limit,omitempty"`                      /*  单个主机网卡个数上限  */
	Private_image_limit                      int32  `json:"private_image_limit,omitempty"`                      /*  私有镜像上限-CS  */
	Snapshot_limit_os                        int32  `json:"snapshot_limit_os,omitempty"`                        /*  快照个数上限-OS  */
	Vm_limit_each_time                       int32  `json:"vm_limit_each_time,omitempty"`                       /*  单次创建云主机个数上限  */
	Vpc_limit                                int32  `json:"vpc_limit,omitempty"`                                /*  VPC上限-CS  */
	Pm_mem_total_limit_per_platform          int32  `json:"pm_mem_total_limit_per_platform,omitempty"`          /*  单资源池物理机内存总额上限  */
	Load_balancer_limit_each_ip              int32  `json:"load_balancer_limit_each_ip,omitempty"`              /*  单个负载均衡下监听器个数上限-CS  */
	Volume_limit_each_time                   int32  `json:"volume_limit_each_time,omitempty"`                   /*  单次创建磁盘个数上限  */
	Load_balancer_limit                      int32  `json:"load_balancer_limit,omitempty"`                      /*  负载均衡个数上限-CS  */
	Disk_backup_amount_limit                 int32  `json:"disk_backup_amount_limit,omitempty"`                 /*  云硬盘备份的数量上限-OS  */
	Max_capacity_of_disk_creation_os         int32  `json:"max_capacity_of_disk_creation_os,omitempty"`         /*  创建单块磁盘时最大容量-OS  */
	Key_pair_limit                           int32  `json:"key_pair_limit,omitempty"`                           /*  密匙对上限  */
	Max_duration_of_host_creation            string `json:"max_duration_of_host_creation,omitempty"`            /*  创建主机时可选的最大时长(年)  */
	Security_group_rules_limit               int32  `json:"security_group_rules_limit,omitempty"`               /*  安全组规则个数上限  */
	Pm_cpu_total_limit_per_platform          int32  `json:"pm_cpu_total_limit_per_platform,omitempty"`          /*  单资源池物理机CPU总配额  */
	Max_duration_of_disk_product_creation    string `json:"max_duration_of_disk_product_creation,omitempty"`    /*  磁盘产品创建时可选最大时长(年)  */
	Max_capacity_of_sys_disk_creation_os     int32  `json:"max_capacity_of_sys_disk_creation_os,omitempty"`     /*  创建系统盘时可选的最大容量-OS(GB)  */
	Snapshot_limit_per_cloud_server          int32  `json:"snapshot_limit_per_cloud_server,omitempty"`          /*  单台云服务器快照个数上限-cs  */
	Network_acl_limit_os                     int32  `json:"network_acl_limit_os,omitempty"`                     /*  ACL规则个数上限-OS  */
	Volume_limit_each_vm                     int32  `json:"volume_limit_each_vm,omitempty"`                     /*  单台云主机可挂载磁盘块数上限  */
	Volume_size_limit                        int32  `json:"volume_size_limit,omitempty"`                        /*  磁盘总容量上限(GB)  */
	Snapshot_limit                           int32  `json:"snapshot_limit,omitempty"`                           /*  快照总个数上限-CS  */
	Public_ip_limit_each_time                int32  `json:"public_ip_limit_each_time,omitempty"`                /*  单次创建公网IP个数上限  */
	Private_image_limit_os                   int32  `json:"private_image_limit_os,omitempty"`                   /*  私有镜像上限-OS  */
	Load_balancer_limit_os                   int32  `json:"load_balancer_limit_os,omitempty"`                   /*  负载均衡个数上限-OS  */
	Volume_size_lower_limit                  int32  `json:"volume_size_lower_limit,omitempty"`                  /*  单块磁盘创建时可选的最小容量(GB)  */
	Monitor_view_limit                       int32  `json:"monitor_view_limit,omitempty"`                       /*  单个监控面板下可添加的监控视图个数上限  */
	Vcpu_limit                               int32  `json:"vcpu_limit,omitempty"`                               /*  VCPU总核数  */
	Self_customized_alerm_model_limit        int32  `json:"self_customized_alerm_model_limit,omitempty"`        /*  自定义告警模板个数上限  */
	Monitor_panel_limit                      int32  `json:"monitor_panel_limit,omitempty"`                      /*  监控面板个数上限  */
	Vm_limit_each_load_balancer              int32  `json:"vm_limit_each_load_balancer,omitempty"`              /*  单个监听器可绑定的主机个数上限-CS  */
	Public_ip_limit                          int32  `json:"public_ip_limit,omitempty"`                          /*  弹性公网IP个数上限  */
	Security_groups_limit                    int32  `json:"security_groups_limit,omitempty"`                    /*  安全组个数上限  */
	Total_volume_limit                       int32  `json:"total_volume_limit,omitempty"`                       /*  磁盘总块数  */
	Backup_policy_limit                      int32  `json:"backup_policy_limit,omitempty"`                      /*  云硬盘备份策略个数上限  */
	Vm_limit                                 int32  `json:"vm_limit,omitempty"`                                 /*  云主机总数上限  */
	Rule_limit_of_direction_out_per_acl_cs   int32  `json:"rule_limit_of_direction_out_per_acl_cs,omitempty"`   /*  单ACL下出方向规则个数上限-CS  */
	Rule_limit_of_direction_out_per_acl_os   int32  `json:"rule_limit_of_direction_out_per_acl_os,omitempty"`   /*  单ACL下出方向规则个数上限-OS  */
	Rule_limit_of_direction_in_per_acl_os    int32  `json:"rule_limit_of_direction_in_per_acl_os,omitempty"`    /*  单ACL下入方向规则个数上限-OS  */
	Rule_limit_of_direction_in_per_acl_cs    int32  `json:"rule_limit_of_direction_in_per_acl_cs,omitempty"`    /*  单ACL下入方向规则个数上限-CS  */
	Public_ip_v6_os_limit                    int32  `json:"public_ip_v6_os_limit,omitempty"`                    /*  ipv6带宽包上限-OS  */
	Csbs_backup_policy_limit                 int32  `json:"csbs_backup_policy_limit,omitempty"`                 /*  云主机备份策略上限  */
	Csbs_backup_policy_instance_limit        int32  `json:"csbs_backup_policy_instance_limit,omitempty"`        /*  云主机备份策略绑定云主机个数上限  */
	Csbs_backup_amount_limit                 int32  `json:"csbs_backup_amount_limit,omitempty"`                 /*  云主机备份上限  */
	Csbs_backup_amount_limit_os              int32  `json:"csbs_backup_amount_limit_os,omitempty"`              /*  OS资源池云主机备份上限  */
	Csbs_backup_capacity_limit               int32  `json:"csbs_backup_capacity_limit,omitempty"`               /*  云主机备份磁盘容量上限(GB)  */
	Csbs_backup_capacity_limit_os            int32  `json:"csbs_backup_capacity_limit_os,omitempty"`            /*  OS资源池云主机备份磁盘容量上限(GB)  */
	Max_count_of_nic_per_vm                  int32  `json:"max_count_of_nic_per_vm,omitempty"`                  /*  单台虚机可添加网卡数量上限  */
	Max_num_of_vm_per_vip                    int32  `json:"max_num_of_vm_per_vip,omitempty"`                    /*  单虚IP可绑定的主机数量上限  */
	Volume_limit_each_vm_os                  int32  `json:"volume_limit_each_vm_os,omitempty"`                  /*  单台云主机可挂载磁盘块数上限-OS  */
	Vm_group_limit                           int32  `json:"vm_group_limit,omitempty"`                           /*  云主机反亲和组个数上限  */
	Vm_limit_per_group                       int32  `json:"vm_limit_per_group,omitempty"`                       /*  单个云主机反亲和组可绑定的主机数量上限  */
	Sdwan_limit                              int32  `json:"sdwan_limit,omitempty"`                              /*  sdwan总数上限  */
	Sdwan_limit_each_edge                    int32  `json:"sdwan_limit_each_edge,omitempty"`                    /*  单个sdwan可包含的翼云edge个数上限  */
	Sdwan_limit_each_site                    int32  `json:"sdwan_limit_each_site,omitempty"`                    /*  单个sdwan可包含的站点个数上限  */
	Edge_limit                               int32  `json:"edge_limit,omitempty"`                               /*  edge个数上限  */
	Site_limit                               int32  `json:"site_limit,omitempty"`                               /*  站点个数上限  */
	Share_bandwidth_count_per_user_limit     int32  `json:"share_bandwidth_count_per_user_limit,omitempty"`     /*  单个用户可以购买的共享带宽数量  */
	Max_duration_of_share_bandwidth_creation string `json:"max_duration_of_share_bandwidth_creation,omitempty"` /*  共享带宽产品创建的最大时长(年)  */
	Max_num_of_share_bandwidth_per_user      int32  `json:"max_num_of_share_bandwidth_per_user,omitempty"`      /*  共享带宽产品创建的带宽最大值  */
	Ip_count_per_share_bandwidth             int32  `json:"ip_count_per_share_bandwidth,omitempty"`             /*  单个共享带宽可添加的公网 IP 最大值  */
	Max_buckets_of_oss                       int32  `json:"max_buckets_of_oss,omitempty"`                       /*  单个资源池下对象存储可创建的存储桶个数  */
	Max_capacity_of_csbs_repo                int32  `json:"max_capacity_of_csbs_repo,omitempty"`                /*  单个云主机备份存储库最大容量(GB)  */
	Min_capacity_of_csbs_repo                int32  `json:"min_capacity_of_csbs_repo,omitempty"`                /*  单个云主机备份存储库最小容量(GB)  */
	Csbs_repo_limit                          int32  `json:"csbs_repo_limit,omitempty"`                          /*  云主机备份存储库个数  */
	Max_duration_of_csbs_repo_creation       string `json:"max_duration_of_csbs_repo_creation,omitempty"`       /*  云主机备份存储库创建的最大时长(年)  */
	Csbs_backup_policy_repository_limit      int32  `json:"csbs_backup_policy_repository_limit,omitempty"`      /*  单个策略可绑定存储库上限  */
	Scaling_group_limit                      int32  `json:"scaling_group_limit,omitempty"`                      /*  弹性伸缩组上限  */
	Scaling_config_limit                     int32  `json:"scaling_config_limit,omitempty"`                     /*  弹性伸缩配置上限  */
	Scaling_rule_limit                       int32  `json:"scaling_rule_limit,omitempty"`                       /*  弹性伸缩策略上限  */
	Max_bandwidth_of_elastic_ip_v6_creation  int32  `json:"max_bandwidth_of_elastic_ip_v6_creation,omitempty"`  /*  创建IPV6时的带宽上限  */
	Site_limit_each_time                     int32  `json:"site_limit_each_time,omitempty"`                     /*  单次创建站点个数上限  */
	Address_limit                            int32  `json:"address_limit,omitempty"`                            /*  收货地址个数上限  */
	Address_limit_each_time                  int32  `json:"address_limit_each_time,omitempty"`                  /*  单次创建收货地址个数上限  */
	Sdwan_acl_limit                          int32  `json:"sdwan_acl_limit,omitempty"`                          /*  SDWAN_ACL个数上限  */
	Sdwan_acl_rule_limit                     int32  `json:"sdwan_acl_rule_limit,omitempty"`                     /*  SDWAN_ACL规则个数上限  */
	Pm_create_num_limit_per_time             int32  `json:"pm_create_num_limit_per_time,omitempty"`             /*  单次物理机创建个数最大值  */
	P_image_share_to_others_quota            int32  `json:"p_image_share_to_others_quota,omitempty"`            /*  私有镜像共享人数上限  */
	Ch_network_instance_limit                int32  `json:"ch_network_instance_limit,omitempty"`                /*  云间高速加载网络实例个数上限  */
	Ch_network_instance_region_limit         int32  `json:"ch_network_instance_region_limit,omitempty"`         /*  云间高速加载网络实例区域个数上限  */
	Ch_limit                                 int32  `json:"ch_limit,omitempty"`                                 /*  云间高速个数上限  */
	SiteTmpl_limit                           int32  `json:"siteTmpl_limit,omitempty"`                           /*  站点模板数量上限  */
	Max_bandwidth_of_elastic_ip_creation_os  int32  `json:"max_bandwidth_of_elastic_ip_creation_os,omitempty"`  /*  创建弹性IP时的带宽上限-OS  */
	Max_num_of_vip_per_vm                    int32  `json:"max_num_of_vip_per_vm,omitempty"`                    /*  单台虚机可绑定的虚IP数量上限  */
	Sdwan_monitor_alarm_rules_limit          int32  `json:"sdwan_monitor_alarm_rules_limit,omitempty"`          /*  SDWAN告警规则个数上限  */
	Max_num_of_vip_per_pm                    int32  `json:"max_num_of_vip_per_pm,omitempty"`                    /*  单台物理机可绑定的虚IP数量上限  */
	Max_num_of_pm_per_vip                    int32  `json:"max_num_of_pm_per_vip,omitempty"`                    /*  单个虚IP可绑定的物理机数量上限  */
	Sfs_fs_count_limit                       int32  `json:"sfs_fs_count_limit,omitempty"`                       /*  弹性文件系统个数上限  */
	Sfs_fs_volume_limit                      int32  `json:"sfs_fs_volume_limit,omitempty"`                      /*  弹性文件系统总容量上限(TB)  */
	Sfs_fs_mount_point_count_limit           int32  `json:"sfs_fs_mount_point_count_limit,omitempty"`           /*  弹性文件系统挂载点个数上限  */
	Sfs_permission_group_count_limit         int32  `json:"sfs_permission_group_count_limit,omitempty"`         /*  弹性文件系统权限组个数上限  */
	Sfs_permission_rule_count_limit          int32  `json:"sfs_permission_rule_count_limit,omitempty"`          /*  弹性文件系统权限组规则个数上限  */
	Elb_cert_limit                           int32  `json:"elb_cert_limit,omitempty"`                           /*  负载均衡证书总个数  */
	Vpc_router_limit_per_table               int32  `json:"vpc_router_limit_per_table,omitempty"`               /*  单个VPC下路由规则个数上限  */
	Bks_repo_limit                           int32  `json:"bks_repo_limit,omitempty"`                           /*  云硬盘备份存储库个数  */
	Max_capacity_of_bks_repo                 int32  `json:"max_capacity_of_bks_repo,omitempty"`                 /*  单个云硬盘备份存储库最大容量  */
	Min_capacity_of_bks_repo                 int32  `json:"min_capacity_of_bks_repo,omitempty"`                 /*  单个硬盘备份存储库最小容量(GB)  */
	Max_duration_of_bks_repo_creation        string `json:"max_duration_of_bks_repo_creation,omitempty"`        /*  云硬盘备份存储库创建的最大时长(年)  */
	Bks_backup_policy_repository_limit       int32  `json:"bks_backup_policy_repository_limit,omitempty"`       /*  单个云硬盘备份策略可绑定存储库上限  */
	Bks_backup_policy_disk_limit             int32  `json:"bks_backup_policy_disk_limit,omitempty"`             /*  云硬盘备份策略绑定云硬盘个数上限  */
	Routing_table_limit                      int32  `json:"routing_table_limit,omitempty"`                      /*  路由表默认配额  */
	Share_ebs_attach_count                   int32  `json:"share_ebs_attach_count,omitempty"`                   /*  共享硬盘可配置数量  */
	P2p_router_count_limit_per_connection    int32  `json:"p2p_router_count_limit_per_connection,omitempty"`    /*  对等连接内路由数量上限  */
	P2p_connection_count_limit               int32  `json:"p2p_connection_count_limit,omitempty"`               /*  对等连接数量上限  */
	P2p_router_count_limit_per_batch         int32  `json:"p2p_router_count_limit_per_batch,omitempty"`         /*  对等连接单次创建路由数量上限  */
	Ch_order_bandwidth_limit                 int32  `json:"ch_order_bandwidth_limit,omitempty"`                 /*  云间高速购买带宽包带宽值上限  */
	Ch_order_bandwidth_num_limit             int32  `json:"ch_order_bandwidth_num_limit,omitempty"`             /*  云间高速购买带宽包个数上限  */
	Oss_bucket_count_limit                   int32  `json:"oss_bucket_count_limit,omitempty"`                   /*  对象存储默认配额  */
	Vpn_user_gate_count_limit                int32  `json:"vpn_user_gate_count_limit,omitempty"`                /*  VPN用户网关个数上限  */
	Vpn_connection_count_limit               int32  `json:"vpn_connection_count_limit,omitempty"`               /*  VPN连接个数上限  */
	Vpn_gate_count_limit                     int32  `json:"vpn_gate_count_limit,omitempty"`                     /*  VPN网关个数上限  */
	Route_limit_per_table                    int32  `json:"route_limit_per_table,omitempty"`                    /*  路由规则  */
	Vpce_limit_per_vpc                       int32  `json:"vpce_limit_per_vpc,omitempty"`                       /*  单个VPC下终端节点个数上限  */
	Vpce_server_limit_per_vpc                int32  `json:"vpce_service_limit_per_vpc,omitempty"`               /*  单个VPC下终端服务节点个数上限  */
	Total_traffic_mirror_limit               int32  `json:"total_traffic_mirror_limit,omitempty"`               /*  流量镜像产品筛选条件配额  */
	Total_traffic_session_limit              int32  `json:"total_traffic_session_limit,omitempty"`              /*  流量镜像产品镜像会话配额  */
	Volume_limit_each_vm_ElasticPM           int32  `json:"volume_limit_each_vm_ElasticPM,omitempty"`           /*  裸金属单块磁盘创建时可选的最小容量(GB)  */
	Max_capacity_of_disk_creation_ElasticPM  int32  `json:"max_capacity_of_disk_creation_ElasticPM,omitempty"`  /*  单块磁盘创建时的最大容量-裸金属(GB)  */
	Cnssl_site_limit                         int32  `json:"cnssl_site_limit,omitempty"`                         /*  云网超级专线站点数量  */
	Total_intranet_dns_limit                 int32  `json:"total_intranet_dns_limit,omitempty"`                 /*  DNS域名配额  */
	Max_count_of_nic_per_pm                  int32  `json:"max_count_of_nic_per_pm,omitempty"`                  /*  单台物理机可添加网卡数量上限  */
	Cnssl_physicsLine_route_limit            int32  `json:"cnssl_physicsLine_route_limit,omitempty"`            /*  SD-WANoe0 0i  */
	Snapshot_policy_limit                    int32  `json:"snapshot_policy_limit,omitempty"`                    /*  云主机快照策略上限  */
	Snapshot_policy_instance_limit           int32  `json:"snapshot_policy_instance_limit,omitempty"`           /*  云主机快照策略绑定云主机上限  */
	Cnssl_physicsLine_snat_limit             int32  `json:"cnssl_physicsLine_snat_limit,omitempty"`             /*  SD-WAN（尊享版）-物理专线SNAT数量  */
	Cnssl_physicsLine_dnat_limit             int32  `json:"cnssl_physicsLine_dnat_limit,omitempty"`             /*  SD-WAN（尊享版）-物理专线DNAT数量  */
	Cnssl_physicsLine_vpc_limit              int32  `json:"cnssl_physicsLine_vpc_limit,omitempty"`              /*  SD-WAN（尊享版）-物理专线入云数量  */
	Cnssl_route_ip_limit                     int32  `json:"cnssl_route_ip_limit,omitempty"`                     /*  SD-WAN（尊享版）-客户侧路由ipv4个数限制  */
	Cnssl_edge_route_limit                   int32  `json:"cnssl_edge_route_limit,omitempty"`                   /*  SD-WAN（尊享版）-智能网关-路由数量  */
	Cnssl_edge_vpc_limit                     int32  `json:"cnssl_edge_vpc_limit,omitempty"`                     /*  SD-WAN（尊享版）-智能网关-入云限制数量  */
	Cnssl_edge_subnet_limit                  int32  `json:"cnssl_edge_subnet_limit,omitempty"`                  /*  SD-WAN（尊享版）-智能网关-子网IP限制数量  */
	Cnssl_physicsLine_app_vpc_limit          int32  `json:"cnssl_physicsLine_app_vpc_limit,omitempty"`          /*  SD-WAN（尊享版）-物理专线应用保障添加VPC数量  */
	Load_balancer_policy_limit_per_listener  int32  `json:"load_balancer_policy_limit_per_listener,omitempty"`  /*  单个监听器下创建的负载均衡转发策略上限  */
	Edge_limit_each_pnet                     int32  `json:"edge_limit_each_pnet,omitempty"`                     /*  单个edge下可配置子网数量  */
	Sdwan_qos_rule_limt                      int32  `json:"sdwan_qos_rule_limt,omitempty"`                      /*  sdwan下Qos规则数量  */
	Sdwan_qos_rule_group_limt                int32  `json:"sdwan_qos_rule_group_limt,omitempty"`                /*  sdwan下Qos规则下五元组数量  */
	Sdwan_qos_limit                          int32  `json:"sdwan_qos_limit,omitempty"`                          /*  sdwan下的qos数量  */
	Sdwan_edge_mpls_ip_limit                 int32  `json:"sdwan_edge_mpls_ip_limit,omitempty"`                 /*  sdwan下edge的过载保护目标检测ip数量上限  */
	Sfs_single_fs_volume_limit               int32  `json:"sfs_single_fs_volume_limit,omitempty"`               /*  单个弹性文件系统容量上限(TB)  */
	Sfs_single_exclusive_fs_volume_limit     int32  `json:"sfs_single_exclusive_fs_volume_limit,omitempty"`     /*  单个专属型文件系统容量上限  */
	Max_duration_of_host_new_creation        string `json:"max_duration_of_host_new_creation,omitempty"`        /*  创建非GPU主机时可选的最大时长(年)  */
	Max_duration_of_network_creation         string `json:"max_duration_of_network_creation,omitempty"`         /*  创建VPN和文件系统时可选的最大时长(年)  */
	Sdwan_edge_static_router_limit           int32  `json:"sdwan_edge_static_router_limit,omitempty"`           /*  sdwan下单个edge里可创建的静态路由数量  */
	Dr_client_limit                          int32  `json:"dr_client_limit,omitempty"`                          /*  单用户单个资源池客户端数  */
	Ch_create_limit                          int32  `json:"ch_create_limit,omitempty"`                          /*  创建云间高速默认配额  */
	Ch_create_net_manage_limit               int32  `json:"ch_create_net_manage_limit,omitempty"`               /*  创建云网关默认配额  */
	Ch_netmanagement_vpc_limit               int32  `json:"ch_netmanagement_vpc_limit,omitempty"`               /*  vpc网络实例默认配额  */
	Ch_netmanagement_cda_limit               int32  `json:"ch_netmanagement_cda_limit,omitempty"`               /*  cda网络实例默认配额  */
	Ch_netmanagement_accountvpc_limit        int32  `json:"ch_netmanagement_accountvpc_limit,omitempty"`        /*  授权vpc网络实例默认配额  */
	Ch_reconsitution_accredit_limit          int32  `json:"ch_reconsitution_accredit_limit,omitempty"`          /*  云间高速（标准版）跨账号授权配额上限  */
	Ch_create_route_limit                    int32  `json:"ch_create_route_limit,omitempty"`                    /*  云间高速（标准版）路由管理创建自定义路由表配额上限  */
	Ch_cda_subnet_limit                      int32  `json:"ch_cda_subnet_limit,omitempty"`                      /*  云间高速（标准版）cda子网选择上限配额  */
	Ch_create_route_num_limit                int32  `json:"ch_create_route_num_limit,omitempty"`                /*  云间高速（标准版）路由条目配额  */
	Ch_vpc_subnet_limit                      int32  `json:"ch_vpc_subnet_limit,omitempty"`                      /*  云间高速（标准版）vpc子网选择上限配额  */
	Ch_vpc_instance_bind_limit               int32  `json:"ch_vpc_instance_bind_limit,omitempty"`               /*  云间高速（标准版）单个vpc被相同云间高速加载次数  */
	Ch_order_bandwidth_num_limit_v2          int32  `json:"ch_order_bandwidth_num_limit_v2,omitempty"`          /*  云间高速购买带宽包个数上限2.0  */
	Ch_order_bandwidth_limit_v2              int32  `json:"ch_order_bandwidth_limit_v2,omitempty"`              /*  云间高速（标准版）2.0订单带宽值上限  */
	Elb_cidr_policy_limit                    int32  `json:"elb_cidr_policy_limit,omitempty"`                    /*  负载均衡访问策略组配额  */
	Elb_cidr_ip_count_limit                  int32  `json:"elb_cidr_ip_count_limit,omitempty"`                  /*  访问策略组-每个IP地址组中IP地址数量  */
	Nic_relate_security_group_limit          int32  `json:"nic_relate_security_group_limit,omitempty"`          /*  网卡可绑定的安全组数量上限  */
	Ssl_vpn_server_limit                     int32  `json:"ssl_vpn_server_limit,omitempty"`                     /*  ssl服务端默认配额  */
	Ssl_vpn_client_limit                     int32  `json:"ssl_vpn_client_limit,omitempty"`                     /*  ssl客户端默认配额  */
	Snap_volume_limit                        int32  `json:"snap_volume_limit,omitempty"`                        /*  快照创建云硬盘个数  */
	Ssl_vpn_gate_count_limit                 int32  `json:"ssl_vpn_gate_count_limit,omitempty"`                 /*  sslvpn网关个数上限  */
	Sfs_oceanfs_volume_limit                 int32  `json:"sfs_oceanfs_volume_limit,omitempty"`                 /*  海量文件系统总容量上限(TB)  */
	Sfs_oceanfs_count_limit                  int32  `json:"sfs_oceanfs_count_limit,omitempty"`                  /*  海量文件系统个数上限  */
	Sfs_hpfs_volume_limit                    int32  `json:"sfs_hpfs_volume_limit,omitempty"`                    /*  并行文件系统总容量上限(TB)  */
	Sfs_hpfs_count_limit                     int32  `json:"sfs_hpfs_count_limit,omitempty"`                     /*  并行文件系统个数上限  */
	Cbr_ecs_limit                            int32  `json:"cbr_ecs_limit,omitempty"`                            /*  云备份客户端配额  */
	Cbr_vault_limit                          int32  `json:"cbr_vault_limit,omitempty"`                          /*  云备份存储库配额  */
	Vip_limit                                int32  `json:"vip_limit,omitempty"`                                /*  单用户单资源池可创建虚拟IP个数  */
	Vpc_create_vip_limit                     int32  `json:"vpc_create_vip_limit,omitempty"`                     /*  单VPC支持创建的VIP数量  */
	Public_ip_cn2_limit                      int32  `json:"public_ip_cn2_limit,omitempty"`                      /*  cn2列表  */
	Rules_limit_of_per_security_group        int32  `json:"rules_limit_of_per_security_group,omitempty"`        /*  单安全组的规则个数上限(不分出入)  */
	Public_ip_v6_limit                       int32  `json:"public_ip_v6_limit,omitempty"`                       /*  4.0创建ipv6带宽的个数上限  */
}

type CtecsQueryCustomerQuotasInRegionV41ReturnObjGlobal_quotaResponse struct {
	Global_public_ip_limit int32 `json:"global_public_ip_limit,omitempty"` /*  弹性公网IP个数上限  */
}

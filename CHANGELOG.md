# Changelog

## v2.0.0 - December 18, 2025

### New Resources

* ctyun_oceanfs
* ctyun_oceanfs_permission_group
* ctyun_oceanfs_permission_rule
* ctyun_oceanfs_permission_group_association
* ctyun_sdwan
* ctyun_sdwan_acl
* ctyun_rabbitmq_vhost
* ctyun_rabbitmq_exchange
* ctyun_rabbitmq_queue
* ctyun_kafka_consumer_group
* ctyun_kafka_topic
* ctyun_kafka_user
* ctyun_kafka_acl
* ctyun_redis_account
* ctyun_redis_backup
* ctyun_redis_instance_whitelist
* ctyun_redis_migration_task
* ctyun_redis_param_template
* ctyun_image_from_ecs
* ctyun_port
* ctyun_ecs_port_association
* ctyun_private_nat
* ctyun_private_nat_dnat
* ctyun_private_nat_snat
* ctyun_private_nat_transit_ip
* ctyun_mongodb_account
* ctyun_mongodb_backup
* ctyun_mongodb_param_template
* ctyun_mongodb_white_list
* ctyun_vip
* ctyun_vip_association
* ctyun_mysql_account
* ctyun_mysql_audit
* ctyun_mysql_backup
* ctyun_mysql_backup_cancel
* ctyun_mysql_backup_recovery
* ctyun_mysql_backup_setting
* ctyun_mysql_copy_parameter_template
* ctyun_mysql_param_template
* ctyun_mysql_readonly_instance
* ctyun_mysql_white_list
* ctyun_postgresql_account
* ctyun_postgresql_backup
* ctyun_postgresql_param_template
* ctyun_postgresql_readonly_instance
* ctyun_postgresql_white_list
* ctyun_mysql_database
* ctyun_ccse_namespace
* ctyun_crs_vpc_attach
* ctyun_acl
* ctyun_acl_rule
* ctyun_prefix_list
* ctyun_subnet_association_acl
* ctyun_express_connect
* ctyun_ec_cloud_gateway
* ctyun_ec_packet
* ctyun_ec_route
* ctyun_ec_region_peer
* ctyun_ec_vpc_instance
* ctyun_iam_user_ak
* ctyun_ecs_data_volume
* ctyun_dhcpoptionset
* ctyun_dhcpoptionset_association_vpc
* ctyun_ec_sdwan_instance
* ctyun_mongodb_restart_db
* ctyun_postgresql_database
* ctyun_mysql_rds_parameter_template
* ctyun_scaling_config
* ctyun_vpc_peer_connection
* ctyun_vpc_peer_connection_attach
* ctyun_private_zone
* ctyun_private_zone_record

### New Data Sources

* ctyun_oceanfs_instances
* ctyun_mysql_character_set
* ctyun_mysql_recoverable_time_points
* ctyun_net_tagss
* ctyun_postgresql_param_templates
* ctyun_prefix_lists
* ctyun_private_nat_cidrs
* ctyun_sdwans
* ctyun_vpc_peer_connections
* ctyun_sdwan_acl_rules
* ctyun_sdwan_acls
* ctyun_rabbitmq_vhosts
* ctyun_rabbitmq_exchanges
* ctyun_rabbitmq_queues
* ctyun_kafka_consumer_groups
* ctyun_kafka_topics
* ctyun_kafka_users
* ctyun_kafka_acls
* ctyun_redis_accounts
* ctyun_redis_backups
* ctyun_redis_instance_whitelists
* ctyun_redis_migration_tasks
* ctyun_redis_param_templates
* ctyun_ports
* ctyun_private_nats
* ctyun_private_nat_dnats
* ctyun_private_nat_snats
* ctyun_private_nat_transit_ips
* ctyun_net_resources_by_tag
* ctyun_mongodb_accounts
* ctyun_mongodb_backups
* ctyun_mongodb_param_templates
* ctyun_mongodb_white_lists
* ctyun_vips
* ctyun_dhcpoptionset_association_vpcs
* ctyun_dhcpoptionsets
* ctyun_ec_sdwan_instances
* ctyun_mysql_accounts
* ctyun_mysql_backups
* ctyun_mysql_databases
* ctyun_mysql_recoverable_time
* ctyun_mysql_white_lists
* ctyun_mysql_param_templates
* ctyun_mysql_parameters
* ctyun_postgresql_accounts
* ctyun_postgresql_backups
* ctyun_postgresql_character_set
* ctyun_postgresql_collation_time_zone
* ctyun_postgresql_databases
* ctyun_postgresql_white_lists
* ctyun_ccse_images
* ctyun_ccse_namespaces
* ctyun_crs_opensource_images
* ctyun_acls
* ctyun_acl_rules
* ctyun_express_connects
* ctyun_ec_region_peers
* ctyun_ec_routes
* ctyun_ec_vpc_instances
* ctyun_ec_cloud_gateways
* ctyun_iam_user_aks
* ctyun_iam_users
* ctyun_iam_policies
* ctyun_private-zones
* ctyun_private_zone_records

### Enhancements

* unify the naming of time fields as create_time, update_time, and expire_time, with values complying with RFC3339
* datasource\ctyun_ebs_backup_policies: fields remove redundant prefixes
* datasource\ctyun_ebs_snapshot_policies: fields remove redundant prefixes
* resource\ctyun_ebs_backup_policy: block to attribute
* resource\ctyun_ebs_snapshot_policy: block to attribute
* datasource\ctyun_ecs_backup_policies: fields remove redundant prefixes
* resource\ctyun_ecs: flavor_name support update
* resource\ctyun_hpfs: fields remove redundant prefixes
* datasource\ctyun_elb_rules: fields remove redundant prefixes
* datasource\ctyun_elb_targets: fields remove redundant prefixes
* resource\ctyun_sfs: fields remove redundant prefixes
* datasource\ctyun_ecs_affinity_groups: fields remove redundant prefixes
* datasource\ctyun_ecs_instances: fields remove redundant prefixes
* datasource\ctyun_hpfs_clusters: fields remove redundant prefixes
* datasource\ctyun_hpfs_instances: fields remove redundant prefixes
* datasource\ctyun_mongodb_instances: fields remove redundant prefixes
* datasource\ctyun_mysql_instances: fields remove redundant prefixes
* datasource\ctyun_postgresql_instances: fields remove redundant prefixes
* datasource\ctyun_mysql_white_lists: fields remove redundant prefixes
* datasource\ctyun_scaling_configs: fields remove redundant prefixes
* datasource\ctyun_scaling_policies: fields remove redundant prefixes
* datasource\ctyun_sfs_instances: fields remove redundant prefixes
* datasource\ctyun_sfs_permission_rules: fields remove redundant prefixes
* resource\ctyun_ccse_cluster: port range validator change and support update
* resource\ctyun_ebm: support bandwidth and remove eip_id
* resource\ctyun_ecs: support bandwidth
* resource\ctyun_ebs: add support fields
* resource\ctyun_elb_acl: source_ips max length change to 50
* resource\ctyun_redis_association_eip: use eip_id replace eip_address
* resource\ctyun_redis_instance: support update
* resource\ctyun_ccse_cluster support cluster update\destroy
* resource\ctyun_keypair support keypair create
* resource\ctyun_kafka_instance support instance restart and billing method modification
* resource\ctyun_redis_instance support instance set ssl config\backup policy\param template id and billing method modification
* resource\ctyun_ebm support metadata
* resource\ctyun_eip_association support ebm

### Bug Fixes
* resource\ctyun_ebs: remove the state record for non-existent remote resources during terraform refresh.
* resource\ctyun_ccse_cluster: fix cycle_count validator.

### Deprecations
* datasource\ctyun_mongodb_association_eips
* datasource\ctyun_mysql_association_eips



## 1.2.0 - September 25, 2025

### New Resources

* ctyun_ebm
* ctyun_ebm_interface
* ctyun_ebm_association_ebs
* ctyun_ecs_affinity_group
* ctyun_ecs_affinity_group_association
* ctyun_vpc_route_table
* ctyun_vpc_route_table_rule
* ctyun_vpce
* ctyun_vpce_service
* ctyun_vpce_service_transit_ip
* ctyun_vpce_service_reverse_rule
* ctyun_vpce_service_connection
* ctyun_zos_bucket
* ctyun_zos_bucket_object
* ctyun_ccse_cluster
* ctyun_ccse_node_pool
* ctyun_redis_instance
* ctyun_redis_association_eip
* ctyun_nat_snat
* ctyun_nat_dnat
* ctyun_nat
* ctyun_elb_loadbalancer
* ctyun_elb_listener
* ctyun_elb_rule
* ctyun_elb_target
* ctyun_elb_target_group
* ctyun_elb_health_check
* ctyun_elb_certificate
* ctyun_elb_acl
* ctyun_mysql_association_eip
* ctyun_mysql_instance
* ctyun_kafka_instance
* ctyun_rabbitmq_instance
* ctyun_postgresql_association_eip
* ctyun_postgresql_instance
* ctyun_ccse_plugin
* ctyun_mongodb_association_eip
* ctyun_mongodb_instance
* ctyun_mysql_white_list
* ctyun_ccse_node_association
* ctyun_ecs_backup
* ctyun_ecs_backup_policy
* ctyun_ecs_backup_policy_bind_instances
* ctyun_ecs_backup_policy_bind_repo
* ctyun_ecs_snapshot
* ctyun_ecs_snapshot_policy
* ctyun_ecs_snapshot_policy_association
* ctyun_hpfs
* ctyun_sfs
* ctyun_sfs_permission_group
* ctyun_sfs_permission_rule
* ctyun_sfs_permission_group_association
* ctyun_scaling
* ctyun_scaling_ecs_protection
* ctyun_scaling_policy
* ctyun_scaling_config
* ctyun_ebs_backup
* ctyun_ebs_backup_policy
* ctyun_ebs_backup_policy_bind_disks
* ctyun_ebs_backup_policy_bind_repo
* ctyun_ebs_snapshot
* ctyun_ebs_snapshot_policy
* ctyun_ebs_snapshot_policy_association

### New Data Sources

* ctyun_ebm_device_types
* ctyun_ebms
* ctyun_ebm_device_images
* ctyun_ebm_device_raids
* ctyun_ebs_volumes
* ctyun_ecs_instances
* ctyun_ecs_affinity_groups
* ctyun_vpcs
* ctyun_subnets
* ctyun_security_groups
* ctyun_vpc_route_tables
* ctyun_vpc_route_table_rules
* ctyun_eips
* ctyun_bandwidths
* ctyun_vpces
* ctyun_vpce_services
* ctyun_vpce_service_transit_ips
* ctyun_vpce_service_reverse_rules
* ctyun_zos_buckets
* ctyun_zos_bucket_objects
* ctyun_ccse_clusters
* ctyun_ccse_node_pools
* ctyun_redis_instances
* ctyun_redis_specs
* ctyun_nats
* ctyun_nat_snats
* ctyun_nat_dnats
* ctyun_elb_loadbalancers
* ctyun_elb_listeners
* ctyun_elb_rules
* ctyun_elb_targets
* ctyun_elb_target_groups
* ctyun_elb_health_checks
* ctyun_elb_certificates
* ctyun_elb_acls
* ctyun_mysql_instances
* ctyun_mysql_specs
* ctyun_kafka_specs
* ctyun_kafka_instances
* ctyun_rabbitmq_specs
* ctyun_rabbitmq_instances
* ctyun_postgresql_instances
* ctyun_postgresql_specs
* ctyun_zones
* ctyun_ccse_plugin_market
* ctyun_mongodb_instances
* ctyun_mongodb_association_eips
* ctyun_mongodb_specs
* ctyun_mysql_white_lists
* ctyun_mysql_association_eips
* ctyun_ccse_template_market
* ctyun_ecs_snapshots
* ctyun_ecs_snapshot_policies
* ctyun_ecs_backups
* ctyun_ecs_backup_policies
* ctyun_hpfs_instances
* ctyun_hpfs_clusters
* ctyun_sfs_permission_rules
* ctyun_sfs_instances
* ctyun_scaling_configs
* ctyun_scaling_policies
* ctyun_scalings
* ctyun_scaling_activities
* ctyun_scaling_ecs_list
* ctyun_ebs_snapshots
* ctyun_ebs_snapshot_policies
* ctyun_ebs_backups
* ctyun_ebs_backup_policies

### Enhancements


### Bug Fixes

* resource/ctyun_ecs fix config
* resource/ctyun_bandwidth fix create

### Deprecations


package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunMysqlInstances{}
	_ datasource.DataSourceWithConfigure = &ctyunMysqlInstances{}
)

type ctyunMysqlInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunMysqlInstances() datasource.DataSource {
	return &ctyunMysqlInstances{}
}
func (c *ctyunMysqlInstances) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunMysqlInstances) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mysql_instances"

}

func (c *ctyunMysqlInstances) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10033813/10134365**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "项目ID",
			},
			"page_now": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "当前页",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "单页记录条数",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"tag_vo_list": schema.ListNestedAttribute{
				Optional:    true,
				Description: "标签列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "标签value",
						},
						"key": schema.StringAttribute{
							Optional:    true,
							Description: "标签key",
						},
						"label_id": schema.StringAttribute{
							Optional:    true,
							Description: "k-v的唯一标识,位于IT那边",
						},
					},
				},
			},
			"res_db_engine": schema.StringAttribute{
				Optional:    true,
				Description: "数据库引擎 枚举5.7, 8.0",
			},
			"prod_inst_name": schema.StringAttribute{
				Optional:    true,
				Description: "实例名称",
			},
			"vip": schema.StringAttribute{
				Optional:    true,
				Description: "连接ip",
			},
			"mysql_instances": schema.ListNestedAttribute{
				Computed:    true,
				Description: "mysql实例列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"prod_inst_name": schema.StringAttribute{
							Computed:    true,
							Description: "实例名称",
						},
						"outer_prod_inst_id": schema.StringAttribute{
							Computed:    true,
							Description: "实例ID",
						},
						"prod_bill_type": schema.Int32Attribute{
							Computed:    true,
							Description: "计费模式 0:按月计费,1:按天计费,2:按年计费,3:按流量计费 4按需计费",
							Validators: []validator.Int32{
								int32validator.Between(0, 4),
							},
						},
						"prod_type": schema.Int32Attribute{
							Computed:    true,
							Description: "0:单机,1:一主一从,2:一主两从,4:只读实例",
							Validators: []validator.Int32{
								int32validator.Between(0, 4),
							},
						},
						"prod_running_status": schema.Int32Attribute{
							Computed:    true,
							Description: "0.正常 1.重启中 2.备份中 3.恢复中 4.修改参数中 5.应用参数组中 6.扩容预处理中 7.扩容预处理完成 8.修改端口中 9.迁移中 10.重置密码中 11.修改数据复制方式中 12.缩容预处理中 13.缩容预处理完成 15.内核小版本升级 17.迁移可用区中 18.修改备份配置中 20.停止中 21.已停止 22.启动中 26.白名单配置中\t",
							Validators: []validator.Int32{
								int32validator.Between(0, 26),
							},
						},
						"prod_order_status": schema.Int32Attribute{
							Computed:    true,
							Description: "0.正常 1.欠费暂停 2.已注销 3.创建中 4.施工失败 5.到期退订状态 6.新增的状态-openApi暂停 7.创建完成等待变更单 8.待注销 9.手动暂停 10.手动退订",
							Validators: []validator.Int32{
								int32validator.Between(0, 10),
							},
						},
						"alive": schema.Int32Attribute{
							Computed:    true,
							Description: "运行状态 0正常 -1异常",
							Validators: []validator.Int32{
								int32validator.Between(0, 1),
							},
						},
						"vip": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟IP地址",
						},
						"vip6": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟IPv6地址",
						},
						"write_port": schema.StringAttribute{
							Computed:    true,
							Description: "写数据端口",
						},
						"read_port": schema.StringAttribute{
							Computed:    true,
							Description: "读端口",
						},
						"create_time": schema.Int64Attribute{
							Computed:    true,
							Description: "创建时间",
						},
						"expire_time": schema.Int64Attribute{
							Computed:    true,
							Description: "到期时间",
						},
						"machine_spec": schema.StringAttribute{
							Computed:    true,
							Description: "实例规格",
						},
						"prod_db_engine": schema.StringAttribute{
							Computed:    true,
							Description: "数据库引擎",
						},
						"disk_size": schema.Int32Attribute{
							Computed:    true,
							Description: "存储空间大小 单位G",
						},
						"new_mysql_version": schema.StringAttribute{
							Computed:    true,
							Description: "mysql版本",
						},
						"disk_type": schema.StringAttribute{
							Computed:    true,
							Description: "存储类型：1.SATA 2.SAS 3.SSD 4.LVM",
						},
						"backup_disk_used_rated": schema.Int32Attribute{
							Computed:    true,
							Description: "备份空间使用率",
						},
						"net_name": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云",
						},
						"subnet": schema.StringAttribute{
							Computed:    true,
							Description: "子网名",
						},
						"prod_inst_flag": schema.StringAttribute{
							Computed:    true,
							Description: "实例标签",
						},
						"order_id": schema.Int64Attribute{
							Computed:    true,
							Description: "订单lD",
						},
						"can_operate": schema.Int32Attribute{
							Computed:    true,
							Description: "是否可升级",
							Validators: []validator.Int32{
								int32validator.Between(0, 1),
							},
						},
						"tpl_name": schema.StringAttribute{
							Computed:    true,
							Description: "实例模板名",
						},
						"audit_log_status": schema.Int32Attribute{
							Computed:    true,
							Description: "日志审计开关",
							Validators: []validator.Int32{
								int32validator.Between(0, 1),
							},
						},
						"parameter_group_used": schema.StringAttribute{
							Computed:    true,
							Description: "参数模板名",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云ID",
						},
						"use_zos": schema.Int32Attribute{
							Computed:    true,
							Description: "是否使用对象存储",
							Validators: []validator.Int32{
								int32validator.Between(0, 1),
							},
						},
						"prod_inst_id": schema.Int64Attribute{
							Computed:    true,
							Description: "内部实例id",
						},
						"db_mysql_version": schema.StringAttribute{
							Computed:    true,
							Description: "mysql版本",
						},
						"resources": schema.StringAttribute{
							Computed:    true,
							Description: "订单来源",
						},
						"user_id": schema.Int64Attribute{
							Computed:    true,
							Description: "用户id",
						},
						"prod_bill_time": schema.Int32Attribute{
							Computed:    true,
							Description: "计费模式 类型",
						},
						"tenant_id": schema.StringAttribute{
							Computed:    true,
							Description: "租户id",
						},
						"tpl_code": schema.StringAttribute{
							Computed:    true,
							Description: "实例模板code",
						},
						"prod_id": schema.Int64Attribute{
							Computed:    true,
							Description: "产品id",
						},
						"prod_inst_set_name": schema.StringAttribute{
							Computed:    true,
							Description: "实例set名称",
						},
						"security_group": schema.StringAttribute{
							Computed:    true,
							Description: "安全组名",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "项目id（可通过接口/teledb-dcp/v2/openapi/dcp-order-info/project/list查询所有项目id）",
						},
						"project_name": schema.StringAttribute{
							Computed:    true,
							Description: "项目名",
						},
						"pause_enable": schema.BoolAttribute{
							Computed:    true,
							Description: "是否允许暂停",
						},
						"inst_release_protection_status": schema.Int32Attribute{
							Computed:    true,
							Description: "实例释放保护开关 1:on,0:off",
							Validators: []validator.Int32{
								int32validator.Between(0, 1),
							},
						},
						"last_manual_back_up": schema.Int64Attribute{
							Computed:    true,
							Description: "最后系统全备id",
						},
						"renewal_enable": schema.BoolAttribute{
							Computed:    true,
							Description: "是否允许续订",
						},
						"is_mgr": schema.Int32Attribute{
							Computed:    true,
							Description: "1：mgr实例 0：非mgr实例",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunMysqlInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunMysqlInstancesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	if config.PageNow.ValueInt32() == 0 {
		config.PageNow = types.Int32Value(1)
	}
	if config.PageSize.ValueInt32() == 0 {
		config.PageSize = types.Int32Value(100)
	}
	params := &mysql.TeledbGetListRequest{
		PageNow:  config.PageNow.ValueInt32(),
		PageSize: config.PageSize.ValueInt32(),
	}
	if !config.TagVOList.IsNull() {
		var tagVoList []TageVoModel
		var tagList []mysql.TagVO
		diag := config.TagVOList.ElementsAs(ctx, &tagVoList, true)
		if diag.HasError() {
			return
		}
		for _, tagItem := range tagVoList {
			var tag mysql.TagVO
			tag.Key = tagItem.Key.ValueStringPointer()
			tag.Value = tagItem.Value.ValueStringPointer()
			tag.LabelId = tagItem.LabelID.ValueStringPointer()
			tagList = append(tagList, tag)
		}
		params.TagVOList = tagList
	}
	if config.ResDbEngine.ValueString() != "" {
		params.ResDbEngine = config.ResDbEngine.ValueStringPointer()
	}
	if config.ProdInstName.ValueString() != "" {
		params.ProdInstName = config.ProdInstName.ValueStringPointer()
	}
	if config.Vip.ValueString() != "" {
		params.Vip = config.Vip.ValueStringPointer()
	}
	header := &mysql.TeledbGetListHeaders{
		RegionID: regionId,
	}
	if config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkCtMysqlApis.TeledbGetListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return
	} else if resp.StatusCode != 0 {
		err = fmt.Errorf("API return error. Message: %s ", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回体
	var mysqlInstances []MysqlInstanceModel
	returnObjMysqlInstances := resp.ReturnObj.List
	for _, instance := range returnObjMysqlInstances {
		var mysqlInstance MysqlInstanceModel
		mysqlInstance.ProdInstName = types.StringValue(instance.ProdInstName)
		mysqlInstance.OuterProdInstId = types.StringValue(instance.OuterProdInstId)
		mysqlInstance.ProdBillType = types.Int32Value(instance.ProdBillType)
		mysqlInstance.ProdType = types.Int32Value(instance.ProdType)
		mysqlInstance.ProdRunningStatus = types.Int32Value(instance.ProdRunningStatus)
		mysqlInstance.ProdOrderStatus = types.Int32Value(instance.ProdOrderStatus)
		mysqlInstance.Alive = types.Int32Value(instance.Alive)
		mysqlInstance.Vip = types.StringValue(instance.Vip)
		mysqlInstance.Vip6 = types.StringValue(instance.Vip6)
		mysqlInstance.WritePort = types.StringValue(instance.WritePort)
		mysqlInstance.ReadPort = types.StringValue(instance.ReadPort)
		mysqlInstance.CreateTime = types.Int64Value(instance.CreateTime)
		mysqlInstance.ExpireTime = types.Int64Value(instance.ExpireTime)
		mysqlInstance.MachineSpec = types.StringValue(instance.MachineSpec)
		mysqlInstance.ProdDbEngine = types.StringValue(instance.ProdDbEngine)
		mysqlInstance.DiskSize = types.Int32Value(instance.DiskSize)
		mysqlInstance.NewMysqlVersion = types.StringValue(instance.NewMysqlVersion)
		mysqlInstance.DiskType = types.StringValue(instance.DiskType)
		mysqlInstance.BackupDiskUsedRated = types.Int32Value(instance.BackupDiskUsedRated)
		mysqlInstance.NetName = types.StringValue(instance.NetName)
		mysqlInstance.Subnet = types.StringValue(instance.Subnet)
		mysqlInstance.ProdInstFlag = types.StringValue(instance.ProdInstFlag)
		mysqlInstance.OrderId = types.Int64Value(instance.OrderId)
		mysqlInstance.CanOperate = types.Int32Value(instance.CanOperate)
		mysqlInstance.TplName = types.StringValue(instance.TplName)
		mysqlInstance.AuditLogStatus = types.Int32Value(instance.AuditLogStatus)
		mysqlInstance.ParameterGroupUsed = types.StringValue(instance.ParameterGroupUsed)
		mysqlInstance.VpcID = types.StringValue(instance.VpcId)
		mysqlInstance.Usezos = types.Int32Value(instance.Usezos)
		mysqlInstance.ProdInstId = types.Int64Value(instance.ProdInstId)
		mysqlInstance.DbMysqlVersion = types.StringValue(instance.DbMysqlVersion)
		mysqlInstance.Resources = types.StringValue(instance.Resources)
		mysqlInstance.UserId = types.Int64Value(instance.UserId)
		mysqlInstance.ProdBillTime = types.Int32Value(instance.ProdBillTime)
		mysqlInstance.TenantId = types.StringValue(instance.TenantId)
		mysqlInstance.TplCode = types.StringValue(instance.TplCode)
		mysqlInstance.ProdId = types.Int64Value(instance.ProdId)
		mysqlInstance.ProdInstSetName = types.StringValue(instance.ProdInstSetName)
		mysqlInstance.SecurityGroup = types.StringValue(instance.SecurityGroup)
		mysqlInstance.ProjectID = types.StringValue(instance.ProjectId)
		mysqlInstance.ProjectName = types.StringValue(instance.ProjectName)
		mysqlInstance.PauseEnable = types.BoolValue(instance.PauseEnable)
		mysqlInstance.InstReleaseProtectionStatus = types.Int32Value(instance.InstReleaseProtectionStatus)
		mysqlInstance.LastManualBackUp = types.Int64Value(instance.LastManualBackUp)
		mysqlInstance.RenewalEnable = types.BoolValue(instance.RenewalEnable)
		mysqlInstance.IsMGR = types.Int32Value(instance.IsMGR)
		mysqlInstances = append(mysqlInstances, mysqlInstance)
	}
	config.MysqlInstances = mysqlInstances
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunMysqlInstancesConfig struct {
	RegionID       types.String         `tfsdk:"region_id"`      // 区域id
	ProjectID      types.String         `tfsdk:"project_id"`     //项目id
	PageNow        types.Int32          `tfsdk:"page_now"`       // 当前页，必填
	PageSize       types.Int32          `tfsdk:"page_size"`      // 单页记录条数，必填
	TagVOList      types.List           `tfsdk:"tag_vo_list"`    // 标签列表，选填
	ResDbEngine    types.String         `tfsdk:"res_db_engine"`  // 数据库引擎，选填
	ProdInstName   types.String         `tfsdk:"prod_inst_name"` // 实例名称，选填
	Vip            types.String         `tfsdk:"vip"`            // 连接ip，选填
	MysqlInstances []MysqlInstanceModel `tfsdk:"mysql_instances"`
}

type TageVoModel struct {
	Value   types.String `tfsdk:"value"`    // 标签value，选填
	Key     types.String `tfsdk:"key"`      // 标签key，选填
	LabelID types.String `tfsdk:"label_id"` // k-v的唯一标识，选填
}

type MysqlInstanceModel struct {
	ProdInstName                types.String `tfsdk:"prod_inst_name"`                 // 实例名称
	OuterProdInstId             types.String `tfsdk:"outer_prod_inst_id"`             // 实例ID
	ProdBillType                types.Int32  `tfsdk:"prod_bill_type"`                 // 计费模式 0:按月计费,1:按天计费,2:按年计费,3:按流量计费 4按需计费
	ProdType                    types.Int32  `tfsdk:"prod_type"`                      // 0:单机,1:一主一从,2:一主两从,4:只读实例
	ProdRunningStatus           types.Int32  `tfsdk:"prod_running_status"`            // 0.正常 1.重启中 2.备份中 3.恢复中 4.修改参数中 5.应用参数组中 6.扩容预处理中 7.扩容预处理完成 8.修改端口中 9.迁移中 10.重置密码中 11.修改数据复制方式中 12.缩容预处理中 13.缩容预处理完成 15.内核小版本升级 17.迁移可用区中 18.修改备份配置中 20.停止中 21.已停止 22.启动中 26.白名单配置中
	ProdOrderStatus             types.Int32  `tfsdk:"prod_order_status"`              // 0.正常 1.欠费暂停 2.已注销 3.创建中 4.施工失败 5.到期退订状态 6.新增的状态-openApi暂停 7.创建完成等待变更单 8.待注销 9.手动暂停 10.手动退订
	Alive                       types.Int32  `tfsdk:"alive"`                          // 运行状态 0正常 -1异常
	Vip                         types.String `tfsdk:"vip"`                            // 虚拟IP地址
	Vip6                        types.String `tfsdk:"vip6"`                           // 虚拟IPv6地址
	WritePort                   types.String `tfsdk:"write_port"`                     // 写数据端口
	ReadPort                    types.String `tfsdk:"read_port"`                      // 读端口
	CreateTime                  types.Int64  `tfsdk:"create_time"`                    // 创建时间
	ExpireTime                  types.Int64  `tfsdk:"expire_time"`                    // 到期时间
	MachineSpec                 types.String `tfsdk:"machine_spec"`                   // 实例规格
	ProdDbEngine                types.String `tfsdk:"prod_db_engine"`                 // 数据库引擎
	DiskSize                    types.Int32  `tfsdk:"disk_size"`                      // 存储空间大小 单位G
	NewMysqlVersion             types.String `tfsdk:"new_mysql_version"`              // mysql版本
	DiskType                    types.String `tfsdk:"disk_type"`                      // 存储类型
	BackupDiskUsedRated         types.Int32  `tfsdk:"backup_disk_used_rated"`         // 备份空间使用率
	NetName                     types.String `tfsdk:"net_name"`                       // 虚拟私有云
	Subnet                      types.String `tfsdk:"subnet"`                         // 子网名
	ProdInstFlag                types.String `tfsdk:"prod_inst_flag"`                 // 实例标签
	OrderId                     types.Int64  `tfsdk:"order_id"`                       // 订单ID
	CanOperate                  types.Int32  `tfsdk:"can_operate"`                    // 是否可升级
	TplName                     types.String `tfsdk:"tpl_name"`                       // 实例模板名
	AuditLogStatus              types.Int32  `tfsdk:"audit_log_status"`               // 日志审计开关
	ParameterGroupUsed          types.String `tfsdk:"parameter_group_used"`           // 参数模板名
	VpcID                       types.String `tfsdk:"vpc_id"`                         // 虚拟私有云id
	Usezos                      types.Int32  `tfsdk:"use_zos"`                        // 是否使用对象存储
	ProdInstId                  types.Int64  `tfsdk:"prod_inst_id"`                   // 内部实例id
	DbMysqlVersion              types.String `tfsdk:"db_mysql_version"`               // mysql版本
	Resources                   types.String `tfsdk:"resources"`                      // 订单来源
	UserId                      types.Int64  `tfsdk:"user_id"`                        // 用户id
	ProdBillTime                types.Int32  `tfsdk:"prod_bill_time"`                 // 计费模式类型
	TenantId                    types.String `tfsdk:"tenant_id"`                      // 租户id
	TplCode                     types.String `tfsdk:"tpl_code"`                       // 实例模板code
	ProdId                      types.Int64  `tfsdk:"prod_id"`                        // 产品id
	ProdInstSetName             types.String `tfsdk:"prod_inst_set_name"`             // 实例set名称
	SecurityGroup               types.String `tfsdk:"security_group"`                 // 安全组名
	ProjectID                   types.String `tfsdk:"project_id"`                     // 项目id
	ProjectName                 types.String `tfsdk:"project_name"`                   // 项目名
	PauseEnable                 types.Bool   `tfsdk:"pause_enable"`                   // 是否允许暂停
	InstReleaseProtectionStatus types.Int32  `tfsdk:"inst_release_protection_status"` // 实例释放保护开关
	LastManualBackUp            types.Int64  `tfsdk:"last_manual_back_up"`            // 最后系统全备id
	RenewalEnable               types.Bool   `tfsdk:"renewal_enable"`                 // 是否允许续订
	IsMGR                       types.Int32  `tfsdk:"is_mgr"`                         // 是否为mgr实例
}

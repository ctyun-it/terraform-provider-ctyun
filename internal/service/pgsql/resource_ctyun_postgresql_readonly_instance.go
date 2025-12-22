package pgsql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunPostgresqlReadOnlyInstance{}
	_ resource.ResourceWithConfigure   = &CtyunPostgresqlReadOnlyInstance{}
	_ resource.ResourceWithImportState = &CtyunPostgresqlReadOnlyInstance{}
)

type CtyunPostgresqlReadOnlyInstance struct {
	meta         *common.CtyunMetadata
	ecsService   *business.EcsService
	pgsqlService *business.PgsqlService
}

func (c *CtyunPostgresqlReadOnlyInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[projectID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunPostgresqlReadOnlyInstanceConfig
	var ID, regionId, projectId string
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		projectId = c.meta.GetExtraIfEmpty(projectId, common.ExtraProjectId)
		ID = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &ID, &projectId, &regionId)
		if err != nil {
			return
		}
	}
	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionId)
	config.ProjectID = types.StringValue(projectId)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunPostgresqlReadOnlyInstance) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_postgresql_readonly_instance"
}
func NewCtyunPostgresqlReadOnlyInstance() resource.Resource {
	return &CtyunPostgresqlReadOnlyInstance{}
}

func (c *CtyunPostgresqlReadOnlyInstance) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ecsService = business.NewEcsService(c.meta)
	c.pgsqlService = business.NewPgsqlService(c.meta)
}

func (c *CtyunPostgresqlReadOnlyInstance) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10034019/10322776",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "mysql数据库实例ID，为该实例管理只读实例",
				Validators: []validator.String{
					stringvalidator.LengthBetween(32, 32),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "订购周期类型，取值范围：month：按月，on_demand：按需。当此值为month时，cycle_count为必填",
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypeOnDemand, business.OrderCycleTypeMonth),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int32Attribute{
				Optional:    true,
				Description: "订购时长，该参数当且仅当在cycle_type为month时填写，支持传递1-36",
				Validators: []validator.Int32{
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					int32validator.Between(1, 36),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"auto_renew": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否自动续订，默认非自动续订，当cycle_type不等于on_demand时才可填写，当cycle_count<12，到期自动续订1个月，当cycle_count>=12，到期自动续订12个月",
				Default:     booldefault.StaticBool(false),
				Validators: []validator.Bool{
					validator2.ConflictsWithEqualBool(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
				},
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"flavor_name": schema.StringAttribute{
				Required:    true,
				Description: "规格名称，形如c7.2xlarge.4，可从data.ctyun_mysql_specs查询支持的规格，支持更新。注：只读规格远小于主实例规格时，可能导致创建只读实例失败、复制延迟等风险。",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "实例名称（实例名称 (长度4到100，必须以字母或中文开头，只能包含字母(不区分大小写)、中文、数字、-或_)）",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(4, 100),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"availability_zone_name": schema.StringAttribute{
				Optional:    true,
				Description: "可用区id，如果不填写，默认为第一个可用区",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "可读实例id",
			},
		},
	}
}

func (c *CtyunPostgresqlReadOnlyInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunPostgresqlReadOnlyInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建前，根据inst_id获取基本信息
	_, err = c.getPostgresqlInstanceDetail(ctx, &plan, plan.InstID.ValueString())
	if err != nil {
		return
	}
	// 创建前检查
	err = c.checkSpec(ctx, &plan)
	if err != nil {
		return
	}
	// 开始创建
	err = c.CreatePostgresqlReadOnlyInstance(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，获取mysql详情
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPostgresqlReadOnlyInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunPostgresqlReadOnlyInstanceConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "未找到实例") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunPostgresqlReadOnlyInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunPostgresqlReadOnlyInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunPostgresqlReadOnlyInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	//// 确保主机在退订之前是处于running状态
	//err = c.StartedLoop(ctx, &state)
	//if err != nil {
	//	return
	//}

	err = c.refund(ctx, state)
	if err != nil {
		return
	}
	// 轮询确认时候退订成功
	err = c.refundLoop(ctx, state)
	if err != nil {
		return
	}
	time.Sleep(30 * time.Second)
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
	//err = c.destroyLoop(ctx, state)
	//if err != nil {
	//	return
	//}
	response.Diagnostics.AddWarning("删除MySql集群成功", "集群退订后，若立即删除子网或安全组可能会失败，需要等待底层资源释放")
}

// checkSpec 检查规格
func (c *CtyunPostgresqlReadOnlyInstance) checkSpec(ctx context.Context, plan *CtyunPostgresqlReadOnlyInstanceConfig) error {
	// 根据父实例版本，确定prod id
	plan.prodID = business.PgsqlReadNodeVersionProdIdDict[plan.prodVersion]
	// 先根据spec_name调用云主机规格接口
	_, err := c.ecsService.GetFlavorByName(ctx, plan.FlavorName.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return err
	}

	f := strings.Split(plan.FlavorName.ValueString(), ".")
	hostType := strings.ToUpper(f[0])
	plan.instanceSeries = string(hostType[0]) // S、M 或 C
	if len(hostType) > 2 {
		plan.instanceSeries = hostType
	}
	// 再调用数据库规格接口
	pgsqlFlavor, err := c.pgsqlService.GetPgsqlFlavorByProdIdAndFlavorName(
		ctx,
		plan.prodID,
		plan.FlavorName.ValueString(),
		plan.RegionID.ValueString(),
		plan.instanceSeries,
	)
	if err != nil {
		return err
	}
	plan.prodPerformanceSpec = pgsqlFlavor.ProdPerformanceSpec
	plan.hostType = pgsqlFlavor.Generation

	// 映射关系
	if strings.HasPrefix(plan.hostType, "K") { // 鲲鹏
		plan.cpuType = "KunPeng"
	} else if strings.HasPrefix(plan.hostType, "H") { // 海光
		plan.cpuType = "Hygon"
	} else if strings.HasPrefix(plan.hostType, "F") {
		plan.cpuType = "Phytium"
	} else {
		plan.cpuType = "Intel"
	}
	plan.osType = "ctyunos"
	return nil
}

func (c *CtyunPostgresqlReadOnlyInstance) getPostgresqlInstanceDetail(ctx context.Context, config *CtyunPostgresqlReadOnlyInstanceConfig, id string) (*pgsql.PgsqlDetailResponse, error) {
	detailParams := &pgsql.PgsqlDetailRequest{
		ProdInstId: id,
	}
	detailHeaders := &pgsql.PgsqlDetailRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() {
		detailHeaders.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeaders)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	returnObj := resp.ReturnObj
	engine := returnObj.ProdDbEngine
	parts := strings.Split(engine, "_")
	if len(parts) > 0 {
		config.prodVersion = parts[0] // 输出: 12.22
	} else {
		err = fmt.Errorf("实例版本有误，版本为：%s", engine)
	}
	config.prodVersion = parts[0]
	config.vpcID = returnObj.VpcId
	config.subnetID = returnObj.SubnetId
	config.securityGroupID = returnObj.SecurityGroupId
	config.storageType = returnObj.DiskType
	config.storageSpace = returnObj.DiskSize
	return resp, nil
}

func (c *CtyunPostgresqlReadOnlyInstance) CreatePostgresqlReadOnlyInstance(ctx context.Context, config *CtyunPostgresqlReadOnlyInstanceConfig) error {
	cycleType := config.CycleType.ValueString()
	params := &pgsql.PgsqlCreateRequest{
		InstId:          config.InstID.ValueStringPointer(),
		Name:            config.Name.ValueString(),
		BillMode:        business.MysqlBillMode[cycleType],
		ProdVersion:     config.prodVersion,
		ProdId:          business.PgsqlProdIDDict[config.prodID],
		RegionId:        config.RegionID.ValueString(),
		VpcId:           config.vpcID,
		SubnetId:        config.subnetID,
		SecurityGroupId: config.securityGroupID,
		HostType:        config.hostType,
		Count:           1,
	}

	cpuType := business.MysqlCpuTypeDict[config.cpuType]
	osType := business.MysqlOSTypeDict[config.osType]
	params.CpuType = &cpuType
	params.OsType = &osType

	header := &pgsql.PgsqlCreateRequestHeader{}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectId = config.ProjectID.ValueStringPointer()
	}
	if cycleType == business.OnDemandCycleType {
		params.AutoRenewStatus = 0
	} else {
		params.AutoRenewStatus = map[bool]int32{true: 1, false: 0}[config.AutoRenew.ValueBool()]
	}
	var pgsqlNodeInfos []pgsql.PgsqlCreateRequestNodeInfoList
	pgsqlNodeInfo := pgsql.PgsqlCreateRequestNodeInfoList{}
	pgsqlNodeInfo.NodeType = business.PgsqlNodeTypeReadNode
	pgsqlNodeInfo.InstSpec = business.MysqlInstanceSeriesDict[config.instanceSeries]
	pgsqlNodeInfo.ProdPerformanceSpec = config.prodPerformanceSpec
	pgsqlNodeInfo.StorageType = config.storageType
	pgsqlNodeInfo.StorageSpace = config.storageSpace
	pgsqlNodeInfo.Disks = 1
	var azInfoList []pgsql.PgsqlCreateRequestAvailabilityZoneInfo
	var azInfo pgsql.PgsqlCreateRequestAvailabilityZoneInfo
	azInfo.AvailabilityZoneCount = 1
	azInfo.NodeType = business.PgsqlNodeTypeReadNode
	// 若 az info不为空，用户指定az
	if !config.AvailabilityZoneName.IsNull() {
		azInfo.AvailabilityZoneName = config.AvailabilityZoneName.ValueString()
	} else {
		// 直接放到az1上
		regionAzList, err := c.getAzInfoByRegion(ctx, config)
		if err != nil {
			return err
		}
		if len(regionAzList) < 1 {
			err = errors.New("该资源池AZ信息获取为空，无法直接分配节点AZ信息")
		}
		azInfo.AvailabilityZoneName = regionAzList[0].AvailabilityZoneId
	}
	azInfoList = append(azInfoList, azInfo)
	pgsqlNodeInfo.AvailabilityZoneInfo = azInfoList
	pgsqlNodeInfos = append(pgsqlNodeInfos, pgsqlNodeInfo)
	params.MysqlNodeInfoList = pgsqlNodeInfos
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlCreateApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return err
	} else if resp == nil {
		err2 := fmt.Errorf("pgsql实例(id=%s)创建只读实例失败，接口返回nil，请联系研发确认问题原因！", config.InstID.ValueString())
		return err2
	} else if resp.StatusCode != 200 {
		err2 := fmt.Errorf("API return error. Message: %s", resp.Message)
		return err2
	}
	return nil
}

func (c *CtyunPostgresqlReadOnlyInstance) getAzInfoByRegion(ctx context.Context, config *CtyunPostgresqlReadOnlyInstanceConfig) (regionAzList []mysql.TeledbGetAvailabilityZoneResponseReturnObjData, err error) {
	params := &mysql.TeledbGetAvailabilityZoneRequest{
		RegionId: config.RegionID.ValueString(),
	}
	header := &mysql.TeledbGetAvailabilityZoneRequestHeader{}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err2 := c.meta.Apis.SdkCtMysqlApis.TeledbGetAvailabilityZone.Do(ctx, c.meta.Credential, params, header)
	if err2 != nil {
		err = err2
		return
	} else if resp == nil {
		err = errors.New("查询该资源池AZ信息时，返回为nil。请稍后再试")
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjError
		return
	}
	regionAzList = resp.ReturnObj.Data
	if regionAzList == nil || len(regionAzList) == 0 {
		err = errors.New("查询该资源池AZ信息时，返回为空。请稍后再试")
		return
	}
	return
}

func (c *CtyunPostgresqlReadOnlyInstance) getAndMerge(ctx context.Context, config *CtyunPostgresqlReadOnlyInstanceConfig) error {
	// 判断id是否为空，如果刚刚创建的话，需要轮询列表，查询id
	if config.ID.IsNull() || config.ID.IsUnknown() {
		// 根据inst name 查询
		instanceList, err := c.getPostgresqlInstanceList(ctx, config, config.Name.ValueStringPointer())
		if err != nil {
			return err
		}
		if len(instanceList) > 1 {
			err = fmt.Errorf("数据库实例名称重复！")
			return err
		} else if len(instanceList) == 0 {
			// 若查询不到，说明未创建成功，需要轮询
			instanceList, err = c.CreateLoop(ctx, config)
			if err != nil {
				return err
			}
		}
		instanceReadNodeInfo := instanceList[0]
		config.ID = types.StringValue(instanceReadNodeInfo.ProdInstId)
	}
	// 根据id查询详情
	resp, err := c.getPostgresqlInstanceDetail(ctx, config, config.ID.ValueString())
	if err != nil {
		return err
	}
	config.Name = types.StringValue(resp.ReturnObj.ProdInstName)
	return nil
}

func (c *CtyunPostgresqlReadOnlyInstance) getPostgresqlInstanceList(ctx context.Context, config *CtyunPostgresqlReadOnlyInstanceConfig, name *string) ([]pgsql.PgsqlListResponsePageInfo, error) {
	params := &pgsql.PgsqlListRequest{
		PageNum:  1,
		PageSize: 100,
	}
	if name != nil {
		params.ProdInstName = name
	}
	header := &pgsql.PgsqlListRequestHeader{
		RegionID: config.RegionID.ValueString(),
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() && config.ProjectID.ValueString() != "" {
		header.ProjectID = config.ProjectID.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlListApi.Do(ctx, c.meta.Credential, params, header)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询inst_name=%s的实例信息失败，接口返回nil，请联系研发确认问题原因！", config.Name.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return nil, err
	}
	return resp.ReturnObj.List, nil
}

func (c *CtyunPostgresqlReadOnlyInstance) CreateLoop(ctx context.Context, config *CtyunPostgresqlReadOnlyInstanceConfig, loopCount ...int) ([]pgsql.PgsqlListResponsePageInfo, error) {
	var err error
	var response []pgsql.PgsqlListResponsePageInfo
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return nil, err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			instanceList, err2 := c.getPostgresqlInstanceList(ctx, config, config.Name.ValueStringPointer())
			if err2 != nil {
				return false
			}
			if len(instanceList) > 1 {
				err = fmt.Errorf("查询到多条为名为%s的只读记录！", config.Name.ValueString())
				return false
			}
			if len(instanceList) == 1 {
				runningStatus := instanceList[0].ProdRunningStatus
				orderStatus := instanceList[0].ProdOrderStatus
				if runningStatus == business.MysqlRunningStatusStarted && orderStatus == business.MysqlRunningStatusStarted {
					response = instanceList
					return false
				}
			}
			// 未查询到，继续轮询
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，资源仍未创建或查询到！")
	}
	return response, nil
}

func (c *CtyunPostgresqlReadOnlyInstance) refund(ctx context.Context, state CtyunPostgresqlReadOnlyInstanceConfig) error {
	params := &pgsql.PgsqlRefundRequest{
		InstId: state.ID.ValueString(),
	}
	headers := &pgsql.PgsqlRefundRequestHeader{}
	if !state.ProjectID.IsNull() && !state.ProjectID.IsUnknown() {
		headers.ProjectID = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlRefundApi.Do(ctx, c.meta.Credential, params, headers)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("退订只读节点id=%s失败，接口返回nil，请联系研发确认问题原因！", state.ID)
		return err
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil

}

func (c *CtyunPostgresqlReadOnlyInstance) refundLoop(ctx context.Context, state CtyunPostgresqlReadOnlyInstanceConfig, loopCount ...int) error {
	var err error
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.getPostgresqlInstanceDetail(ctx, &state, state.ID.ValueString())
			if err2 != nil {
				err = err2
				return false
			}
			orderStatus := resp.ReturnObj.ProdOrderStatus
			if orderStatus != business.PgsqlProdOrderStatusRunning {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未创建或查询到！")
	}
	return err
}

func (c *CtyunPostgresqlReadOnlyInstance) destroy(ctx context.Context, state CtyunPostgresqlReadOnlyInstanceConfig) error {
	deleteParams := &pgsql.TeledbDestroyRequest{
		InstId: state.ID.ValueString(),
	}
	deleteHeader := &pgsql.TeledbDestroyRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		deleteHeader.ProjectID = state.ProjectID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtPgsqlApis.PgsqlDestroyApi.Do(ctx, c.meta.Credential, deleteParams, deleteHeader)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

type CtyunPostgresqlReadOnlyInstanceConfig struct {
	InstID               types.String `tfsdk:"instance_id"`
	CycleType            types.String `tfsdk:"cycle_type"`             // 计费模式： 支持on_demand和month
	CycleCount           types.Int32  `tfsdk:"cycle_count"`            // 购买时长：单位月（范围：1-12，24，36）
	AutoRenew            types.Bool   `tfsdk:"auto_renew"`             // 自动续订状态
	FlavorName           types.String `tfsdk:"flavor_name"`            // 规格名称
	RegionID             types.String `tfsdk:"region_id"`              // 资源池id
	ProjectID            types.String `tfsdk:"project_id"`             // 项目id
	Name                 types.String `tfsdk:"name"`                   // 只读实例名称
	AvailabilityZoneName types.String `tfsdk:"availability_zone_name"` // 可用区信息
	ID                   types.String `tfsdk:"id"`
	storageType          string       // 存储类型
	storageSpace         int32        // 存储空间, 磁盘大小100G-2T 步长10G
	vpcID                string
	subnetID             string
	securityGroupID      string
	prodID               string
	prodVersion          string
	osType               string
	cpuType              string
	prodPerformanceSpec  string
	hostType             string
	instanceSeries       string
}

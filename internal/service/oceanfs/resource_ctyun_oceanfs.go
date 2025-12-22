package oceanfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/oceanfs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

type CtyunOceanfs struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunOceanfs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_oceanfs"
}

func (c *CtyunOceanfs) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunOceanfs() resource.Resource {
	return &CtyunOceanfs{}
}

func (c *CtyunOceanfs) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[projectID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunOceanfsConfig
	var ID, regionID, projectID string
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		ID = request.ID
	} else {
		err = terraform_extend.Split(request.ID, &ID, &projectID, &regionID)
		if err != nil {
			return
		}
	}
	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionID)
	config.ProjectID = types.StringValue(projectID)
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunOceanfs) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10088966/10115906",
		Attributes: map[string]schema.Attribute{
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
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议类型，nfs/cifs。nfs 适用于 Linux；cifs 适用于 Windows",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"nfs", "cifs"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "文件系统名称；单账户单资源池下，命名需唯一，只能由数字、“-”、字母组成，不能以数字和“-”开头、且不能以“-”结尾，2~255字符",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 255),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.Int32Attribute{
				Required:    true,
				Description: "文件系统大小（GB），支持更新。取值范围默认为[100,1048576]，实际取值受限于用户剩余容量配额大小。为避免资源浪费，单用户单资源池默认分配500TB容量配额，可提交工单提升配额。",
				Validators: []validator.Int32{
					int32validator.Between(100, 1048576),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "计费类型，year/month/on_demand。不支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.SfsCycleType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "包周期数，cycle_type是year或month时必须指定，周期最大长度不能超过3年",
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeMonth),
						types.StringValue(business.OrderCycleTypeYear),
					),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("cycle_type"),
						types.StringValue(business.OrderCycleTypeOnDemand),
					),
					validator2.CycleCount(1, 11, 1, 3),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称，若未填写，默认从环境变量中读取。",
				// az有必要设定默认值
				Default: defaults.AcquireFromGlobalString(common.ExtraAzName, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_id": schema.StringAttribute{
				Description: "VPC ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Description: "子网ID，当isVpce为true时必填",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"is_vpce": schema.BoolAttribute{
				Description: "创建文件系统时是否自动创建VPC终端节点。开启后本服务将为您创建免费的VPC终端节点（VPCE），连接文件存储服务。创建VPCE后将返回该VPC专属的挂载地址，通常需要1~3分钟。注：物理机必须通过VPCE专属挂载地址访问文件系统，其它计算服务如云主机、容器为非必须",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.SetNestedAttribute{
				Description: "标签列表",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "标签键",
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"value": schema.StringAttribute{
							Description: "标签值",
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description: "资源ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Description: "文件系统状态",
				Computed:    true,
			},
			"used_size": schema.Int32Attribute{
				Description: "已使用大小（GB）",
				Computed:    true,
			},
			"create_time": schema.StringAttribute{
				Description: "创建时间，为UTC格式",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Description: "更新时间，为UTC格式",
				Computed:    true,
			},
			"expire_time": schema.StringAttribute{
				Description: "到期时间，为UTC格式，按需时为空",
				Computed:    true,
			},
			"share_path": schema.StringAttribute{
				Computed:    true,
				Description: "挂载路径",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"share_path_windows": schema.StringAttribute{
				Computed:    true,
				Description: "挂载路径（windows）",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunOceanfs) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunOceanfsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunOceanfs) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunOceanfsConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "未找到") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunOceanfs) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunOceanfsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunOceanfsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunOceanfs) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunOceanfsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunOceanfs) create(ctx context.Context, config *CtyunOceanfsConfig) error {
	params := &oceanfs.OceanfsNewSfsRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		SfsType:     "massive",
		SfsProtocol: config.SfsProtocol.ValueString(),
		SfsName:     config.Name.ValueString(),
		SfsSize:     config.SfsSize.ValueInt32(),
		Vpc:         config.VpcID.ValueString(),
	}
	if config.CycleType.ValueString() == business.SfsOnDemandCycleType {
		trueVar := true
		params.OnDemand = &trueVar
	} else {
		falseVar := false
		params.OnDemand = &falseVar
		params.CycleType = config.CycleType.ValueString()
		params.CycleCount = int32(config.CycleCount.ValueInt64())
	}
	if !config.AzName.IsNull() && !config.AzName.IsUnknown() {
		params.AzName = config.AzName.ValueString()
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() && config.ProjectID.ValueString() != "" {
		params.ProjectID = config.ProjectID.ValueString()
	}
	if config.IsVpce.ValueBool() {
		params.IsVpce = config.IsVpce.ValueBoolPointer()

	}
	if !config.SubnetID.IsNull() && !config.SubnetID.IsUnknown() {
		params.Subnet = config.SubnetID.ValueString()
	}
	// 处理tags
	if !config.Tags.IsNull() && !config.Tags.IsUnknown() {
		var tags []CtyunOceanfsTagModel
		diags := config.Tags.ElementsAs(ctx, &tags, false)
		if diags.HasError() {
			err := fmt.Errorf(diags[0].Detail())
			return err
		}

		var tagsParam []*oceanfs.LabelRequest
		for _, tag := range tags {
			tagsParam = append(tagsParam, &oceanfs.LabelRequest{
				Key:   tag.Key.ValueString(),
				Value: tag.Value.ValueString(),
			})
		}
	}

	_, err := c.createReq(ctx, params)
	if err != nil {
		return err
	}

	// 轮询确认是否创建成功
	_, err = c.createLoop(ctx, config, params, 60)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunOceanfs) createLoop(ctx context.Context, config *CtyunOceanfsConfig, params *oceanfs.OceanfsNewSfsRequest, loopCount ...int) (*oceanfs.OceanfsNewSfsResponse, error) {
	var err error
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
			// 轮询创建接口，当返回的id不为空，这跳出循环
			resp, err2 := c.createReq(ctx, params)
			if err2 != nil {
				if !strings.Contains(err2.Error(), "order in progress") {
					err = err2
					return false
				}
				return true
			}
			if len(resp.ReturnObj.Resources) > 1 {
				err = fmt.Errorf("轮询创建时，接口返回多个海量文件服务实例信息")
				return false
			}
			if len(resp.ReturnObj.Resources) == 0 {
				return true
			}
			id := resp.ReturnObj.Resources[0].SfsUID
			if id != nil {
				config.ID = types.StringValue(*id)
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return nil, errors.New("轮询已达最大次数，资源仍未创建成功！")
	}
	return nil, err
}

func (c *CtyunOceanfs) createReq(ctx context.Context, params *oceanfs.OceanfsNewSfsRequest) (*oceanfs.OceanfsNewSfsResponse, error) {
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsNewSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("创建海量文件服务Oceanfs失败，返回结果为空，请联系研发确认问题原因！")
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode && !strings.Contains(resp.ErrorCode, "sfs.order.inProgress") {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

func (c *CtyunOceanfs) getAndMerge(ctx context.Context, config *CtyunOceanfsConfig) error {
	resp, err := c.getOceanfsDetail(ctx, config)
	if err != nil {
		return err
	}
	returnObj := resp.ReturnObj
	config.Name = types.StringValue(returnObj.SfsName)
	//config.SfsType = types.StringValue(returnObj.SfsType)
	config.SfsProtocol = types.StringValue(returnObj.SfsProtocol)
	config.SfsSize = types.Int32Value(returnObj.SfsSize)
	config.AzName = types.StringValue(returnObj.AzName)
	config.Status = types.StringValue(returnObj.SfsStatus)
	config.UsedSize = types.Int32Value(returnObj.UsedSize)
	config.CreateTime = types.StringValue(utils.FromUnixToUTC(returnObj.CreateTime))
	config.UpdateTime = types.StringValue(utils.FromUnixToUTC(returnObj.UpdateTime))
	config.ExpireTime = types.StringValue(utils.FromUnixToUTC(returnObj.ExpireTime))
	config.SharePath = types.StringValue(returnObj.SharePath)
	config.SharePathWin = types.StringValue(returnObj.WindowsSharePath)
	if config.Tags.IsNull() || config.Tags.IsUnknown() {
		config.Tags = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"key":   types.StringType,
				"value": types.StringType,
			},
		})
	}
	return nil
}

func (c *CtyunOceanfs) update(ctx context.Context, state *CtyunOceanfsConfig, plan *CtyunOceanfsConfig) error {
	if plan.SfsSize.Equal(state.SfsSize) {
		return nil
	}
	params := &oceanfs.OceanfsResizeSfsRequest{
		SfsSize:     plan.SfsSize.ValueInt32(),
		SfsUID:      state.ID.ValueString(),
		RegionID:    state.RegionID.ValueString(),
		ClientToken: uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsResizeSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("扩容海量文件服务Oceanfs失败(id=%s)，返回结果为空，请联系研发确认问题原因！", state.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode && !strings.Contains(resp.ErrorCode, "inProgress") {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}

	// 轮询确认是否扩容成功
	err = c.resizeLoop(ctx, state, plan, 60)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunOceanfs) delete(ctx context.Context, config CtyunOceanfsConfig) error {
	params := &oceanfs.OceanfsRefundSfsRequest{
		ClientToken: uuid.NewString(),
		SfsUID:      config.ID.ValueString(),
		RegionID:    config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsRefundSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除海量文件服务Oceanfs失败(id=%s)，返回结果为空，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return err
	}
	return nil
}

func (c *CtyunOceanfs) resizeLoop(ctx context.Context, state *CtyunOceanfsConfig, plan *CtyunOceanfsConfig, loopCount ...int) error {
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
			// 轮询详情接口，确认sfs size是否与plan.sfsSize对应
			resp, err2 := c.getOceanfsDetail(ctx, state)
			if err2 != nil {
				if !strings.Contains(err2.Error(), "order in progress") {
					err = err2
					return false
				}
				return true
			}
			sfsSize := resp.ReturnObj.SfsSize
			if sfsSize == plan.SfsSize.ValueInt32() {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，弹性文件系统仍未扩容成功！")
	}
	return err
}

func (c *CtyunOceanfs) getOceanfsDetail(ctx context.Context, config *CtyunOceanfsConfig) (*oceanfs.OceanfsInfoSfsResponse, error) {
	params := &oceanfs.OceanfsInfoSfsRequest{
		SfsUID:   config.ID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkOceanfsApis.OceanfsInfoSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取海量文件服务Oceanfs详情失败(id=%s)，返回结果为空，请联系研发确认问题原因！", config.ID.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

type CtyunOceanfsConfig struct {
	RegionID  types.String `tfsdk:"region_id"`
	ProjectID types.String `tfsdk:"project_id"`
	//SfsType     types.String `tfsdk:"type"`
	SfsProtocol  types.String `tfsdk:"protocol"`
	Name         types.String `tfsdk:"name"`
	SfsSize      types.Int32  `tfsdk:"size"`
	CycleType    types.String `tfsdk:"cycle_type"`
	CycleCount   types.Int64  `tfsdk:"cycle_count"`
	AzName       types.String `tfsdk:"az_name"`
	VpcID        types.String `tfsdk:"vpc_id"`
	SubnetID     types.String `tfsdk:"subnet_id"`
	IsVpce       types.Bool   `tfsdk:"is_vpce"`
	Tags         types.Set    `tfsdk:"tags"`
	ID           types.String `tfsdk:"id"`
	Status       types.String `tfsdk:"status"`
	UsedSize     types.Int32  `tfsdk:"used_size"`
	CreateTime   types.String `tfsdk:"create_time"`
	UpdateTime   types.String `tfsdk:"update_time"`
	ExpireTime   types.String `tfsdk:"expire_time"`
	SharePath    types.String `tfsdk:"share_path"`
	SharePathWin types.String `tfsdk:"share_path_windows"`
}

type CtyunOceanfsTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

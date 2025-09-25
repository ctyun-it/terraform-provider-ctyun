package hpfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/hpfs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

type ctyunHpfs struct {
	meta *common.CtyunMetadata
}

func (c *ctyunHpfs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_hpfs"
}

func (c *ctyunHpfs) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func NewCtyunHpfsInstance() resource.Resource {
	return &ctyunHpfs{}
}

func (c *ctyunHpfs) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10088932/10090437**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
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
			//"sfs_type": schema.StringAttribute{
			//	Required:    true,
			//	Description: "并行文件类型 (HPC性能型)",
			//	Validators: []validator.String{
			//		stringvalidator.OneOf("HPC性能型"),
			//	},
			//	PlanModifiers: []planmodifier.String{
			//		stringplanmodifier.RequiresReplace(),
			//	},
			//},
			"sfs_protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议类型，可以根据data.ctyun_hpfs_clusters接口查询，也可访问网页查询：https://www.ctyun.cn/document/10088932/10510589",
				Validators: []validator.String{
					stringvalidator.OneOf("nfs", "hpfs"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "订购周期类型，只支持on_demand",
				Default:     stringdefault.StaticString(business.OrderCycleTypeOnDemand),
				Validators: []validator.String{
					stringvalidator.OneOf(business.OrderCycleTypeOnDemand),
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "并行文件名，仅允许英文字母数字及-，开头必须为字母，结尾不允许为-，且长度为2-255字符，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 255),
				},
			},
			"sfs_size": schema.Int32Attribute{
				Required:    true,
				Description: "文件大小（GB），范围: 500-32768。支持更新",
				Validators: []validator.Int32{
					// 范围验证
					int32validator.Between(500, 32768),
					// 自定义步长验证
					validator2.SfsSize(),
				},
			},
			"az_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "可用区名称",
				Default:     defaults.AcquireFromGlobalString(common.ExtraAzName, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"cluster_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "集群名称，仅资源池支持指定集群时可传入该参数。可以根据data.ctyun_hpfs_clusters接口查询，也可访问网页查询：https://www.ctyun.cn/document/10088932/10510589",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"baseline": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "性能基线（MB/s/TB），仅资源池支持性能基线时可传入该参数。可以根据data.ctyun_hpfs_clusters接口查询，也可访问网页查询：https://www.ctyun.cn/document/10088932/10510589",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "虚拟网 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Optional:    true,
				Description: "子网 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "订单ID",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "资源 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sfs_status": schema.StringAttribute{
				Computed:    true,
				Description: "并行文件状态",
			},
			"used_size": schema.Int32Attribute{
				Computed:    true,
				Description: "已用大小（MB）",
			},
			"dataflow_list": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "HPFS文件系统下的数据流动策略ID列表",
			},
		},
	}
}

func (c *ctyunHpfs) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunHpfsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	//创建前检查,检查证书有效性
	isValid, err := c.checkBeforeHpfs(ctx, plan)
	if !isValid || err != nil {
		return
	}
	createParams, err := c.createHpfs(ctx, &plan)
	if err != nil {
		return
	}
	// 创建后，通过创建的请求轮询，确认创建成功
	err = c.createLoop(ctx, &plan, createParams, 60)
	if err != nil {
		return
	}
	// 创建后反查创建后的证书信息
	err = c.getAndMergeHpfs(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunHpfs) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunHpfsConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeHpfs(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotExists") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunHpfs) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取 plan -tf文件中配置
	var plan CtyunHpfsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunHpfsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}

	// 更新基本信息
	err = c.updateHfps(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeHpfs(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunHpfs) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunHpfsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &hpfs.HpfsRefundSfsRequest{
		RegionID: state.RegionID.ValueString(),
		SfsUID:   state.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkHpfsApis.HpfsRefundSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = errors.New("hpfs退订失败，返回值为nil")
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	// 异步接口，需要轮询查看是否退订成功
	err = c.deleteLoop(ctx, &state, 600)
	if err != nil {
		return
	}
}

func (c *ctyunHpfs) checkBeforeHpfs(ctx context.Context, plan CtyunHpfsConfig) (inValid bool, err error) {
	// 判断sfs_type，sfs_protocol是否合理，/v4/hpfs/list-cluster

	return true, nil
}

func (c *ctyunHpfs) createHpfs(ctx context.Context, config *CtyunHpfsConfig) (*hpfs.HpfsNewSfsRequest, error) {
	params := &hpfs.HpfsNewSfsRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		SfsType:     "hpfs_perf",
		SfsProtocol: config.SfsProtocol.ValueString(),
		CycleType:   config.CycleType.ValueString(),
		SfsName:     config.Name.ValueString(),
		SfsSize:     config.SfsSize.ValueInt32(),
		Vpc:         config.VpcID.ValueString(),
		Subnet:      config.SubnetID.ValueString(),
	}
	if config.CycleType.ValueString() == business.HpfsCycleTypeOnDemand {
		onDemand := true
		params.OnDemand = &onDemand
	} else {
		params.CycleCount = config.CycleCount.ValueInt32()
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	if !config.AzName.IsNull() && !config.AzName.IsUnknown() {
		params.AzName = config.AzName.ValueString()
	}
	if !config.ClusterName.IsNull() && !config.ClusterName.IsUnknown() {
		params.ClusterName = config.ClusterName.ValueString()
	}
	if !config.Baseline.IsNull() && !config.Baseline.IsUnknown() {
		params.Baseline = config.Baseline.ValueString()
	}
	resp, err := c.meta.Apis.SdkHpfsApis.HpfsNewSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = errors.New("开通hpfs失败，返回nil")
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		if strings.Contains(resp.Message, "in progress") {
			err = nil
		} else {
			err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
			return nil, err
		}
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	config.MasterOrderID = types.StringValue(resp.ReturnObj.MasterOrderID)
	//config.ID = types.StringValue(resp.ReturnObj.Resources[0].SfsUID)
	return params, nil
}

func (c *ctyunHpfs) getAndMergeHpfs(ctx context.Context, config *CtyunHpfsConfig) error {
	// 获取hpfs详情
	hpfsResp, err := c.getHpfsDetail(ctx, config)
	if err != nil {
		return err
	}
	hpfsDetail := hpfsResp.ReturnObj
	config.Name = types.StringValue(hpfsDetail.SfsName)
	config.SfsSize = types.Int32Value(hpfsDetail.SfsSize)
	config.SfsStatus = types.StringValue(hpfsDetail.SfsStatus)
	config.ClusterName = types.StringValue(hpfsDetail.ClusterName)
	config.UsedSize = types.Int32Value(hpfsDetail.UsedSize)
	config.Baseline = types.StringValue(hpfsDetail.Baseline)
	dataFlowList, diags := types.SetValueFrom(ctx, types.StringType, hpfsDetail.DataflowList)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return err
	}
	config.DataflowList = dataFlowList
	return nil
}

func (c *ctyunHpfs) getHpfsDetail(ctx context.Context, config *CtyunHpfsConfig) (*hpfs.HpfsInfoSfsResponse, error) {
	params := &hpfs.HpfsInfoSfsRequest{
		SfsUID:   config.ID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkHpfsApis.HpfsInfoSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = errors.New("获取hpfs详情失败，返回为nil")
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

func (c *ctyunHpfs) updateHfps(ctx context.Context, state *CtyunHpfsConfig, plan *CtyunHpfsConfig) error {
	// 并行文件重命名
	err := c.hfpsRename(ctx, state, plan)
	if err != nil {
		return err
	}
	// 并行文件修改规格
	err = c.updateHpfsSize(ctx, state, plan)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunHpfs) hfpsRename(ctx context.Context, state *CtyunHpfsConfig, plan *CtyunHpfsConfig) error {
	if plan.Name.IsNull() || state.Name == plan.Name {
		return nil
	}
	params := &hpfs.HpfsRenameSfsRequest{
		RegionID: state.RegionID.ValueString(),
		SfsUID:   state.ID.ValueString(),
		SfsName:  plan.Name.ValueString(),
	}
	resp, err := c.meta.Apis.SdkHpfsApis.HpfsRenameSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("hpfs 更名失败，返回nil")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	err = c.renameLoop(ctx, state, plan, 60)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunHpfs) updateHpfsSize(ctx context.Context, state *CtyunHpfsConfig, plan *CtyunHpfsConfig) error {
	// 判断是否需要进行修改
	if plan.SfsSize.IsNull() || state.SfsSize == plan.SfsSize {
		return nil
	}
	if plan.SfsSize.ValueInt32() < state.SfsSize.ValueInt32() {
		return errors.New("并行文件暂不支持缩容能力")
	}
	// state和plan阶段sfs_size相同，不触发变配
	if plan.SfsSize.ValueInt32() == state.SfsSize.ValueInt32() {
		return nil
	}
	// 配置修改参数
	params := &hpfs.HpfsResizeSfsRequest{
		SfsSize:     plan.SfsSize.ValueInt32(),
		SfsUID:      state.ID.ValueString(),
		RegionID:    state.RegionID.ValueString(),
		ClientToken: uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkHpfsApis.HpfsResizeSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("hpfs sfs_size修改失败，返回值为Nil。")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		if strings.Contains(resp.Message, " in progress") {
			err = nil
		} else {
			err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
			return err
		}
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}

	// 变配完成后，轮询确认升级完成
	err = c.updateHpfsSizeLoop(ctx, state, plan)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunHpfs) updateHpfsSizeLoop(ctx context.Context, state *CtyunHpfsConfig, plan *CtyunHpfsConfig, loopCount ...int) error {
	var err error
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.getHpfsDetail(ctx, state)
			if err2 != nil {
				err = err2
				return false
			}
			if resp.ReturnObj.SfsSize == plan.SfsSize.ValueInt32() {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未变配成功！")
	}
	return err
}

func (c *ctyunHpfs) deleteLoop(ctx context.Context, config *CtyunHpfsConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.getHpfsDetail(ctx, config)
			if err2 != nil {
				if strings.Contains(err2.Error(), "资源不存在") {
					err = nil
				} else {
					err = err2
				}
				return false
			} else if resp == nil {
				return false
			} else if resp.StatusCode != common.NormalStatusCode {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未退订成功！")
	}
	return
}

func (c *ctyunHpfs) createLoop(ctx context.Context, plan *CtyunHpfsConfig, params *hpfs.HpfsNewSfsRequest, loopCount ...int) error {
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
			// 轮询创建接口，当返回的id不为空，这跳出循环
			resp, err2 := c.meta.Apis.SdkHpfsApis.HpfsNewSfsApi.Do(ctx, c.meta.SdkCredential, params)
			if err2 != nil {
				err2 = err
				return false
			} else if resp == nil {
				err = errors.New("开通hpfs失败，返回nil")
				return false
			} else if resp.StatusCode != common.NormalStatusCode {
				if strings.Contains(resp.Message, "in progress") {
					err = nil
					return true
				} else {
					err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
					return false
				}
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			id := resp.ReturnObj.Resources[0].SfsUID
			if id != "" {
				plan.ID = types.StringValue(id)
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未创建成功！")
	}
	return err

}

func (c *ctyunHpfs) renameLoop(ctx context.Context, state *CtyunHpfsConfig, plan *CtyunHpfsConfig, loopCount ...int) error {
	var err error
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*5, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.getHpfsDetail(ctx, state)
			if err2 != nil {
				err = err2
				return false
			}
			if resp.ReturnObj.SfsName == plan.Name.ValueString() {
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未更名成功！")
	}
	return err
}

type CtyunHpfsConfig struct {
	RegionID  types.String `tfsdk:"region_id"`  // 资源池 ID
	ProjectID types.String `tfsdk:"project_id"` // 资源所属企业项目 ID
	//SfsType       types.String `tfsdk:"sfs_type"`        // 并行文件类型
	SfsProtocol   types.String `tfsdk:"sfs_protocol"`    // 协议类型
	CycleType     types.String `tfsdk:"cycle_type"`      // 包周期类型
	CycleCount    types.Int32  `tfsdk:"cycle_count"`     // 包周期数
	Name          types.String `tfsdk:"name"`            // 并行文件名
	SfsSize       types.Int32  `tfsdk:"sfs_size"`        // 文件大小（GB）
	AzName        types.String `tfsdk:"az_name"`         // 可用区名称
	ClusterName   types.String `tfsdk:"cluster_name"`    // 集群名称
	Baseline      types.String `tfsdk:"baseline"`        // 性能基线
	VpcID         types.String `tfsdk:"vpc_id"`          // 虚拟网 ID
	SubnetID      types.String `tfsdk:"subnet_id"`       // 子网 ID
	MasterOrderID types.String `tfsdk:"master_order_id"` // 订单id
	ID            types.String `tfsdk:"id"`              // 资源 ID
	SfsStatus     types.String `tfsdk:"sfs_status"`      // 并行文件状态
	UsedSize      types.Int32  `tfsdk:"used_size"`       // 已用大小（MB）
	DataflowList  types.Set    `tfsdk:"dataflow_list"`   // HPFS文件系统下的数据流动策略ID列表
}

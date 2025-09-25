package sfs

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sfs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

type ctyunSfs struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *ctyunSfs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_sfs"
}

func (c *ctyunSfs) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunSfs() resource.Resource {
	return &ctyunSfs{}
}

func (c *ctyunSfs) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027350**`,
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
			"is_encrypt": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否加密盘，默认false，支持更新。目前仅少量资源池支持加密。具体可查看产品能力地图：https://www.ctyun.cn/document/10027350/10693922",
				Default:     booldefault.StaticBool(false),
			},
			"kms_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "如果是加密盘，需要提供kms的uuid，支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("is_encrypt"),
						types.BoolValue(true),
					),
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
			"sfs_type": schema.StringAttribute{
				Required:    true,
				Description: "存储类型，capacity(标准型)或performance（性能型）",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"capacity", "performance"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sfs_protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议类型，nfs/cifs",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"nfs", "cifs"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "文件系统名称；单账户单资源池下，命名需唯一，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"sfs_size": schema.Int32Attribute{
				Required:    true,
				Description: "大小，单位GB，取值范围：[500GB, 32768GB]。支持更新。弹性文件只支持扩容，不支持缩容",
				Validators: []validator.Int32{
					int32validator.Between(500, 32768),
				},
			},
			"cycle_type": schema.StringAttribute{
				Required:    true,
				Description: "包周期类型，year/month/on_demand；onDemand为false时，必须指定。不支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.SfsCycleType...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle_count": schema.Int64Attribute{
				Optional:    true,
				Description: "包周期数。onDemand为false时必须指定；周期最大长度不能超过3年",
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
				Description: "实例部署的az信息。多可用区资源池下，若不填写，将随机分配AZ",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "虚拟私有云ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:    true,
				Description: "子网ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.SubnetValidate(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "弹性文件系统id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "弹性文件系统状态",
			},
			"used_size": schema.Int32Attribute{
				Computed:    true,
				Description: "弹性文件系统已使用大小（MB）",
			},
			"read_only": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "弹性文件系统是否只读。默认为false。支持更新，true-只读；false-可读写。sfs_protocol=cifs仅支持为false",
				Validators: []validator.Bool{
					validator2.ConflictsWithEqualBool(
						path.MatchRoot("sfs_protocol"),
						types.StringValue("cifs"),
					),
				},
			},
		},
	}
}

func (c *ctyunSfs) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunSfsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.createSfs(ctx, &plan)
	if err != nil {
		return
	}
	// 如果创建时，read_only不为空，需要调用设置下
	err = c.setSfsRw(ctx, &plan, &plan)
	if err != nil {
		return
	}
	// 创建后反查创建后的证书信息
	err = c.getAndMergeSfs(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunSfs) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunSfsConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeSfs(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "未找到") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunSfs) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取 plan -tf文件中配置
	var plan CtyunSfsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunSfsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}
	err = c.updateSfs(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeSfs(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunSfs) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunSfsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &sfs.SfsSfsRefundSfsRequest{
		ClientToken: uuid.NewString(),
		SfsUID:      state.ID.ValueString(),
		RegionID:    state.RegionID.ValueString(),
	}

	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsRefundSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("sfs id为%s的弹性文件服务退订失败，接口返回nil。请与研发联系确认原因。", state.ID.ValueString())
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	time.Sleep(30 * time.Second)
	return
}

func (c *ctyunSfs) createSfs(ctx context.Context, config *CtyunSfsConfig) error {
	params := &sfs.SfsSfsNewSfsRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
		SfsType:     config.SfsType.ValueString(),
		SfsProtocol: config.SfsProtocol.ValueString(),
		SfsName:     config.Name.ValueString(),
		SfsSize:     config.SfsSize.ValueInt32(),
		Vpc:         config.VpcID.ValueString(),
		Subnet:      config.SubnetID.ValueString(),
	}

	if config.CycleType.ValueString() == business.SfsOnDemandCycleType {
		onDemand := true
		params.OnDemand = &onDemand
	} else {
		onDemand := false
		params.OnDemand = &onDemand
		params.CycleType = config.CycleType.ValueString()
		params.CycleCount = int32(config.CycleCount.ValueInt64())
	}
	// 确认是否为多az
	zones, err2 := c.regionService.GetZonesByRegionID(ctx, config.RegionID.ValueString())
	if err2 != nil {
		return err2
	}
	isNaz := false
	if len(zones) > 1 {
		isNaz = true
	}
	if isNaz {
		if config.AzName.IsNull() || config.AzName.IsUnknown() {
			params.AzName = zones[0]
			//err := fmt.Errorf("当资源池为多AZ，创建sfs需要指定AZ。")
			//return err
		} else {
			params.AzName = config.AzName.ValueString()
		}
	}
	if !config.IsEncrypt.IsNull() && !config.IsEncrypt.IsUnknown() {
		params.IsEncrypt = config.IsEncrypt.ValueBoolPointer()
		if config.IsEncrypt.ValueBool() {
			params.KmsUUID = config.KmsUUID.ValueString()
		}
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	_, err := c.requestCreateSfsApi(ctx, params)
	if err != nil {
		if !strings.Contains(err.Error(), "order in progress") {
			return err
		}
	}
	// 轮询确认是否创建成功
	err = c.createLoop(ctx, config, params, 60)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunSfs) requestCreateSfsApi(ctx context.Context, params *sfs.SfsSfsNewSfsRequest) (*sfs.SfsSfsNewSfsResponse, error) {
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsNewSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("创建sfs失败，接口返回nil。可联系研发确认原因")
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

func (c *ctyunSfs) createLoop(ctx context.Context, config *CtyunSfsConfig, params *sfs.SfsSfsNewSfsRequest, loopCount ...int) error {
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
			resp, err2 := c.requestCreateSfsApi(ctx, params)
			if err2 != nil {
				if !strings.Contains(err2.Error(), "order in progress") {
					err = err2
					return false
				}
				return true
			}
			id := resp.ReturnObj.Resources[0].SfsUID
			if id != "" {
				config.ID = types.StringValue(id)
				return false
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，资源仍未创建成功！")
	}
	return err
}

func (c *ctyunSfs) getAndMergeSfs(ctx context.Context, config *CtyunSfsConfig) error {
	resp, err := c.getSfsDetail(ctx, config)
	if err != nil {
		return err
	}
	returnObj := resp.ReturnObj
	config.Name = types.StringValue(returnObj.SfsName)
	config.UsedSize = types.Int32Value(returnObj.UsedSize)
	config.Status = types.StringValue(returnObj.SfsStatus)
	config.SfsSize = types.Int32Value(returnObj.SfsSize)
	config.SfsProtocol = types.StringValue(returnObj.SfsProtocol)
	config.SfsType = types.StringValue(returnObj.SfsType)
	//config.AzName = types.StringValue(returnObj.AzName)
	// 获取是否只读
	rwResp, err := c.getSfsRwDetail(ctx, config)
	if err != nil {
		return err
	}
	config.ReadOnly = types.BoolValue(*rwResp.ReturnObj.List[0].ReadOnly.Value)
	//if rwResp.ReturnObj.List[0].ReadOnly == "true" {
	//	config.ReadOnly = types.BoolValue(true)
	//} else {
	//	config.ReadOnly = types.BoolValue(false)
	//}

	return nil
}

func (c *ctyunSfs) getSfsDetail(ctx context.Context, config *CtyunSfsConfig) (*sfs.SfsSfsCreateSfsResponse, error) {
	params := &sfs.SfsSfsCreateSfsRequest{
		SfsUID:   config.ID.ValueString(),
		RegionID: config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsCreateSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询sfs失败，接口返回nil。sfs id 为：%s可联系研发确认原因", config.ID.ValueString())
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

func (c *ctyunSfs) getSfsRwDetail(ctx context.Context, config *CtyunSfsConfig) (*sfs.SfsSfsListReadWriteSfs1Response, error) {
	params := &sfs.SfsSfsListReadWriteSfs1Request{
		RegionID: config.RegionID.ValueString(),
		SfsUID:   config.ID.ValueString(),
		PageSize: 10,
		PageNo:   1,
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsListReadWriteSfs1Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询sfs id 为%s的读写信息失败，接口返回为nil，请联系研发确认失败原因。", config.ID.ValueString())
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

func (c *ctyunSfs) updateSfs(ctx context.Context, state *CtyunSfsConfig, plan *CtyunSfsConfig) error {
	// 更新sfs名称
	err := c.updateSfsName(ctx, state, plan)
	if err != nil {
		return err
	}
	// 扩容sfs
	err = c.ResizeSfs(ctx, state, plan)
	if err != nil {
		return err
	}

	// 设置文件系统是否已读
	err = c.setSfsRw(ctx, state, plan)
	if err != nil {
		return err
	}

	return nil
}

func (c *ctyunSfs) updateSfsName(ctx context.Context, state *CtyunSfsConfig, plan *CtyunSfsConfig) error {
	if plan.Name.Equal(state.Name) {
		return nil
	}
	params := &sfs.SfsRenameSFSRequest{
		RegionID: state.RegionID.ValueString(),
		SfsUID:   state.ID.ValueString(),
		SfsName:  plan.Name.ValueString(),
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsRenameSFSApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		return fmt.Errorf("修改id为%s的弹性文件系统名称失败，接口返回nil。请与研发联系确认问题原因。", state.ID.ValueString())
	} else if resp.StatusCode != common.NormalStatusCode {
		return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
	}

	// 复核sfs name是否修改成功
	err = c.renameLoop(ctx, *plan)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunSfs) renameLoop(ctx context.Context, plan CtyunSfsConfig) error {
	var err error
	retryer, err := business.NewRetryer(time.Second*10, 60)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			var resp *sfs.SfsSfsCreateSfsResponse
			resp, err = c.getSfsDetail(ctx, &plan)
			if err != nil {
				return false
			}
			if resp.ReturnObj.SfsName != plan.Name.ValueString() {
				return true
			}
			return false
		})
	if err != nil {
		return err
	}
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，弹性文件系统名称未修改成功！")
	}
	return err
}

func (c *ctyunSfs) ResizeSfs(ctx context.Context, state *CtyunSfsConfig, plan *CtyunSfsConfig) error {
	if plan.SfsSize.Equal(state.SfsSize) {
		return nil
	}
	params := &sfs.SfsSfsResizeSfsRequest{
		SfsSize:     plan.SfsSize.ValueInt32(),
		SfsUID:      state.ID.ValueString(),
		RegionID:    state.RegionID.ValueString(),
		ClientToken: uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkSfsApi.SfsSfsResizeSfsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		if !strings.Contains(err.Error(), "order in progress") {
			return err
		}
	} else if resp == nil {
		return fmt.Errorf("扩容id为%s的弹性文件系统失败，接口返回nil。请与研发联系确认问题原因", state.ID.ValueString())
	} else if resp.StatusCode != common.NormalStatusCode {
		if !strings.Contains(resp.Error, "Sfs.Order.InProgress") {
			return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		}
	}
	// 轮询确认是否扩容成功
	err = c.resizeLoop(ctx, state, plan, 60)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunSfs) resizeLoop(ctx context.Context, state *CtyunSfsConfig, plan *CtyunSfsConfig, loopCount ...int) error {
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
			resp, err2 := c.getSfsDetail(ctx, state)
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

func (c *ctyunSfs) setSfsRw(ctx context.Context, state *CtyunSfsConfig, plan *CtyunSfsConfig) error {
	if plan.ReadOnly.IsNull() || plan.ReadOnly.IsUnknown() {
		return nil
	}

	// 如果需要已读
	if plan.ReadOnly.ValueBool() {
		params := &sfs.SfsSfsSetReadSfsRequest{
			RegionID: state.RegionID.ValueString(),
			SfsUID:   state.ID.ValueString(),
		}
		resp, err := c.meta.Apis.SdkSfsApi.SfsSfsSetReadSfsApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return err
		} else if resp == nil {
			return fmt.Errorf("设置id为%s的弹性文件服务为只读实例失败，接口返回nil。请联系研发确认问题原因。", state.ID.ValueString())
		} else if resp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		}
	} else {
		params := &sfs.SfsSfsSetReadWriteSfsRequest{
			RegionID: state.RegionID.ValueString(),
			SfsUID:   state.ID.ValueString(),
		}
		resp, err := c.meta.Apis.SdkSfsApi.SfsSfsSetReadWriteSfsApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return err
		} else if resp == nil {
			return fmt.Errorf("设置id为%s的弹性文件服务为读写实例失败，接口返回nil。请联系研发确认问题原因。", state.ID.ValueString())
		} else if resp.StatusCode != common.NormalStatusCode {
			return fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		}
	}

	// 确认是否更新成功
	err := c.rwLoop(ctx, *plan)
	if err != nil {
		return err
	}
	return nil
}

func (c *ctyunSfs) rwLoop(ctx context.Context, plan CtyunSfsConfig) error {
	var err error
	retryer, err := business.NewRetryer(time.Second*10, 60)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			var rwResp *sfs.SfsSfsListReadWriteSfs1Response
			rwResp, err = c.getSfsRwDetail(ctx, &plan)
			if err != nil {
				return false
			}
			readOnly := rwResp.ReturnObj.List[0].ReadOnly.Value
			//if readOnly != fmt.Sprintf("%t", plan.ReadOnly.ValueBool()) {
			if *readOnly != plan.ReadOnly.ValueBool() {
				return true
			}
			return false
		})
	if err != nil {
		return err
	}
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，弹性文件系统名称未修改成功！")
	}
	return err
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionId],[projectId]
func (c *ctyunSfs) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunSfsConfig
	var ID, regionId, projectId string
	err = terraform_extend.Split(request.ID, &ID, &regionId, &projectId)
	if err != nil {
		return
	}

	cfg.ID = types.StringValue(ID)
	cfg.RegionID = types.StringValue(regionId)
	cfg.ProjectID = types.StringValue(projectId)

	err = c.getAndMergeSfs(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

type CtyunSfsConfig struct {
	RegionID    types.String `tfsdk:"region_id"`
	IsEncrypt   types.Bool   `tfsdk:"is_encrypt"`
	KmsUUID     types.String `tfsdk:"kms_uuid"`
	ProjectID   types.String `tfsdk:"project_id"`
	SfsType     types.String `tfsdk:"sfs_type"`
	SfsProtocol types.String `tfsdk:"sfs_protocol"`
	Name        types.String `tfsdk:"name"`
	SfsSize     types.Int32  `tfsdk:"sfs_size"`
	CycleType   types.String `tfsdk:"cycle_type"`
	CycleCount  types.Int64  `tfsdk:"cycle_count"`
	AzName      types.String `tfsdk:"az_name"`
	VpcID       types.String `tfsdk:"vpc_id"`
	SubnetID    types.String `tfsdk:"subnet_id"`
	ID          types.String `tfsdk:"id"`
	Status      types.String `tfsdk:"status"`
	UsedSize    types.Int32  `tfsdk:"used_size"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
}

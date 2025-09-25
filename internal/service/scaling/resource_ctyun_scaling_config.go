package scaling

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

type ctyunScalingConfig struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
	imageService  *business.ImageService
}

func (c *ctyunScalingConfig) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scaling_config"
}

func (c *ctyunScalingConfig) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)
	c.imageService = business.NewImageService(c.meta)
}

func NewCtyunScalingConfig() resource.Resource {
	return &ctyunScalingConfig{}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[regionId],[projectId]
func (c *ctyunScalingConfig) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var cfg CtyunScalingConfigModel
	var ID, regionId string
	err = terraform_extend.Split(request.ID, &ID, &regionId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	cfg.ID = types.Int64Value(id)
	cfg.RegionID = types.StringValue(regionId)

	err = c.getAndMergeScalingConfig(ctx, &cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunScalingConfig) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027725/10241446**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "伸缩配置名称,长度为 2～15 个字符，允许使用大小写字母、数字或连字符（-）。不能以点号（.）或连字符（-）开头或结尾，不能连续使用点号（.）或连字符（-），也不能仅使用数字，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 15),
					validator2.ScalingConfigNameValidate(),
				},
			},
			"image_id": schema.StringAttribute{
				Required:    true,
				Description: "镜像ID，可以通过data.ctyun_images(datasource)获取，支持更新",
			},
			//"security_group_id_list": schema.SetAttribute{
			//	ElementType: types.StringType,
			//	Optional:    true,
			//	Description: "安全组ID列表，非多可用区资源池该参数为必填",
			//},
			// todo 描述清楚如何获取
			"flavor_name": schema.StringAttribute{
				Required:    true,
				Description: "规格名称，形如c7.2xlarge.4，支持更新。",
			},
			"volumes": schema.ListNestedAttribute{
				Required:    true,
				Description: "磁盘类型和大小列表，最多添加9块硬盘。系统盘仅支持1块。数据盘最多支持8块，支持更新。",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"volume_type": schema.StringAttribute{
							Required:    true,
							Description: "磁盘类型: SATA/SAS/SSD/SATA-KUNPENG/SATA-HAIGUANG/SAS-KUNPENG/SAS-HAIGUANG/SSD-genric，支持更新",
						},
						"volume_size": schema.Int32Attribute{
							Required:    true,
							Description: "磁盘大小(GB)，支持更新",
						},
						"disk_mode": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "磁盘模式: VBD（虚拟块存储设备）/ISCSI（小型计算机系统接口）。当flag=OS情况下，不可填写。数据盘磁盘模式，默认为VBD，支持更新",
							Validators: []validator.String{
								validator2.ConflictsWithEqualString(
									path.MatchRoot("flag"),
									types.StringValue(business.ScalingVolumeFlagOSStr),
								),
							},
						},
						"flag": schema.StringAttribute{
							Required:    true,
							Description: "磁盘类型: OS-系统盘, DATA-数据盘，系统盘限制1块。，支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf(business.ScalingVolumeFlag...),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.UniqueValues(),
				},
			},
			"use_floatings": schema.StringAttribute{
				Required:    true,
				Description: "是否使用弹性IP: diable-不使用, auto-自动分配。支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ScalingUseFloatings...),
				},
			},
			"bandwidth": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "弹性IP带宽(Mbps)，范围1-3000，支持更新",
				Validators: []validator.Int32{
					int32validator.Between(1, 3000),
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("use_floatings"),
						types.StringValue(business.ScalingUseFloatingsAutoStr),
					),
					validator2.ConflictsWithEqualInt32(
						path.MatchRoot("use_floatings"),
						types.StringValue(business.ScalingUseFloatingsDisableStr),
					),
				},
			},
			"login_mode": schema.StringAttribute{
				Required:    true,
				Description: "登录方式: password-密码, key_pair-密钥对，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ScalingLoginMode...),
				},
			},
			"username": schema.StringAttribute{
				Computed:    true,
				Description: "用户名，windows系统为administrator,linux系统为root。不可修改",
				//PlanModifiers: []planmodifier.String{
				//	stringplanmodifier.UseStateForUnknown(),
				//},
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "密码，login_mode为password时必填。密码规则：（1）8～30 个字符（2）必须同时包含三项（大写字母、小写字母、数字、 ()`~!@#$%^&*_-+=|{}[]:;'<>,.?/ 中的特殊符号）（3）不能以斜线号（/）开头 （4）不能包含3个及以上连续字符，如abc、123 （5）Windows镜像不能包含镜像用户名（Administrator）、用户名大小写变化（adminiSTrator），支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 30),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("login_mode"),
						types.StringValue(business.ScalingLoginModePasswordStr),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("login_mode"),
						types.StringValue(business.ScalingLoginModeKeyPairStr),
					),
					validator2.ScalingConfigPasswordValidate(),
				},
			},
			"key_pair_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "密钥对ID，login_mode为key_pair时必填，支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("login_mode"),
						types.StringValue(business.ScalingLoginModeKeyPairStr),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("login_mode"),
						types.StringValue(business.ScalingLoginModePasswordStr),
					),
				},
			},
			"tags": schema.ListNestedAttribute{
				Optional:    true,
				Description: "标签集，支持更新",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:    true,
							Description: "标签键，支持更新",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "标签值，支持更新",
						},
					},
				},
			},
			"az_names": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "可用区列表，仅多可用区资源池支持，支持更新",
			},
			"monitor_service": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否开启详细监控，支持更新",
				Default:     booldefault.StaticBool(true),
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "弹性伸缩配置ID",
			},
		},
	}
}

func (c *ctyunScalingConfig) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunScalingConfigModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	//创建前检查,检查证书有效性
	isValid, err := c.checkBeforeScalingConfig(ctx, plan)
	if !isValid || err != nil {
		return
	}
	err = c.createScalingConfig(ctx, &plan)
	if err != nil {
		return
	}
	// 创建后，通过创建的请求轮询，确认创建成功
	if err != nil {
		return
	}
	// 创建后反查创建后的证书信息
	err = c.getAndMergeScalingConfig(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunScalingConfig) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunScalingConfigModel
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeScalingConfig(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "未找到") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunScalingConfig) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 读取 plan -tf文件中配置
	var plan CtyunScalingConfigModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunScalingConfigModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
	}

	// 更新基本信息
	err = c.updateScalingConfig(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端数据，并同步本地state
	err = c.getAndMergeScalingConfig(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunScalingConfig) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunScalingConfigModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &scaling.ScalingConfigDeleteRequest{
		RegionID: state.RegionID.ValueString(),
		ConfigID: state.ID.ValueInt64(),
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingConfigDeleteApi.Do(ctx, c.meta.SdkCredential, params)
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
	err = c.checkDelete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *ctyunScalingConfig) checkBeforeScalingConfig(ctx context.Context, plan CtyunScalingConfigModel) (bool, error) {
	return true, nil
}

func (c *ctyunScalingConfig) createScalingConfig(ctx context.Context, config *CtyunScalingConfigModel) error {
	params := &scaling.ScalingConfigCreateRequest{
		RegionID:       config.RegionID.ValueString(),
		Name:           config.Name.ValueString(),
		ImageID:        config.ImageID.ValueString(),
		SpecName:       config.FlavorName.ValueString(),
		UseFloatings:   business.ScalingUseFloatingsDict[config.UseFloatings.ValueString()],
		LoginMode:      business.ScalingLoginModeDict[config.LoginMode.ValueString()],
		MonitorService: config.MonitorService.ValueBoolPointer(),
	}

	// 判断资源池是否为多AZ
	zones, err2 := c.regionService.GetZonesByRegionID(ctx, config.RegionID.ValueString())
	if err2 != nil {
		return err2
	}
	isNaz := false
	if len(zones) > 1 {
		isNaz = true
	}
	//// 当AZ必填
	//if !config.SecurityGroupIDList.IsNull() && !config.SecurityGroupIDList.IsUnknown() && !isNaz {
	//	var securityGroupIDList []string
	//	diags := config.SecurityGroupIDList.ElementsAs(ctx, &securityGroupIDList, true)
	//	if diags.HasError() {
	//		err := errors.New(diags[0].Detail())
	//		return err
	//	}
	//	params.SecurityGroupIDList = securityGroupIDList
	//}

	if !config.AzNames.IsNull() && !config.AzNames.IsUnknown() && isNaz {
		var azNames []string
		diags := config.AzNames.ElementsAs(ctx, &azNames, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		params.AzNames = azNames
	}
	if !config.Volumes.IsNull() && !config.Volumes.IsUnknown() {
		var volumeList []CtyunVolumesModel
		var volumeListReq []*scaling.ScalingConfigCreateVolumesRequest
		diags := config.Volumes.ElementsAs(ctx, &volumeList, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		//todo 极速型SSD云硬盘仅支持挂载至vCPU数量至少为16且为6代以上的通用计算增强型和内存优化型云主机
		for _, volumeItem := range volumeList {
			var volume scaling.ScalingConfigCreateVolumesRequest
			volume.VolumeSize = volumeItem.VolumeSize.ValueInt32()
			volume.VolumeType = volumeItem.VolumeType.ValueString()
			volume.Flag = business.ScalingVolumeFlagDict[volumeItem.Flag.ValueString()]
			if !volumeItem.DiskMode.IsNull() && !volumeItem.DiskMode.IsUnknown() {
				volume.DiskMode = volumeItem.DiskMode.ValueString()
			}
			volumeListReq = append(volumeListReq, &volume)
		}
		params.Volumes = volumeListReq
	}

	if config.UseFloatings.ValueString() == business.ScalingUseFloatingsAutoStr && !config.BandWidth.IsNull() {
		params.BandWidth = config.BandWidth.ValueInt32()
	}
	// 当登录类型为password，通过image确定username
	// os_type = linux, username = root
	// os_type = windows, username = administrator
	if config.LoginMode.ValueString() == business.ScalingLoginModePasswordStr {
		imageInfo, err := c.imageService.GetImageInfo(ctx, config.ImageID.ValueString(), config.RegionID.ValueString())
		if err != nil {
			return err
		}
		if imageInfo.OsType == "linux" {
			config.Username = types.StringValue("root")
		} else if imageInfo.OsType == "windows" {
			config.Username = types.StringValue("administrator")
		}
		params.Username = config.Username.ValueString()

		//if !config.Username.IsNull() {
		//
		//	params.Username = config.Username.ValueString()
		//} else {
		//	err := errors.New("当login_mode取值范围为password。username必填")
		//	return err
		//}
		if !config.Password.IsNull() {
			params.Password = config.Password.ValueString()
			var result bool
			result, err = c.checkPassword(params.Password, config.Username)
			if err != nil {
				return err
			}
			if !result {
				return err
			}
		} else {
			err := errors.New("当login_mode取值范围为password。password必填")
			return err
		}
	} else if config.LoginMode.ValueString() == business.ScalingLoginModeKeyPairStr {
		if !config.KeyPairID.IsNull() {
			params.KeyPairID = config.KeyPairID.ValueString()
		} else {
			err := errors.New("当login_mode取值范围为key_pair。key_pair_id必填")
			return err
		}
	} else {
		err := errors.New("login_mode输入有误，取值范围：password和key_pair")
		return err
	}

	if !config.Tags.IsNull() {
		var tags []CtyunTagModel
		var tagList []*scaling.ScalingConfigCreateTagsRequest
		diags := config.Tags.ElementsAs(ctx, &tags, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		for _, tagItem := range tags {
			var tag scaling.ScalingConfigCreateTagsRequest
			tag.Key = tagItem.Key.ValueString()
			tag.Value = tagItem.Value.ValueString()
			tagList = append(tagList, &tag)
		}
		params.Tags = tagList
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingConfigCreateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = errors.New("创建伸缩配置失败，接口返回nil。请稍后重试！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}

	// 保存弹性伸缩配置id
	config.ID = types.Int64Value(resp.ReturnObj.Id)
	return nil
}

func (c *ctyunScalingConfig) getAndMergeScalingConfig(ctx context.Context, config *CtyunScalingConfigModel) error {
	// 获取伸缩配置详情
	scalingConfigDetailResp, err := c.getScalingDetail(ctx, config)
	if err != nil {
		return err
	}
	var diags diag.Diagnostics

	// 同步云端值
	detail := scalingConfigDetailResp.ReturnObj[0]
	config.Name = types.StringValue(detail.Name)
	config.ImageID = types.StringValue(detail.ImageID)
	config.FlavorName = types.StringValue(detail.SpecName)
	config.UseFloatings = types.StringValue(business.ScalingUseFloatingsDictRev[detail.UseFloatings])
	if config.UseFloatings.ValueString() == business.ScalingUseFloatingsAutoStr {
		config.BandWidth = types.Int32Value(detail.Bandwidth)
	}
	config.LoginMode = types.StringValue(business.ScalingLoginModeDictRev[detail.LoginMode])
	config.MonitorService = types.BoolValue(*detail.MonitorService)
	var volumeList []CtyunVolumesModel
	for _, volumeItem := range detail.Volumes {
		var volume CtyunVolumesModel
		volume.VolumeSize = types.Int32Value(volumeItem.VolumeSize)
		volume.VolumeType = types.StringValue(volumeItem.VolumeType)
		volume.DiskMode = types.StringValue(volumeItem.DiskMode)
		volume.Flag = types.StringValue(business.ScalingVolumeFlagDictRev[volumeItem.Flag])
		volumeList = append(volumeList, volume)
	}
	config.Volumes, diags = types.ListValueFrom(ctx, utils.StructToTFObjectTypes(CtyunVolumesModel{}), volumeList)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return err
	}

	var tagList []CtyunTagModel
	for _, tagItem := range detail.Tags {
		var tag CtyunTagModel
		tag.Key = types.StringValue(tagItem.Key)
		tag.Value = types.StringValue(tagItem.Value)
		tagList = append(tagList, tag)
	}
	config.Tags, diags = types.ListValueFrom(ctx, utils.StructToTFObjectTypes(CtyunTagModel{}), tagList)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return err
	}
	azNames := strings.Split(detail.AzNames, ",")
	config.AzNames, diags = types.SetValueFrom(ctx, types.StringType, azNames)
	if diags.HasError() {
		err = errors.New(diags[0].Detail())
		return err
	}

	if config.LoginMode.ValueString() == business.ScalingLoginModePasswordStr {
		config.KeyPairID = types.StringValue("")
	} else if config.LoginMode.ValueString() == business.ScalingLoginModeKeyPairStr {
		config.Password = types.StringValue("")
		config.Username = types.StringValue("")
	}
	if config.UseFloatings.ValueString() == business.ScalingUseFloatingsDisableStr {
		config.BandWidth = types.Int32Value(0)
	}
	return nil

}

func (c *ctyunScalingConfig) getScalingDetail(ctx context.Context, config *CtyunScalingConfigModel) (*scaling.ScalingConfigListResponse, error) {
	params := &scaling.ScalingConfigListRequest{
		RegionID: config.RegionID.ValueString(),
		ConfigID: config.ID.ValueInt64(),
		PageNo:   1,
		PageSize: 10,
	}
	resp, err := c.meta.Apis.SdkScalingApis.ScalingConfigListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取id为 %d 的伸缩配置详情失败，接口返回nil。请稍后重试！", config.ID.ValueInt64())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	if len(resp.ReturnObj) > 1 {
		err = fmt.Errorf("根据id：%d查询绳索配置详细信息，返回条数大于1。具体如下:%#v\n", config.ID.ValueInt64(), resp.ReturnObj)
		return resp, err
	}
	if len(resp.ReturnObj) < 1 {
		err = fmt.Errorf("根据id：%d查询绳索配置详细信息，返回条数为0", config.ID.ValueInt64())
		return nil, err
	}
	return resp, nil
}

func (c *ctyunScalingConfig) updateScalingConfig(ctx context.Context, state *CtyunScalingConfigModel, plan *CtyunScalingConfigModel) error {
	params := &scaling.ScalingConfigUpdateRequest{
		RegionID:       state.RegionID.ValueString(),
		ConfigID:       state.ID.ValueInt64(),
		Name:           plan.Name.ValueString(),
		ImageID:        plan.ImageID.ValueString(),
		SpecName:       plan.FlavorName.ValueString(),
		UseFloatings:   business.ScalingUseFloatingsDict[plan.UseFloatings.ValueString()],
		LoginMode:      business.ScalingLoginModeDict[plan.LoginMode.ValueString()],
		Tags:           nil,
		AzNames:        nil,
		MonitorService: plan.MonitorService.ValueBoolPointer(),
	}
	if !plan.Name.Equal(state.Name) {
		params.Name = plan.Name.ValueString()
	}
	if !plan.ImageID.Equal(state.ImageID) {
		params.ImageID = plan.ImageID.ValueString()
	}
	if !plan.FlavorName.Equal(state.FlavorName) {
		params.SpecName = plan.FlavorName.ValueString()
	}
	if params.UseFloatings == business.ScalingUseFloatingsAuto {
		params.BandWidth = plan.BandWidth.ValueInt32()
	}
	if params.LoginMode == business.ScalingLoginModePassword {
		// username 根据image_id获取
		imageInfo, err := c.imageService.GetImageInfo(ctx, state.ImageID.ValueString(), state.RegionID.ValueString())
		if err != nil {
			return err
		}
		if imageInfo.OsType == "linux" {
			state.Username = types.StringValue("root")
		} else if imageInfo.OsType == "windows" {
			state.Username = types.StringValue("administrator")
		}
		params.Username = state.Username.ValueString()

		params.Password = plan.Password.ValueString()
	} else if params.LoginMode == business.ScalingLoginModeKeyPair {
		params.KeyPairID = plan.KeyPairID.ValueString()
		state.KeyPairID = plan.KeyPairID
	}
	if !plan.Volumes.IsNull() && !plan.Volumes.IsUnknown() {
		var volumeList []CtyunVolumesModel
		var volumeListReq []*scaling.ScalingConfigUpdateVolumesRequest
		diags := plan.Volumes.ElementsAs(ctx, &volumeList, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		//todo 极速型SSD云硬盘仅支持挂载至vCPU数量至少为16且为6代以上的通用计算增强型和内存优化型云主机
		for _, volumeItem := range volumeList {
			var volume scaling.ScalingConfigUpdateVolumesRequest
			volume.VolumeSize = volumeItem.VolumeSize.ValueInt32()
			volume.VolumeType = volumeItem.VolumeType.ValueString()
			volume.DiskMode = volumeItem.DiskMode.ValueString()
			volume.Flag = business.ScalingVolumeFlagDict[volumeItem.Flag.ValueString()]
			volumeListReq = append(volumeListReq, &volume)
		}
		params.Volumes = volumeListReq
	}

	if !plan.Tags.IsNull() {
		var tags []CtyunTagModel
		var tagList []*scaling.ScalingConfigUpdateTagsRequest
		diags := plan.Tags.ElementsAs(ctx, &tags, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		for _, tagItem := range tags {
			var tag scaling.ScalingConfigUpdateTagsRequest
			tag.Key = tagItem.Key.ValueString()
			tag.Value = tagItem.Value.ValueString()
			tagList = append(tagList, &tag)
		}
		params.Tags = tagList
	}

	if !plan.AzNames.IsNull() && !plan.AzNames.IsUnknown() {
		var azNames []string
		diags := plan.AzNames.ElementsAs(ctx, &azNames, true)
		if diags.HasError() {
			err := errors.New(diags[0].Detail())
			return err
		}
		params.AzNames = azNames
	}

	resp, err := c.meta.Apis.SdkScalingApis.ScalingConfigUpdateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新弹性伸缩配置失败，接口返回nil。请联系研发确认！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return err
	}
	return nil
}

func (c *ctyunScalingConfig) checkDelete(ctx context.Context, state *CtyunScalingConfigModel) error {
	// 查询列表，是否还能查询到
	resp, err := c.getScalingDetail(ctx, state)
	if err != nil {
		return nil
	}
	if len(resp.ReturnObj) == 1 {
		return errors.New("弹性伸缩配置尚未被删除")
	}
	return nil
}

func (c *ctyunScalingConfig) checkPassword(password string, username types.String) (bool, error) {
	lowerPass := strings.ToLower(password)
	if strings.Contains(lowerPass, username.ValueString()) {
		err := errors.New("密码不能包含管理员用户名")
		return false, err
	}
	return true, nil
}

type CtyunScalingConfigModel struct {
	RegionID types.String `tfsdk:"region_id"` // 资源池ID
	Name     types.String `tfsdk:"name"`      // 伸缩配置名称
	ImageID  types.String `tfsdk:"image_id"`  // 镜像ID
	//SecurityGroupIDList types.Set    `tfsdk:"security_group_id_list"` // 安全组ID列表
	FlavorName     types.String `tfsdk:"flavor_name"`     // 规格名称
	Volumes        types.List   `tfsdk:"volumes"`         // 磁盘类型和大小列表
	UseFloatings   types.String `tfsdk:"use_floatings"`   // 是否使用弹性IP
	BandWidth      types.Int32  `tfsdk:"bandwidth"`       // 弹性IP带宽
	LoginMode      types.String `tfsdk:"login_mode"`      // 登录方式
	Username       types.String `tfsdk:"username"`        // 用户名
	Password       types.String `tfsdk:"password"`        // 密码
	KeyPairID      types.String `tfsdk:"key_pair_id"`     // 密钥对ID
	Tags           types.List   `tfsdk:"tags"`            // 标签集
	AzNames        types.Set    `tfsdk:"az_names"`        // 可用区列表
	MonitorService types.Bool   `tfsdk:"monitor_service"` // 是否开启详细监控
	ID             types.Int64  `tfsdk:"id"`              // 伸缩配置ID
}

type CtyunVolumesModel struct {
	VolumeType types.String `tfsdk:"volume_type"` // 磁盘类型
	VolumeSize types.Int32  `tfsdk:"volume_size"` // 磁盘大小(GB)
	DiskMode   types.String `tfsdk:"disk_mode"`   // 磁盘模式
	Flag       types.String `tfsdk:"flag"`        // 磁盘类型(1-系统盘,2-数据盘)
}

type CtyunTagModel struct {
	Key   types.String `tfsdk:"key"`   // 标签键
	Value types.String `tfsdk:"value"` // 标签值
}

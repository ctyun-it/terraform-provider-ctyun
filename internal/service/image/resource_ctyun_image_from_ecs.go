package image

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctimage"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"log"
	"regexp"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunImageFromEcs{}
	_ resource.ResourceWithConfigure   = &ctyunImageFromEcs{}
	_ resource.ResourceWithImportState = &ctyunImageFromEcs{}
)

func NewCtyunImageFromEcs() resource.Resource {
	return &ctyunImageFromEcs{}
}

type ctyunImageFromEcs struct {
	meta *common.CtyunMetadata
}

func (c *ctyunImageFromEcs) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_image_from_ecs"
}

func (c *ctyunImageFromEcs) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10027726/10031013`,
		Attributes: map[string]schema.Attribute{
			// 新增：资源ID（由API返回，自动生成）
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "私有镜像唯一标识（镜像ID）",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// 公共必填参数（所有创建方式均需）
			"image_name": schema.StringAttribute{
				Required:    true,
				Description: "镜像名称。长度2~32字符，仅数字、字母、-组成，不能以数字或-开头/结尾，且不与已有私有镜像重名。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 32),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$`),
						"镜像名称格式错误：需2~32字符，仅数字、字母、-组成，不以数字或-开头/结尾。",
					),
				},
			},
			// 镜像类型字段
			"image_type": schema.StringAttribute{
				Required:    true,
				Description: "镜像类型，可选值：system_disk（系统盘镜像）、data_disk（数据盘镜像）、entire_machine（整机镜像）",
				Validators: []validator.String{
					stringvalidator.OneOf("system_disk", "data_disk", "entire_machine"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// 公共可选参数（所有创建方式均支持）
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "镜像描述。长度1~128字符，不能以空格开头或结尾。 支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^\\s].*[^\\s]$|^[^\\s]$`),
						"描述不能以空格开头或结尾。",
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
				Validators: []validator.String{
					validator2.Project(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"labels": schema.ListNestedAttribute{
				Optional:    true,
				Description: "标签列表。最多10个标签，标签键不可重复，键值长度1~32字符，不能换行或以空格开头/结尾。",
				Validators: []validator.List{
					listvalidator.SizeAtMost(10),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"label_key": schema.StringAttribute{
							Required:    true,
							Description: "标签键。",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
								//stringvalidator.RegexMatches(
								//	regexp.MustCompile(`^[^\\s\\n].*[^\\s\\n]$|^[^\\s\\n]$`),
								//	"标签键不能换行或以空格开头/结尾。",
								//),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"label_value": schema.StringAttribute{
							Required:    true,
							Description: "标签值。",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
								//stringvalidator.RegexMatches(
								//	regexp.MustCompile(`^[^\\s\\n].*[^\\s\\n]$|^[^\\s\\n]$`),
								//	"标签值不能换行或以空格开头/结尾。",
								//),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
				// TODO 标签变更需重建（暂无无动态更新标签能力）
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			//
			"enable_image_integrity_check": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "是否启用镜像完整性校验，仅资源池支持时生效。",
				// 云主机变更需重建
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},

			// 系统盘/数据盘/整机共用参数（云主机ID）
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "云主机ID，系统盘、数据盘、整机创建方式必填。",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
				// 云主机变更需重建
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// 系统盘特有参数
			"minimum_ram": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				//Default:     int64default.StaticInt64(0),
				Description: "最小内存限制（GiB），仅系统盘、快照创建方式支持，取值0/1/2/4/8/16/32/64/128/256/512。支持更新",
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("image_type"),
						types.StringValue("system_disk"),
					),
					int64validator.OneOf(0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512),
				},
			},
			"maximum_ram": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				//Default:     int64default.StaticInt64(0),
				Description: "最大内存限制（GiB），仅系统盘、快照创建方式支持，需≥最小内存。 支持更新",
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("image_type"),
						types.StringValue("system_disk"),
					),
					int64validator.OneOf(0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512),
				},
			},

			// 数据盘特有参数
			"data_disk_id": schema.StringAttribute{
				Optional:    true,
				Description: "数据盘ID，仅数据盘创建方式必填，需挂载于指定云主机。",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
				// 数据盘变更需重建
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// 整机特有参数
			"repository_id": schema.StringAttribute{
				Optional:    true,
				Description: "云主机备份存储库ID，仅整机创建方式在非多可用区资源池时必填。",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
				// 存储库变更需重建
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
func (c *ctyunImageFromEcs) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunImageFromEcsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 根据传入的不同参数确定创建方式
	imageType := plan.ImageType.ValueString()

	switch imageType {
	case "system_disk":
		err = c.createSystemDiskImage(ctx, &plan, response)
	case "data_disk":
		err = c.createDataDiskImage(ctx, &plan, response)
	case "entire_machine":
		err = c.createEntireMachineImage(ctx, &plan, response)
	}
	if err != nil {
		return
	}
	// 查询镜像状态信息
	err = c.getAndMergeImage(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (c *ctyunImageFromEcs) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunImageFromEcsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.getAndMergeImage(ctx, &state)
	if err != nil {
		if err.Error() == common.ImageImageCheckNotFound {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunImageFromEcs) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// tf文件中的计划状态
	var plan CtyunImageFromEcsConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// state中的当前状态
	var state CtyunImageFromEcsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 检查更新前条件
	err = c.checkBeforeUpdate(ctx, &plan, &state)
	if err != nil {
		return
	}

	// 按字段逐个更新
	err = c.updateImageAttributes(ctx, &plan, &state)
	if err != nil {
		return
	}

	// 查询远端信息确保更新成功，并将最新状态设置到state中
	err = c.getAndMergeImage(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunImageFromEcs) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunImageFromEcsConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.SdkCtImageApis.CtimageDeleteImageApi.Do(ctx, c.meta.SdkCredential, &ctimage.CtimageDeleteImageRequest{
		ImageID:  state.Id.ValueString(),
		RegionID: state.RegionId.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 等待镜像删除成功
	e := c.waitForImageDeleted(ctx, &state)
	if e != nil {
		response.Diagnostics.AddError(e.Error(), e.Error())
		return
	}
}

func (c *ctyunImageFromEcs) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [imageId],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunImageFromEcsConfig
	var imageId, regionId string
	// 根据分隔符数量判断是否输入了regionId
	if strings.Count(request.ID, common.ImportSeparator) < 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &imageId)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &imageId, &regionId)
		if err != nil {
			return
		}
	}

	if imageId == "" {
		err = fmt.Errorf("imageId不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}

	cfg.Id = types.StringValue(imageId)
	cfg.RegionId = types.StringValue(regionId)

	err = c.getAndMergeImage(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &cfg)...)
}

func (c *ctyunImageFromEcs) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
func (c *ctyunImageFromEcs) createSystemDiskImage(ctx context.Context, plan *CtyunImageFromEcsConfig, response *resource.CreateResponse) (err error) {
	// 系统盘创建方式
	// 构造标签列表
	// 根据传入的不同参数确定创建方式
	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	var labels []*ctimage.CtimageCreateEcsSystemDiskImageLabelsRequest
	if plan.Labels != nil {
		labels = make([]*ctimage.CtimageCreateEcsSystemDiskImageLabelsRequest, len(plan.Labels))
		for i, label := range plan.Labels {
			labels[i] = &ctimage.CtimageCreateEcsSystemDiskImageLabelsRequest{
				LabelKey:   label.LabelKey.ValueString(),
				LabelValue: label.LabelValue.ValueString(),
			}
		}
	}

	// 系统盘创建方式
	systemReq := &ctimage.CtimageCreateEcsSystemDiskImageRequest{
		ImageName:                 plan.ImageName.ValueString(),
		RegionID:                  regionId,
		InstanceID:                plan.InstanceId.ValueString(),
		Description:               plan.Description.ValueString(),
		EnableImageIntegrityCheck: plan.EnableImageIntegrityCheck.ValueBoolPointer(),
		Labels:                    labels,
		MaximumRAM:                int32(plan.MaximumRAM.ValueInt64()),
		MinimumRAM:                int32(plan.MinimumRAM.ValueInt64()),
		ProjectID:                 projectId,
	}

	resp, err := c.meta.Apis.SdkCtImageApis.CtimageCreateEcsSystemDiskImageApi.Do(ctx, c.meta.SdkCredential, systemReq)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	if resp.ReturnObj == nil || len(resp.ReturnObj.Images) == 0 {
		log.Printf("[DEBUG] 创建镜像失败，API返回: %+v", resp)
		response.Diagnostics.AddError("创建镜像失败", fmt.Sprintf("API返回数据为空，完整响应: %+v", resp))
		return
	}

	image := resp.ReturnObj.Images[0]
	plan.ImageName = types.StringValue(image.ImageName)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	plan.Id = types.StringValue(image.ImageID)
	plan.Description = types.StringValue(image.Description)

	// 设置其他属性
	if image.EnableImageIntegrityCheck != nil {
		plan.EnableImageIntegrityCheck = types.BoolPointerValue(image.EnableImageIntegrityCheck)
	}

	plan.MaximumRAM = types.Int64Value(int64(image.MaximumRAM))
	plan.MinimumRAM = types.Int64Value(int64(image.MinimumRAM))

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 轮询镜像状态直到active
	return c.waitForUploadImageActive(ctx, plan)

}
func (c *ctyunImageFromEcs) createDataDiskImage(ctx context.Context, plan *CtyunImageFromEcsConfig, response *resource.CreateResponse) (err error) {
	// 构造标签列表
	// 根据传入的不同参数确定创建方式
	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	// 数据盘创建方式
	// 构造标签列表
	var labels []*ctimage.CtimageCreateEcsDataDiskImageLabelsRequest
	if plan.Labels != nil {
		labels = make([]*ctimage.CtimageCreateEcsDataDiskImageLabelsRequest, len(plan.Labels))
		for i, label := range plan.Labels {
			labels[i] = &ctimage.CtimageCreateEcsDataDiskImageLabelsRequest{
				LabelKey:   label.LabelKey.ValueString(),
				LabelValue: label.LabelValue.ValueString(),
			}
		}
	}

	// 数据盘创建方式
	dataDiskReq := &ctimage.CtimageCreateEcsDataDiskImageRequest{
		ImageName:                 plan.ImageName.ValueString(),
		RegionID:                  regionId,
		InstanceID:                plan.InstanceId.ValueString(),
		DataDiskID:                plan.DataDiskId.ValueString(),
		Description:               plan.Description.ValueString(),
		EnableImageIntegrityCheck: plan.EnableImageIntegrityCheck.ValueBoolPointer(),
		Labels:                    labels,
		ProjectID:                 projectId,
	}

	resp, err := c.meta.Apis.SdkCtImageApis.CtimageCreateEcsDataDiskImageApi.Do(ctx, c.meta.SdkCredential, dataDiskReq)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	if resp.ReturnObj == nil || len(resp.ReturnObj.Images) == 0 {
		log.Printf("[DEBUG] 创建镜像失败，API返回: %+v", resp)
		response.Diagnostics.AddError("创建镜像失败", fmt.Sprintf("API返回数据为空，完整响应: %+v", resp))
		return
	}
	image := resp.ReturnObj.Images[0]
	plan.ImageName = types.StringValue(image.ImageName)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	plan.Id = types.StringValue(image.ImageID)
	plan.Description = types.StringValue(image.Description)

	// 设置其他属性
	if image.EnableImageIntegrityCheck != nil {
		plan.EnableImageIntegrityCheck = types.BoolPointerValue(image.EnableImageIntegrityCheck)
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 轮询镜像状态直到active
	return c.waitForUploadImageActive(ctx, plan)

}
func (c *ctyunImageFromEcs) createEntireMachineImage(ctx context.Context, plan *CtyunImageFromEcsConfig, response *resource.CreateResponse) (err error) {

	// 构造标签列表
	// 根据传入的不同参数确定创建方式
	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	// 整机创建方式
	// 构造标签列表
	var labels []*ctimage.CtimageCreateFullEcsImageLabelsRequest
	if plan.Labels != nil {
		labels = make([]*ctimage.CtimageCreateFullEcsImageLabelsRequest, len(plan.Labels))
		for i, label := range plan.Labels {
			labels[i] = &ctimage.CtimageCreateFullEcsImageLabelsRequest{
				LabelKey:   label.LabelKey.ValueString(),
				LabelValue: label.LabelValue.ValueString(),
			}
		}
	}

	// 整机创建方式
	entireReq := &ctimage.CtimageCreateFullEcsImageRequest{
		ImageName:   plan.ImageName.ValueString(),
		RegionID:    regionId,
		InstanceID:  plan.InstanceId.ValueString(),
		Description: plan.Description.ValueString(),
		Labels:      labels,
		ProjectID:   projectId,
	}
	// repository_id在非多可用区资源池时必填
	if !plan.RepositoryId.IsNull() {
		entireReq.RepositoryID = plan.RepositoryId.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtImageApis.CtimageCreateFullEcsImageApi.Do(ctx, c.meta.SdkCredential, entireReq)
	if err != nil {
		return
	}

	if resp.ReturnObj == nil || len(resp.ReturnObj.Images) == 0 {
		log.Printf("[DEBUG] 创建镜像失败，API返回: %+v", resp)
		response.Diagnostics.AddError("创建镜像失败", fmt.Sprintf("API返回数据为空，完整响应: %+v", resp))
		return
	}

	image := resp.ReturnObj.Images[0]
	plan.ImageName = types.StringValue(image.ImageName)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	plan.Id = types.StringValue(image.ImageID)
	plan.Description = types.StringValue(image.Description)

	// 设置其他属性
	if image.EnableImageIntegrityCheck != nil {
		plan.EnableImageIntegrityCheck = types.BoolPointerValue(image.EnableImageIntegrityCheck)
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 轮询镜像状态直到active
	return c.waitForUploadImageActive(ctx, plan)
}

// waitForImageDeleted 等待镜像删除成功
func (c *ctyunImageFromEcs) waitForImageDeleted(ctx context.Context, cfg *CtyunImageFromEcsConfig) error {
	executeSuccessFlag := false
	retryer, _ := business.NewRetryer(time.Second*5, 60)
	retryer.Start(
		func(currentTime int) bool {
			response, err := c.meta.Apis.SdkCtImageApis.CtimageDetailImageApi.Do(ctx, c.meta.SdkCredential, &ctimage.CtimageDetailImageRequest{
				ImageID:  cfg.Id.ValueString(),
				RegionID: cfg.RegionId.ValueString(),
			})

			if err != nil {
				// 执行完成后，查询不到镜像会抛错，这个是正常的出口
				if err.Error() == common.ImageImageCheckNotFound {
					executeSuccessFlag = true
					return false
				}
				return true
			}
			log.Printf("response: %+v", response)
			// 执行完成后，可能查询不到镜像的信息了，这个也是正常出口
			if response.ReturnObj == nil || len(response.ReturnObj.Images) == 0 {
				executeSuccessFlag = true
				return false
			}

			// 其余的情况，需要按照镜像的状态来判断是否需要继续执行
			// 注意：这里需要根据实际API响应中的状态字段进行判断
			// 由于当前API响应结构中未看到状态字段，暂时留空
			// 可以根据实际需要添加状态判断逻辑
			return true
		},
	)

	if !executeSuccessFlag {
		return errors.New("轮询镜像" + cfg.Id.ValueString() + "状态失败")
	}
	return nil
}

func (c *ctyunImageFromEcs) waitForUploadImageActive(ctx context.Context, cfg *CtyunImageFromEcsConfig) error {
	executeSuccessFlag := false
	retryer, _ := business.NewRetryer(time.Second*10, 180)

	retryer.Start(
		func(currentTime int) bool {

			response, err := c.meta.Apis.SdkCtImageApis.CtimageDetailImageApi.Do(ctx, c.meta.SdkCredential, &ctimage.CtimageDetailImageRequest{
				ImageID:  cfg.Id.ValueString(),
				RegionID: cfg.RegionId.ValueString(),
			})

			if err != nil {
				return true
			}

			if response.ReturnObj == nil || len(response.ReturnObj.Images) != 1 {
				return true
			}
			if response.ReturnObj.Images[0].ImageStatus == business.ImageStatusActive {
				executeSuccessFlag = true
				return false
			}
			return true

		},
	)
	if !executeSuccessFlag {
		return errors.New("轮询镜像" + cfg.Id.ValueString() + "状态失败")
	}
	return nil
}

// checkBeforeUpdate 更新前检查
func (c *ctyunImageFromEcs) checkBeforeUpdate(ctx context.Context, plan, state *CtyunImageFromEcsConfig) (err error) {
	// 检查镜像状态是否允许更新
	instance, err := c.getImageByID(ctx, state)
	if err != nil {
		return err
	}

	if instance.ImageStatus != business.ImageStatusActive {
		return fmt.Errorf("镜像状态不是active，无法进行更新操作")
	}

	return nil
}

// updateImageAttributes 按字段更新镜像属性
func (c *ctyunImageFromEcs) updateImageAttributes(ctx context.Context, plan, state *CtyunImageFromEcsConfig) (err error) {
	// 更新镜像名称
	err = c.updateImage(ctx, plan, state)
	if err != nil {
		return
	}

	return
}

// updateImage 更新镜像描述
func (c *ctyunImageFromEcs) updateImage(ctx context.Context, plan, state *CtyunImageFromEcsConfig) (err error) {
	params := &ctimage.CtimageUpdateImageRequest{
		ImageID:  state.Id.ValueString(),
		RegionID: state.RegionId.ValueString(),
	}
	if !plan.Description.IsNull() {
		params.Description = plan.Description.ValueString()
	}

	if !plan.ImageName.IsNull() {
		params.ImageName = plan.ImageName.ValueString()
	}

	// 仅在系统盘镜像类型下更新内存限制
	if plan.ImageType.ValueString() == "system_disk" {
		if !plan.MaximumRAM.IsNull() {
			params.MaximumRAM = int32(plan.MaximumRAM.ValueInt64())
		}
		if !plan.MinimumRAM.IsNull() {
			params.MinimumRAM = int32(plan.MinimumRAM.ValueInt64())
		}
	}
	_, err = c.meta.Apis.SdkCtImageApis.CtimageUpdateImageApi.Do(ctx, c.meta.SdkCredential, params)

	if err != nil {
		return fmt.Errorf("更新镜像描述失败: %w", err)
	}

	return
}

// updateImageMemoryLimits 更新内存限制

// getImageByID 根据ID获取镜像详情
func (c *ctyunImageFromEcs) getImageByID(ctx context.Context, cfg *CtyunImageFromEcsConfig) (*ctimage.CtimageDetailImageReturnObjImagesResponse, error) {
	response, err := c.meta.Apis.SdkCtImageApis.CtimageDetailImageApi.Do(ctx, c.meta.SdkCredential, &ctimage.CtimageDetailImageRequest{
		ImageID:  cfg.Id.ValueString(),
		RegionID: cfg.RegionId.ValueString(),
	})

	if err != nil {
		return nil, err
	}

	if response.ReturnObj == nil || len(response.ReturnObj.Images) == 0 {
		return nil, fmt.Errorf(common.ImageImageCheckNotFound)
	}

	return response.ReturnObj.Images[0], nil
}

// getAndMergeImage 查询合并镜像
func (c *ctyunImageFromEcs) getAndMergeImage(ctx context.Context, cfg *CtyunImageFromEcsConfig) (err error) {

	resp, err := c.getImageByID(ctx, cfg)
	if err != nil {
		return
	}

	cfg.Id = types.StringValue(resp.ImageID)
	cfg.ImageName = types.StringValue(resp.ImageName)
	cfg.Description = types.StringValue(resp.Description)

	// 设置其他属性
	if resp.EnableImageIntegrityCheck != nil {
		cfg.EnableImageIntegrityCheck = types.BoolPointerValue(resp.EnableImageIntegrityCheck)
	}

	cfg.MaximumRAM = types.Int64Value(int64(resp.MaximumRAM))
	cfg.MinimumRAM = types.Int64Value(int64(resp.MinimumRAM))
	cfg.ProjectId = types.StringValue(resp.ProjectID)
	// 如果有标签信息，也需要设置
	// 注意：根据API文档，详情接口可能不返回标签信息，需要根据实际情况调整

	return
}

// *CtyunImageFromEcsConfig 映射从云主机/快照创建私有镜像的配置参数，适配四种创建方式：
// 1. 系统盘镜像（API 4765）
// 2. 数据盘镜像（API 5230）
// 3. 整机镜像（API 18058）
// 4. 快照镜像（API 18057）
// 同时兼容修改接口（API 5085）的可更新字段特性
type CtyunImageFromEcsConfig struct {
	// ID是资源的唯一标识符，由Terraform框架自动管理
	Id types.String `tfsdk:"id"`
	// 镜像名称（所有创建方式必填，修改接口支持更新）
	// 约束：2~32字符，仅数字、字母、-组成，不以数字或-开头/结尾，不与同资源池私有镜像重名
	ImageName types.String `tfsdk:"image_name"`
	// 镜像类型
	ImageType types.String `tfsdk:"image_type"`

	// 资源池ID（所有创建方式必填，不可跨池修改）
	// 约束：需为用户可见的资源池ID，修改时必须与镜像所属资源池一致
	RegionId types.String `tfsdk:"region_id"`

	// 镜像描述（可选，修改接口支持更新）
	// 约束：1~128字符，不能以空格开头/结尾
	Description types.String `tfsdk:"description"`

	// 企业项目ID（可选，默认0，修改接口支持更新但需重建资源）
	// 约束：默认值为"0"（default项目），需为用户有权限的企业项目ID
	ProjectId types.String `tfsdk:"project_id"`

	// 标签列表（可选，修改接口支持更新但需重建资源）
	// 约束：最多10个标签，标签键不可重复，键值长度1~32字符，不能换行或以空格开头/结尾
	Labels []Label `tfsdk:"labels"`

	// 云主机ID（系统盘/数据盘/整机镜像创建必填，快照镜像创建不可填）
	// 约束：云主机状态需为stopped（部分资源池支持running），与数据盘/整机创建强关联
	InstanceId types.String `tfsdk:"instance_id"`

	// 数据盘ID（仅数据盘镜像创建必填）
	// 约束：需挂载于指定InstanceId的云主机，磁盘模式不为FCSAN/ISCSI，状态为in-use且未加密
	DataDiskId types.String `tfsdk:"data_disk_id"`

	// 备份存储库ID（仅整机镜像创建时，非多可用区资源池必填）
	// 约束：需为未到期、未冻结且容量充足的云主机备份存储库
	RepositoryId types.String `tfsdk:"repository_id"`

	// 最小内存限制（GiB，仅系统盘/快照镜像支持，修改接口支持更新）
	// 约束：取值为0（不限制）、1、2、4、8、16、32、64、128、256、512，需≤MaximumRAM
	MinimumRAM types.Int64 `tfsdk:"minimum_ram"`

	// 最大内存限制（GiB，仅系统盘/快照镜像支持，修改接口支持更新）
	// 约束：取值同MinimumRAM，需≥MinimumRAM
	MaximumRAM types.Int64 `tfsdk:"maximum_ram"`

	// 是否启用镜像完整性校验（可选，默认false，修改接口支持更新）
	// 约束：仅部分资源池支持，默认值为false
	EnableImageIntegrityCheck types.Bool `tfsdk:"enable_image_integrity_check"`
}

// Label 映射标签键值对，用于镜像的标签管理
type Label struct {
	// 标签键（必填）
	// 约束：1~32字符，不能换行或以空格开头/结尾，同一镜像的标签键不可重复
	LabelKey types.String `tfsdk:"label_key"`

	// 标签值（必填）
	// 约束：1~32字符，不能换行或以空格开头/结尾
	LabelValue types.String `tfsdk:"label_value"`
}

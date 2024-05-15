package resource

import (
	"context"
	"errors"
	"github.com/ctyun-it/ctyun-sdk-go/ctyun-sdk-endpoint/ctimage"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-ctyun/internal/business"
	"terraform-provider-ctyun/internal/common"
	terraform_extend "terraform-provider-ctyun/internal/extend/terraform"
	defaults2 "terraform-provider-ctyun/internal/extend/terraform/defaults"
	"time"
)

func NewCtyunImage() resource.Resource {
	return &ctyunImage{}
}

type ctyunImage struct {
	meta *common.CtyunMetadata
}

func (c *ctyunImage) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_image"
}

func (c *ctyunImage) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10027726**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "id",
			},
			"file_source": schema.StringAttribute{
				Required:    true,
				Description: "镜像文件地址，格式应为{internetEndpoint}/{bucket}/{key}。可使用访问控制endpoint查询接口来查询外网访问endpoint，可使用获取桶列表接口来查询您拥有的桶的列表，可使用查看对象列表接口来查询存储桶内所有对象",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "镜像名称，长度为2-32个字符，只能由数字、字母、-组成，不能以数字、-开头，且不能以-结尾",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$"), "不满足镜像名称要求"),
				},
			},
			"os_distro": schema.StringAttribute{
				Required:    true,
				Description: "操作系统的发行版名称。注意：对于Windows系操作系统，应确保参数值是windows，否则视作Linux系操作系统；对于Linux系操作系统，参数值的取值应根据系统实际情况，建议参照cloud-init的配置文件（/etc/cloud/cloud.cfg）中system_info.distro的取值或以下取值：anolis、centos、ctyunos、debian、fedora、kylin、openEuler、ubuntu、UnionTech、windows",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.ImageOsDistros...),
				},
			},
			"os_version": schema.StringAttribute{
				Required:    true,
				Description: "操作系统版本。注意：参数值的取值应根据系统实际情况，建议参考（以下列出osDistro所列取值对应的osVersion参考取值）：anolis：7.9、centos：7.8、ctyunos：2.0.1、debian：9.0.0、fedora：36、kylin：V10_sp1、openEuler：20.03、ubuntu：18.04、UnionTech：V20_1050u1e、windows：2008",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"architecture": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "镜像系统架构，aarch64：AArch64架构，仅支持UEFI启动方式、x86_64：x86_64架构，支持BIOS和UEFI启动方式，注意：所指定的镜像系统架构应受所指定的资源池支持",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.ImageArchitectures...),
				},
			},
			"boot_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "启动方式，bios：BIOS启动方式、uefi：UEFI启动方式，注意：若镜像系统架构为aarch64，则对启动方式的指定不生效。此参数无默认值，不指定则表示使用镜像系统架构的默认启动方式（x86_64架构的默认启动方式为BIOS）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ImageBootModes...),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "镜像描述信息。注意：长度为1~128个字符。",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(1, 128),
				},
			},
			"disk_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "磁盘容量，单位为GB，取值范围：最小5（默认值），最大1024。注意：磁盘容量不能小于镜像文件的大小；若小于镜像文件的大小，则实际的磁盘容量将使用镜像文件的大小",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Between(5, 1024),
				},
				Default: int64default.StaticInt64(5),
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "镜像种类，system：系统盘镜像，data：数据盘镜像，默认为系统盘镜像system",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.ImageTypes...),
				},
				Default: stringdefault.StaticString(business.ImageTypeSystemDiskImage),
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目id，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraProjectId, false),
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults2.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "镜像状态，accepted：已接受共享镜像，active：正常，deactivated：已弃用，deactivating：弃用中，deleted：已删除，deleting：删除中，error：错误，importing：导入中，killed：上传出错，镜像不可读，pending_delete：等待删除中，queued：排队中，reactivating：取消弃用中，rejected：已拒绝共享镜像，saving：保存中，syncing：同步中，uploading：上传中，waiting：等待接受/拒绝共享镜像",
			},
			// "maximum_ram": schema.Int64Attribute{
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Description: "最大内存，单位为GB，取值范围：0（默认值，即不限制）/1/2/4/8/16/32/64/128/256/512。注意：若取值不为0且所指定的最小内存也不为不限制时，则取值应大于或等于所指定的最小内存",
			// 	Validators: []validator.Int64{
			// 		int64validator.OneOf(0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512),
			// 	},
			// 	Default: int64default.StaticInt64(0),
			// },
			// "minimum_ram": schema.Int64Attribute{
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Description: "最小内存，单位为GB，取值范围：0（默认值，即不限制）/1/2/4/8/16/32/64/128/256/512。注意：若取值不为0且所指定的最小内存也不为不限制时，则取值应大于或等于所指定的最小内存",
			// 	Validators: []validator.Int64{
			// 		int64validator.OneOf(0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512),
			// 	},
			// 	Default: int64default.StaticInt64(0),
			// },
		},
	}
}

func (c *ctyunImage) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan CtyunImageConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建实例
	regionId := plan.RegionId.ValueString()
	projectId := plan.ProjectId.ValueString()
	imageType, err := business.ImageTypeMap.FromOriginalScene(plan.Type.ValueString(), business.ImageTypeMapScene1)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	resp, err := c.meta.Apis.CtImageApis.ImageImportApi.Do(ctx, c.meta.Credential, &ctimage.ImageImportRequest{
		RegionId:        regionId,
		ProjectId:       projectId,
		ImageFileSource: plan.FileSource.ValueString(),
		ImageProperties: ctimage.ImageImportImagePropertiesRequest{
			ImageName:    plan.Name.ValueString(),
			OsDistro:     plan.OsDistro.ValueString(),
			OsVersion:    plan.OsVersion.ValueString(),
			Architecture: plan.Architecture.ValueString(),
			BootMode:     plan.BootMode.ValueString(),
			Description:  plan.Description.ValueString(),
			DiskSize:     int(plan.DiskSize.ValueInt64()),
			ImageType:    imageType.(string),
			// MaximumRam:   int(plan.MaximumRam.ValueInt64()),
			// MinimumRam:   int(plan.MinimumRam.ValueInt64()),
		},
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	plan.Id = types.StringValue(resp.Images[0].ImageId)
	plan.RegionId = types.StringValue(regionId)
	plan.ProjectId = types.StringValue(projectId)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 轮询镜像上传的状态
	err = c.waitForUploadImageActive(ctx, plan)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 查询镜像状态信息
	instance, ctyunRequestError := c.getAndMergeImage(ctx, plan)
	if ctyunRequestError != nil {
		response.Diagnostics.AddError(ctyunRequestError.Error(), ctyunRequestError.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunImage) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state CtyunImageConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := c.getAndMergeImage(ctx, state)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	if instance == nil {
		response.State.RemoveResource(ctx)
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunImage) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state CtyunImageConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	var plan CtyunImageConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)

	id := state.Id.ValueString()
	_, err := c.meta.Apis.CtImageApis.ImageUpdateApi.Do(ctx, c.meta.Credential, &ctimage.ImageUpdateRequest{
		ImageId:     id,
		RegionId:    state.RegionId.ValueString(),
		BootMode:    plan.BootMode.ValueString(),
		Description: plan.Description.ValueString(),
		ImageName:   plan.Name.ValueString(),
		// MaximumRam:  int(plan.MaximumRam.ValueInt64()),
		// MinimumRam:  int(plan.MinimumRam.ValueInt64()),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	instance, err2 := c.getAndMergeImage(ctx, state)
	if err2 != nil {
		response.Diagnostics.AddError(err2.Error(), err2.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunImage) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state CtyunImageConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := c.meta.Apis.CtImageApis.ImageDeleteApi.Do(ctx, c.meta.Credential, &ctimage.ImageDeleteRequest{
		ImageId:  state.Id.ValueString(),
		RegionId: state.RegionId.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	// 等待镜像删除成功
	e := c.waitForImageDeleted(ctx, state)
	if e != nil {
		response.Diagnostics.AddError(e.Error(), e.Error())
		return
	}
}

// 导入命令：terraform import [配置标识].[导入配置名称] [imageId],[regionId]
func (c *ctyunImage) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var cfg CtyunImageConfig
	var imageId, regionId string
	err := terraform_extend.Split(request.ID, &imageId, &regionId)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	cfg.Id = types.StringValue(imageId)
	cfg.RegionId = types.StringValue(regionId)

	instance, err := c.getAndMergeImage(ctx, cfg)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, instance)...)
}

func (c *ctyunImage) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// waitForImageDeleted 等待镜像删除成功
func (c *ctyunImage) waitForImageDeleted(ctx context.Context, cfg CtyunImageConfig) error {
	executeSuccessFlag := false
	retryer, _ := business.NewRetryer(time.Second*5, 60)
	retryer.Start(
		func(currentTime int) bool {
			response, err := c.meta.Apis.CtImageApis.ImageDetailApi.Do(ctx, c.meta.Credential, &ctimage.ImageDetailRequest{
				RegionId: cfg.RegionId.ValueString(),
				ImageId:  cfg.Id.ValueString(),
			})
			if err != nil {
				// 执行完成后，查询不到镜像会抛错，这个是正常的出口
				if err.ErrorCode() == common.ImageImageCheckNotFound {
					executeSuccessFlag = true
					return false
				}
				return false
			}
			// 执行完成后，可能查询不到镜像的信息了，这个也是正常出口
			if len(response.Images) == 0 {
				executeSuccessFlag = true
				return false
			}
			// 其余的情况，需要按照镜像的状态来判断是否需要继续执行
			switch response.Images[0].Status {
			case business.ImageStatusDeleting:
				// 执行中
				return true
			default:
				// 默认为执行失败
				return false
			}
		},
	)

	if !executeSuccessFlag {
		return errors.New("轮询镜像" + cfg.Id.ValueString() + "状态失败")
	}
	return nil
}

// waitForUploadImageActive 等待镜像上传完成
func (c *ctyunImage) waitForUploadImageActive(ctx context.Context, cfg CtyunImageConfig) error {
	executeSuccessFlag := false
	retryer, _ := business.NewRetryer(time.Second*5, 60)
	retryer.Start(
		func(currentTime int) bool {
			response, err := c.meta.Apis.CtImageApis.ImageDetailApi.Do(ctx, c.meta.Credential, &ctimage.ImageDetailRequest{
				ImageId:  cfg.Id.ValueString(),
				RegionId: cfg.RegionId.ValueString(),
			})
			if err != nil {
				return false
			}
			if len(response.Images) != 1 {
				return false
			}
			switch response.Images[0].Status {
			case business.ImageStatusQueued:
				return true
			case business.ImageStatusActive:
				executeSuccessFlag = true
				return false
			default:
				// 默认为执行失败
				return false
			}
		},
	)

	if !executeSuccessFlag {
		return errors.New("轮询镜像" + cfg.Id.ValueString() + "状态失败")
	}
	return nil
}

// getAndMergeImage 查询合并镜像
func (c *ctyunImage) getAndMergeImage(ctx context.Context, cfg CtyunImageConfig) (*CtyunImageConfig, error) {
	response, err := c.meta.Apis.CtImageApis.ImageDetailApi.Do(ctx, c.meta.Credential, &ctimage.ImageDetailRequest{
		ImageId:  cfg.Id.ValueString(),
		RegionId: cfg.RegionId.ValueString(),
	})
	if err != nil {
		if err.ErrorCode() == common.ImageImageCheckNotFound {
			return nil, nil
		}
		return nil, err
	}
	if len(response.Images) == 0 {
		return nil, nil
	}
	resp := response.Images[0]

	imageType, err2 := business.ImageTypeMap.ToOriginalScene(resp.ImageType, business.ImageTypeMapScene1)
	if err2 != nil {
		return nil, err2
	}
	cfg.Id = types.StringValue(resp.ImageId)
	cfg.Name = types.StringValue(resp.ImageName)
	cfg.OsDistro = types.StringValue(resp.OsDistro)
	cfg.OsVersion = types.StringValue(resp.OsVersion)
	cfg.Architecture = types.StringValue(resp.Architecture)
	cfg.BootMode = types.StringValue(resp.BootMode)
	cfg.Description = types.StringValue(resp.Description)
	cfg.DiskSize = types.Int64Value(int64(resp.DiskSize))
	cfg.Type = types.StringValue(imageType.(string))
	cfg.Status = types.StringValue(resp.Status)
	// cfg.MaximumRam = types.Int64Value(int64(resp.MaximumRam))
	// cfg.MinimumRam = types.Int64Value(int64(resp.MinimumRam))
	return &cfg, nil
}

type CtyunImageConfig struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	FileSource   types.String `tfsdk:"file_source"`
	OsDistro     types.String `tfsdk:"os_distro"`
	OsVersion    types.String `tfsdk:"os_version"`
	Architecture types.String `tfsdk:"architecture"`
	BootMode     types.String `tfsdk:"boot_mode"`
	Description  types.String `tfsdk:"description"`
	DiskSize     types.Int64  `tfsdk:"disk_size"`
	Type         types.String `tfsdk:"type"`
	ProjectId    types.String `tfsdk:"project_id"`
	RegionId     types.String `tfsdk:"region_id"`
	Status       types.String `tfsdk:"status"`
	// MaximumRam   types.Int64  `tfsdk:"maximum_ram"`
	// MinimumRam   types.Int64  `tfsdk:"minimum_ram"`
}

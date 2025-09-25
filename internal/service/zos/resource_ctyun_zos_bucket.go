package zos

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctzos"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunZosBucket{}
	_ resource.ResourceWithConfigure   = &ctyunZosBucket{}
	_ resource.ResourceWithImportState = &ctyunZosBucket{}
)

type ctyunZosBucket struct {
	meta *common.CtyunMetadata
}

func NewCtyunZosBucket() resource.Resource {
	return &ctyunZosBucket{}
}

func (c *ctyunZosBucket) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_zos_bucket"
}

type CtyunZosBucketConfig struct {
	ID             types.String `tfsdk:"id"`
	RegionID       types.String `tfsdk:"region_id"`
	ACL            types.String `tfsdk:"acl"`
	Name           types.String `tfsdk:"name"`
	Bucket         types.String `tfsdk:"bucket"`
	ProjectID      types.String `tfsdk:"project_id"`
	StorageType    types.String `tfsdk:"storage_type"`
	IsEncrypted    types.Bool   `tfsdk:"is_encrypted"`
	CmkUUID        types.String `tfsdk:"cmk_uuid"`
	AzPolicy       types.String `tfsdk:"az_policy"`
	VersionEnabled types.Bool   `tfsdk:"version_enabled"`
	Tags           types.Map    `tfsdk:"tags"`
	LogEnabled     types.Bool   `tfsdk:"log_enabled"`
	LogBucket      types.String `tfsdk:"log_bucket"`
	LogPrefix      types.String `tfsdk:"log_prefix"`
	RetentionMode  types.String `tfsdk:"retention_mode"`
	RetentionDay   types.Int64  `tfsdk:"retention_day"`
	RetentionYear  types.Int64  `tfsdk:"retention_year"`
}

func (c *ctyunZosBucket) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026735/10181237**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "名称",
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
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"bucket": schema.StringAttribute{
				Required:    true,
				Description: "桶名称，不可为空。长度3-63个字符内（含）字符只能有大小写字母、数字以及中划线（-）。禁止两个中划线（-）相邻。禁止中划线（-）作为开头或结尾。",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 63),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9](?:[a-zA-Z0-9]|[-][a-zA-Z0-9])+$"), "桶名称不符合规则"),
				},
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "标签，支持更新",
				Validators: []validator.Map{
					mapvalidator.SizeAtMost(10),
				},
			},
			"acl": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "桶权限，可选值为'private'、'public-read'、'public-read-write'，分别表示私有、公共读、公共读写，默认为'private'，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ZosAclPrivate, business.ZosAclPublicRead, business.ZosAclPublicReadWrite),
				},
				Default: stringdefault.StaticString(business.ZosAclPrivate),
			},
			"version_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否启用版本控制，默认不启用。若启用后暂停，将无法在桶内创建新的历史版本，之前创建的历史版本会保留，支持更新",
				Default:     booldefault.StaticBool(false),
				Validators: []validator.Bool{
					validator2.CrossFieldBool(
						path.MatchRoot("retention_mode"),
						[]attr.Value{types.StringValue("COMPLIANCE")},
						[]attr.Value{types.BoolValue(true)}),
				},
			},
			"log_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否启用日志转存策略，默认不启用。支持更新",
				Default:     booldefault.StaticBool(false),
			},
			"log_bucket": schema.StringAttribute{
				Optional:    true,
				Description: "日志存储桶，当log_enabled为true时必填，可以让日志传递到拥有的任何桶（包括当前桶），支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("log_enabled"),
						types.BoolValue(true),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("log_enabled"),
						types.BoolValue(false),
					),
				},
			},
			"log_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "日志生成的目录+前缀，如log/logfile，当log_enabled为true时必填，支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("log_enabled"),
						types.BoolValue(true),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("log_enabled"),
						types.BoolValue(false),
					),
				},
			},
			"retention_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "合规保留模式，创建后不支持修改。默认为空，表示不开启合规保留，若填写，目前只支持为COMPLIANCE，且version_enabled必须为true",
				Validators: []validator.String{
					stringvalidator.OneOf("COMPLIANCE"),
					stringvalidator.Any(
						stringvalidator.AlsoRequires(path.MatchRoot("retention_year")),
						stringvalidator.AlsoRequires(path.MatchRoot("retention_day")),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"retention_day": schema.Int64Attribute{
				Optional:    true,
				Description: "合规保留天数，支持1-1000，当retention_mode不为空时，retention_day与retention_year只能填写其中之一，支持更新",
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("retention_mode")),
					int64validator.ConflictsWith(path.MatchRoot("retention_year")),
					int64validator.Between(1, 1000),
				},
			},
			"retention_year": schema.Int64Attribute{
				Optional:    true,
				Description: "合规保留年数，支持1-60，当retention_mode不为空时，retention_day与retention_year只能填写其中之一，支持更新",
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("retention_mode")),
					int64validator.ConflictsWith(path.MatchRoot("retention_day")),
					int64validator.Between(1, 60),
				},
			},
			"cmk_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "密钥管理服务中创建的密钥ID，当is_encrypted为true时，会自动创建密钥",
			},
			"is_encrypted": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "加密状态，默认false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"storage_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "存储类型，可选的值STANDARD、STANDARD_IA、GLACIER，分别表示标准、低频、归档，默认STANDARD",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.ZosStorageTypeStandard, business.ZosStorageTypeStandardIA, business.ZosStorageTypeGlacier),
				},
				Default: stringdefault.StaticString(business.ZosStorageTypeStandard),
			},
			"az_policy": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "az策略，可选值为single-az、multi-az，分别表示单az、多az，默认为single-az",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.ZosAzPolicySingle, business.ZosAzPolicyMulti),
				},
				Default: stringdefault.StaticString(business.ZosAzPolicySingle),
			},
		},
	}
}

func (c *ctyunZosBucket) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunZosBucketConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	err = c.createExtra(ctx, plan)
	if err != nil {
		return
	}
	// 调用查询接口触发服务端缓存刷新
	business.NewZosService(c.meta).GetZosBucketInfo(ctx, plan.Bucket.ValueString(), plan.RegionID.ValueString())
	time.Sleep(30 * time.Second)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunZosBucket) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunZosBucketConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found bucket") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunZosBucket) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunZosBucketConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunZosBucketConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	state.RetentionMode = plan.RetentionMode
	state.RetentionDay = plan.RetentionDay
	state.RetentionYear = plan.RetentionYear
	state.ACL = plan.ACL // 没有接口查询acl，所以只能这样更新
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunZosBucket) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunZosBucketConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	//response.State.RemoveResource(ctx)
}

func (c *ctyunZosBucket) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [bucket],[regionID]
func (c *ctyunZosBucket) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunZosBucketConfig
	var bucket, regionID string
	err = terraform_extend.Split(request.ID, &bucket, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.Bucket = types.StringValue(bucket)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunZosBucket) checkBeforeCreate(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosGetOssServiceStatusRequest{
		RegionID: plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosGetOssServiceStatusApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	if resp.ReturnObj.State != "true" {
		err = fmt.Errorf("您尚未在该资源池开通对象存储服务，请前往控制台开通后使用")
		return
	}
	return
}

// create 创建
func (c *ctyunZosBucket) create(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosCreateBucketRequest{
		RegionID:    plan.RegionID.ValueString(),
		ACL:         plan.ACL.ValueString(),
		Bucket:      plan.Bucket.ValueString(),
		ProjectID:   plan.ProjectID.ValueString(),
		CmkUUID:     plan.CmkUUID.ValueString(),
		IsEncrypted: plan.IsEncrypted.ValueBoolPointer(),
		StorageType: plan.StorageType.ValueString(),
		AZPolicy:    plan.AzPolicy.ValueString(),
	}
	if plan.RetentionMode.ValueString() != "" {
		t := true
		params.OtherBucketInfo = &ctzos.ZosCreateBucketOtherBucketInfoRequest{
			ObjectLockEnabledForBucket: &t,
		}
	}

	resp, err := c.meta.Apis.SdkCtZosApis.ZosCreateBucketApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// createExtra 创建进阶配置
func (c *ctyunZosBucket) createExtra(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	// 版本控制
	if plan.VersionEnabled.ValueBool() {
		err = c.setVersionEnabled(ctx, plan)
		if err != nil {
			return
		}
	}
	// 标签
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		err = c.setTags(ctx, plan)
		if err != nil {
			return
		}
	}
	// 日志转存设置
	if plan.LogEnabled.ValueBool() {
		err = c.setLogConfig(ctx, plan)
		if err != nil {
			return
		}
	}
	// 合规保留设置
	if plan.RetentionMode.ValueString() != "" {
		err = c.setRetention(ctx, plan)
		if err != nil {
			return
		}
	}

	return
}

// update 更新桶
func (c *ctyunZosBucket) update(ctx context.Context, plan, state CtyunZosBucketConfig) (err error) {
	plan.ID = state.ID
	if !plan.ACL.Equal(state.ACL) {
		err = c.setAcl(ctx, plan)
		if err != nil {
			return
		}
	}
	if !plan.VersionEnabled.Equal(state.VersionEnabled) {
		err = c.setVersionEnabled(ctx, plan)
		if err != nil {
			return
		}
	}
	if !plan.Tags.Equal(state.Tags) {
		err = c.setTags(ctx, plan)
		if err != nil {
			return
		}
	}
	if !plan.LogEnabled.Equal(state.LogEnabled) || !plan.LogBucket.Equal(state.LogBucket) || !plan.LogPrefix.Equal(state.LogPrefix) {
		err = c.setLogConfig(ctx, plan)
		if err != nil {
			return
		}
	}
	if !plan.RetentionDay.Equal(state.RetentionDay) || !plan.RetentionYear.Equal(state.RetentionYear) {
		err = c.setRetention(ctx, plan)
		if err != nil {
			return
		}
	}
	return
}

// setAcl 设置桶acl
func (c *ctyunZosBucket) setAcl(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosPutBucketAclRequest{
		Bucket:   plan.Bucket.ValueString(),
		RegionID: plan.RegionID.ValueString(),
		ACL:      plan.ACL.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosPutBucketAclApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// setRetention 设置合规保留策略
func (c *ctyunZosBucket) setRetention(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosPutObjectLockConfRequest{
		Bucket:        plan.Bucket.ValueString(),
		RegionID:      plan.RegionID.ValueString(),
		RetentionMode: plan.RetentionMode.ValueString(),
		Days:          plan.RetentionDay.ValueInt64(),
		Years:         plan.RetentionYear.ValueInt64(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosPutObjectLockConfApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// setLog 设置转存日志配置
func (c *ctyunZosBucket) setLogConfig(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	if !plan.LogEnabled.ValueBool() {
		if plan.ID.ValueString() != "" {
			// 删除日志转存
			params := &ctzos.ZosDeleteBucketLoggingRequest{
				Bucket:   plan.Bucket.ValueString(),
				RegionID: plan.RegionID.ValueString(),
			}
			resp, err2 := c.meta.Apis.SdkCtZosApis.ZosDeleteBucketLoggingApi.Do(ctx, c.meta.SdkCredential, params)
			if err2 != nil {
				err = err2
				return
			} else if resp.StatusCode == common.ErrorStatusCode {
				err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
				return
			}
			return
		}
		return
	}
	params := &ctzos.ZosPutBucketLoggingRequest{
		Bucket:   plan.Bucket.ValueString(),
		RegionID: plan.RegionID.ValueString(),
		BucketLoggingStatus: &ctzos.ZosPutBucketLoggingBucketLoggingStatusRequest{
			LoggingEnabled: &ctzos.ZosPutBucketLoggingBucketLoggingStatusLoggingEnabledRequest{
				TargetPrefix: plan.LogPrefix.ValueString(),
				TargetBucket: plan.LogBucket.ValueString(),
			}},
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosPutBucketLoggingApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// getLog 查询转存日志配置
func (c *ctyunZosBucket) getLogConfig(ctx context.Context, plan *CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosGetBucketLoggingRequest{
		Bucket:   plan.Bucket.ValueString(),
		RegionID: plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosGetBucketLoggingApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	if resp.ReturnObj.LoggingEnabled == nil {
		plan.LogEnabled = types.BoolValue(false)
		plan.LogBucket = types.StringNull()
		plan.LogPrefix = types.StringNull()
	} else {
		plan.LogEnabled = types.BoolValue(true)
		plan.LogBucket = types.StringValue(resp.ReturnObj.LoggingEnabled.TargetBucket)
		plan.LogPrefix = types.StringValue(resp.ReturnObj.LoggingEnabled.TargetPrefix)
	}

	return
}

// setTags 设置标签
func (c *ctyunZosBucket) setTags(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	if plan.ID.ValueString() != "" {
		// 清空旧标签
		params := &ctzos.ZosDeleteBucketTaggingRequest{
			Bucket:   plan.Bucket.ValueString(),
			RegionID: plan.RegionID.ValueString(),
		}
		resp, err2 := c.meta.Apis.SdkCtZosApis.ZosDeleteBucketTaggingApi.Do(ctx, c.meta.SdkCredential, params)
		if err2 != nil {
			err = err2
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
			return
		}
	}

	var tags map[string]string
	tags, err = utils.TypesMapToStringMap(ctx, plan.Tags)
	if err != nil {
		return
	}
	if len(tags) == 0 {
		return
	}
	params := &ctzos.ZosPutBucketTaggingRequest{
		Bucket:   plan.Bucket.ValueString(),
		RegionID: plan.RegionID.ValueString(),
		Tagging:  &ctzos.ZosPutBucketTaggingTaggingRequest{TagSet: []*ctzos.ZosPutBucketTaggingTaggingTagSetRequest{}},
	}
	for k, v := range tags {
		params.Tagging.TagSet = append(params.Tagging.TagSet, &ctzos.ZosPutBucketTaggingTaggingTagSetRequest{
			Key:   k,
			Value: v,
		})
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosPutBucketTaggingApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// getTags 查询桶标签
func (c *ctyunZosBucket) getTags(ctx context.Context, plan *CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosGetBucketTaggingRequest{
		Bucket:   plan.Bucket.ValueString(),
		RegionID: plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosGetBucketTaggingApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode && !strings.Contains(resp.Message, "NoSuchTagSetError") {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	tags := map[string]string{}
	if resp.ReturnObj != nil {
		for _, t := range resp.ReturnObj.TagSet {
			tags[t.Key] = t.Value
		}
	}
	t, err := utils.MapStringToTypesMap(ctx, tags)
	if err != nil {
		return
	}
	plan.Tags = t

	return
}

// setVersionEnabled 设置桶版本控制管理
func (c *ctyunZosBucket) setVersionEnabled(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosPutBucketVersioningRequest{
		Bucket:                  plan.Bucket.ValueString(),
		RegionID:                plan.RegionID.ValueString(),
		VersioningConfiguration: &ctzos.ZosPutBucketVersioningVersioningConfigurationRequest{},
	}
	if plan.VersionEnabled.ValueBool() {
		params.VersioningConfiguration.Status = "Enabled"
	} else {
		params.VersioningConfiguration.Status = "Suspended"
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosPutBucketVersioningApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

// getVersionEnabled 查询桶版本控制管理
func (c *ctyunZosBucket) getVersionEnabled(ctx context.Context, plan *CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosGetBucketVersioningRequest{
		Bucket:   plan.Bucket.ValueString(),
		RegionID: plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosGetBucketVersioningApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	if resp.ReturnObj.Status == "Enabled" {
		plan.VersionEnabled = types.BoolValue(true)
	} else {
		plan.VersionEnabled = types.BoolValue(false)
	}

	return
}

// getAndMerge 从远端查询
func (c *ctyunZosBucket) getAndMerge(ctx context.Context, plan *CtyunZosBucketConfig) (err error) {
	b, err := business.NewZosService(c.meta).GetZosBucketInfo(ctx, plan.Bucket.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return
	}
	plan.AzPolicy = types.StringValue(b.AZPolicy)
	plan.StorageType = types.StringValue(b.StorageType)
	plan.CmkUUID = utils.SecStringValue(b.CmkUUID)
	if b.CmkUUID != nil {
		plan.IsEncrypted = types.BoolValue(true)
	} else {
		plan.IsEncrypted = types.BoolValue(false)
	}
	plan.Name = plan.Bucket
	plan.ID = plan.Bucket
	// 获取version_enabled
	err = c.getVersionEnabled(ctx, plan)
	if err != nil {
		return
	}
	// 获取tags
	err = c.getTags(ctx, plan)
	if err != nil {
		return
	}
	// 获取日志转存
	err = c.getLogConfig(ctx, plan)
	if err != nil {
		return
	}

	// 以下字段无接口查询
	// plan.Acl
	return
}

// delete 删除
func (c *ctyunZosBucket) delete(ctx context.Context, plan CtyunZosBucketConfig) (err error) {
	params := &ctzos.ZosDeleteBucketRequest{
		Bucket:   plan.Bucket.ValueString(),
		RegionID: plan.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtZosApis.ZosDeleteBucketApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	}
	return
}

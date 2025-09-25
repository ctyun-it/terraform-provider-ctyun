package zos

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &ctyunZosBucketObject{}
	_ resource.ResourceWithConfigure   = &ctyunZosBucketObject{}
	_ resource.ResourceWithImportState = &ctyunZosBucketObject{}
)

type ctyunZosBucketObject struct {
	meta *common.CtyunMetadata
}

func NewCtyunZosBucketObject() resource.Resource {
	return &ctyunZosBucketObject{}
}

func (c *ctyunZosBucketObject) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_zos_bucket_object"
}

type CtyunZosBucketObjectConfig struct {
	ID                 types.String `tfsdk:"id"`
	RegionID           types.String `tfsdk:"region_id"`
	Bucket             types.String `tfsdk:"bucket"`
	Key                types.String `tfsdk:"key"`
	Source             types.String `tfsdk:"source"`
	Content            types.String `tfsdk:"content"`
	ACL                types.String `tfsdk:"acl"`
	CacheControl       types.String `tfsdk:"cache_control"`
	ContentDisposition types.String `tfsdk:"content_disposition"`
	ContentEncoding    types.String `tfsdk:"content_encoding"`
	ContentType        types.String `tfsdk:"content_type"`
	Tags               types.Map    `tfsdk:"tags"`
	StorageType        types.String `tfsdk:"storage_type"`
	Etag               types.String `tfsdk:"etag"`
	VersionID          types.String `tfsdk:"version_id"`
	client             *s3.S3
}

func (c *ctyunZosBucketObject) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026735/10181324**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
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
			"bucket": schema.StringAttribute{
				Required:    true,
				Description: "桶名称",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 63),
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-zA-Z0-9](?:[a-zA-Z0-9]|[-][a-zA-Z0-9])+$"), "桶名称不符合规则"),
				},
			},
			"key": schema.StringAttribute{
				Required:    true,
				Description: "对象名称，长度1-1024",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1024),
				},
			},
			"source": schema.StringAttribute{
				Optional:    true,
				Description: "文件路径，和content有且只能有一个",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("content"),
					}...),
				},
			},
			"content": schema.StringAttribute{
				Optional:    true,
				Description: "内容，和source有且只能有一个",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("source"),
					}...),
				},
			},
			"acl": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "权限，可选值为'private'、'public-read'、'public-read-write'，分别表示私有、公共读、公共读写，默认为'private'，支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ZosAclPrivate, business.ZosAclPublicRead, business.ZosAclPublicReadWrite),
				},
				Default: stringdefault.StaticString(business.ZosAclPrivate),
			},
			"cache_control": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "指定缓存行为，对应S3协议Header中的Cache-Control",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"content_disposition": schema.StringAttribute{
				Computed:    true,
				Description: "该对象的表示性信息，对应S3协议Header中的Content-Disposition",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"content_encoding": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "指定已对该对象应用哪些内容编码方式，对应S3协议Header中的Content-Encoding",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"content_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述对象类型，对应S3协议Header中的Content-Type",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"storage_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "存储类型，可选的值STANDARD、STANDARD_IA、GLACIER，分别表示标准、低频、归档，默认使用桶的storage_type",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(business.ZosStorageTypeStandard, business.ZosStorageTypeStandardIA, business.ZosStorageTypeGlacier),
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
			"etag": schema.StringAttribute{
				Computed:    true,
				Description: "该对象生成的实体标签（ETag）（即该对象内容的 MD5 哈希值）",
			},
			"version_id": schema.StringAttribute{
				Computed:    true,
				Description: "对象版本号",
			},
		},
	}
}

func (c *ctyunZosBucketObject) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunZosBucketObjectConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	// 获取s3 client
	err = c.initClient(ctx, &plan)
	if err != nil {
		return
	}
	// 创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunZosBucketObject) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunZosBucketObjectConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 获取s3 client
	err = c.initClient(ctx, &state)
	if err != nil {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") || strings.Contains(err.Error(), "NotFound") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunZosBucketObject) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunZosBucketObjectConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunZosBucketObjectConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 获取s3 client
	err = c.initClient(ctx, &state)
	if err != nil {
		return
	}
	// 更新
	err = c.update(ctx, &plan, &state)
	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunZosBucketObject) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunZosBucketObjectConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 获取s3 client
	err = c.initClient(ctx, &state)
	if err != nil {
		return
	}
	// 删除
	err = c.delete(ctx, &state)
	if err != nil {
		return
	}
}

func (c *ctyunZosBucketObject) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [key],[bucket],[regionID]
func (c *ctyunZosBucketObject) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunZosBucketObjectConfig
	var key, bucket, regionID string
	err = terraform_extend.Split(request.ID, &key, &bucket, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.Key = types.StringValue(key)
	cfg.Bucket = types.StringValue(bucket)
	err = c.initClient(ctx, &cfg)
	if err != nil {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeCreate 创建前检查
func (c *ctyunZosBucketObject) checkBeforeCreate(ctx context.Context, plan CtyunZosBucketObjectConfig) (err error) {
	// 检查桶是否存在
	zosService := business.NewZosService(c.meta)
	_, err = zosService.GetZosBucketInfo(ctx, plan.Bucket.ValueString(), plan.RegionID.ValueString())
	if err != nil {
		return
	}
	return
}

// create 创建
func (c *ctyunZosBucketObject) create(ctx context.Context, plan CtyunZosBucketObjectConfig) (err error) {
	source, content := plan.Source.ValueString(), plan.Content.ValueString()

	var body io.ReadSeeker
	if source != "" {
		file, e := os.Open(source)
		if e != nil {
			err = fmt.Errorf("文件 (%s) 打开时出错: %s", source, e.Error())
			return
		}
		body = file
		defer file.Close()
	} else if content != "" {
		body = bytes.NewReader([]byte(content))
	} else {
		err = fmt.Errorf("source 和 content 不能都为空")
		return
	}

	input := &s3.PutObjectInput{
		Bucket: plan.Bucket.ValueStringPointer(),
		Key:    plan.Key.ValueStringPointer(),
		Body:   body,
	}

	input.ACL = plan.ACL.ValueStringPointer()
	input.CacheControl = plan.CacheControl.ValueStringPointer()
	input.ContentType = plan.ContentType.ValueStringPointer()
	input.ContentEncoding = plan.ContentEncoding.ValueStringPointer()
	input.StorageClass = plan.StorageType.ValueStringPointer()

	tags, err := utils.TypesMapToStringMap(ctx, plan.Tags)
	if err != nil {
		return
	}
	q := utils.MapToQueryString(tags)
	if len(q) > 0 {
		input.Tagging = &q
	}

	_, err = plan.client.PutObject(input)
	if err != nil {
		return
	}

	return
}

// getAndMerge 从远端查询
func (c *ctyunZosBucketObject) getAndMerge(ctx context.Context, plan *CtyunZosBucketObjectConfig) (err error) {
	input := &s3.HeadObjectInput{
		Bucket: plan.Bucket.ValueStringPointer(),
		Key:    plan.Key.ValueStringPointer(),
	}
	output, err := plan.client.HeadObject(input)
	if err != nil {
		return
	}
	plan.ContentDisposition = utils.SecStringValue(output.ContentDisposition)
	plan.ContentEncoding = utils.SecStringValue(output.ContentEncoding)
	plan.ContentType = utils.SecStringValue(output.ContentType)
	plan.CacheControl = utils.SecStringValue(output.CacheControl)
	sType := utils.SecString(output.StorageClass)
	if sType == "" {
		sType = business.ZosStorageTypeStandard
	}
	plan.StorageType = types.StringValue(sType)

	plan.Etag = utils.SecStringValue(output.ETag)
	plan.VersionID = utils.SecStringValue(output.VersionId)
	plan.ID = types.StringValue(
		fmt.Sprintf("%s,%s,%s",
			plan.Key.ValueString(),
			plan.Bucket.ValueString(),
			plan.RegionID.ValueString(),
		),
	)

	err = c.getTags(ctx, plan)
	if err != nil {
		return
	}

	return
}

// update 更新
func (c *ctyunZosBucketObject) update(ctx context.Context, plan, state *CtyunZosBucketObjectConfig) (err error) {
	err = c.updateACL(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.updateTags(ctx, plan, state)
	if err != nil {
		return
	}
	return
}

// delete 删除
func (c *ctyunZosBucketObject) delete(ctx context.Context, plan *CtyunZosBucketObjectConfig) (err error) {
	input := &s3.DeleteObjectInput{
		Bucket: plan.Bucket.ValueStringPointer(),
		Key:    plan.Key.ValueStringPointer(),
	}
	_, err = plan.client.DeleteObject(input)
	if err != nil {
		return
	}
	return
}

// getTags 查询标签
func (c *ctyunZosBucketObject) getTags(ctx context.Context, plan *CtyunZosBucketObjectConfig) (err error) {
	input := &s3.GetObjectTaggingInput{
		Bucket: plan.Bucket.ValueStringPointer(),
		Key:    plan.Key.ValueStringPointer(),
	}
	output, err := plan.client.GetObjectTagging(input)
	if err != nil {
		return
	}
	tags := map[string]string{}
	for _, t := range output.TagSet {
		tags[*t.Key] = *t.Value
	}
	t, err := utils.MapStringToTypesMap(ctx, tags)
	if err != nil {
		return
	}
	plan.Tags = t
	return
}

// updateACL 更新ACL
func (c *ctyunZosBucketObject) updateACL(ctx context.Context, plan, state *CtyunZosBucketObjectConfig) (err error) {
	if plan.ACL.Equal(state.ACL) {
		return
	}
	input := &s3.PutObjectAclInput{
		Bucket: state.Bucket.ValueStringPointer(),
		Key:    state.Key.ValueStringPointer(),
		ACL:    plan.ACL.ValueStringPointer(),
	}
	_, err = state.client.PutObjectAcl(input)
	if err != nil {
		return
	}
	state.ACL = plan.ACL
	return
}

// updateTags 更新标签
func (c *ctyunZosBucketObject) updateTags(ctx context.Context, plan, state *CtyunZosBucketObjectConfig) (err error) {
	if plan.Tags.Equal(state.Tags) {
		return
	}

	tags, err := utils.TypesMapToStringMap(ctx, plan.Tags)
	if err != nil {
		return
	}
	if len(tags) == 0 {
		input := &s3.DeleteObjectTaggingInput{
			Bucket: state.Bucket.ValueStringPointer(),
			Key:    state.Key.ValueStringPointer(),
		}
		_, err = state.client.DeleteObjectTagging(input)
		return
	}

	input := &s3.PutObjectTaggingInput{
		Bucket: state.Bucket.ValueStringPointer(),
		Key:    state.Key.ValueStringPointer(),
		Tagging: &s3.Tagging{
			TagSet: []*s3.Tag{},
		},
	}

	for key, value := range tags {
		k, v := key, value
		input.Tagging.TagSet = append(input.Tagging.TagSet, &s3.Tag{Key: &k, Value: &v})
	}
	_, err = state.client.PutObjectTagging(input)
	return

}

// initClient 获取s3 client
func (c *ctyunZosBucketObject) initClient(ctx context.Context, plan *CtyunZosBucketObjectConfig) (err error) {
	if plan.client == nil {
		// 获取s3 client
		plan.client, err = business.NewZosService(c.meta).BuildS3Client(ctx, plan.RegionID.ValueString())
		if err != nil {
			return
		}
	}
	return
}

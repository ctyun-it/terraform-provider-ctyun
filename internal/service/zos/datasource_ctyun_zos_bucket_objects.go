package zos

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunZosBucketObjects{}
	_ datasource.DataSourceWithConfigure = &ctyunZosBucketObjects{}
)

type ctyunZosBucketObjects struct {
	meta *common.CtyunMetadata
}

func NewCtyunZosBucketObjects() datasource.DataSource {
	return &ctyunZosBucketObjects{}
}

func (c *ctyunZosBucketObjects) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_zos_bucket_objects"
}

type CtyunZosBucketObjectsModel struct {
	Key          types.String `tfsdk:"key"`
	Size         types.Int64  `tfsdk:"size"`
	LastModified types.String `tfsdk:"last_modified"`
	StorageType  types.String `tfsdk:"storage_type"`
	Etag         types.String `tfsdk:"etag"`
}

type CtyunZosBucketObjectsConfig struct {
	RegionID  types.String `tfsdk:"region_id"`
	Bucket    types.String `tfsdk:"bucket"`
	Prefix    types.String `tfsdk:"prefix"`
	Delimiter types.String `tfsdk:"delimiter"`
	Marker    types.String `tfsdk:"marker"`
	MaxKeys   types.Int64  `tfsdk:"max_keys"`

	IsTruncated types.Bool   `tfsdk:"is_truncated"`
	NextMarker  types.String `tfsdk:"next_marker"`

	Objects []CtyunZosBucketObjectsModel `tfsdk:"objects"`
	client  *s3.S3
}

func (c *ctyunZosBucketObjects) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026735/10181324**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"bucket": schema.StringAttribute{
				Required:    true,
				Description: "桶名称",
			},
			"prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Key前缀",
			},
			"delimiter": schema.StringAttribute{
				Optional:    true,
				Description: "分隔符",
			},
			"marker": schema.StringAttribute{
				Optional:    true,
				Description: "指定开始键",
			},
			"max_keys": schema.Int64Attribute{
				Optional:    true,
				Description: "最大数量",
				Validators: []validator.Int64{
					int64validator.Between(1, 1000),
				},
			},
			"next_marker": schema.StringAttribute{
				Computed:    true,
				Description: "下一次查询开始键",
			},
			"is_truncated": schema.BoolAttribute{
				Computed:    true,
				Description: "是否被截断",
			},
			"objects": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:    true,
							Description: "对象名称",
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "大小",
						},
						"last_modified": schema.StringAttribute{
							Computed:    true,
							Description: "上次修改时间",
						},
						"storage_type": schema.StringAttribute{
							Computed:    true,
							Description: "存储类型，可选的值STANDARD、STANDARD_IA、GLACIER，分别表示标准、低频、归档，默认STANDARD，",
						},
						"etag": schema.StringAttribute{
							Computed:    true,
							Description: "该对象生成的实体标签（ETag）（即该对象内容的 MD5 哈希值）",
						},
					},
				},
			},
		},
	}
}
func (c *ctyunZosBucketObjects) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunZosBucketObjectsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	config.RegionID = types.StringValue(regionId)
	err = c.initClient(ctx, &config)
	if err != nil {
		return
	}
	// 组装请求体
	input := &s3.ListObjectsInput{
		Bucket:    config.Bucket.ValueStringPointer(),
		Prefix:    config.Prefix.ValueStringPointer(),
		Delimiter: config.Delimiter.ValueStringPointer(),
		Marker:    config.Marker.ValueStringPointer(),
		MaxKeys:   config.MaxKeys.ValueInt64Pointer(),
	}
	// 调用API
	output, err := config.client.ListObjects(input)
	if err != nil {
		return
	}
	// 解析返回值
	config.Objects = []CtyunZosBucketObjectsModel{}
	config.IsTruncated = utils.SecBoolValue(output.IsTruncated)
	config.NextMarker = utils.SecStringValue(output.NextMarker)
	for _, obj := range output.Contents {
		item := CtyunZosBucketObjectsModel{
			Key:          utils.SecStringValue(obj.Key),
			StorageType:  utils.SecStringValue(obj.StorageClass),
			Etag:         utils.SecStringValue(obj.ETag),
			LastModified: types.StringValue(obj.LastModified.String()),
		}
		if obj.Size != nil {
			item.Size = types.Int64Value(*obj.Size)
		}
		config.Objects = append(config.Objects, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunZosBucketObjects) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// initClient 获取s3 client
func (c *ctyunZosBucketObjects) initClient(ctx context.Context, plan *CtyunZosBucketObjectsConfig) (err error) {
	if plan.client == nil {
		// 获取s3 client
		plan.client, err = business.NewZosService(c.meta).BuildS3Client(ctx, plan.RegionID.ValueString())
		if err != nil {
			return
		}
	}
	return
}

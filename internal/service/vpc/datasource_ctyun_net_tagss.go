package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCtyunNetTagss() datasource.DataSource {
	return &ctyunNetTagss{}
}

type ctyunNetTagss struct {
	meta        *common.CtyunMetadata
	TagsService *business.TagsService
}

func (c *ctyunNetTagss) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_net_tagss"
}

func (c *ctyunNetTagss) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028310`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID，值",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，默认使用provider ctyun总region_id 或者环境变量",
			},
			"resource_type": schema.StringAttribute{
				Required:    true,
				Description: "资源类型，resourceType only support vpc / subnet / acl / security_group / route_table / havip / port / multicast_domain / vpc_peer / vpce_endpoint / vpce_endpoint_service / ipv6_gateway / elb / private_nat / nat / eip / bandwidth /ipv6_bandwidth",
				Validators: []validator.String{
					stringvalidator.OneOf(business.NetResourceTypes...),
				},
			},
			"resource_id": schema.StringAttribute{
				Required:    true,
				Description: "资源 ID",
			},
			"page_no": schema.Int64Attribute{
				Optional:    true,
				Description: "页码，从1开始，默认为1",
			},
			"page_size": schema.Int64Attribute{
				Optional:    true,
				Description: "每页记录数，取值范围1-50，默认为10",
				Validators: []validator.Int64{
					int64validator.Between(1, 50),
				},
			},
			"tags": schema.SetNestedAttribute{
				Optional:    true,
				Description: "标签列表。最多10个标签，标签键不可重复，键值长度1~32字符，不能换行或以空格开头/结尾。",
				Validators: []validator.Set{
					setvalidator.SizeAtMost(10),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "标签id。",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
							},
						},
						"key": schema.StringAttribute{
							Computed:    true,
							Description: "标签键。支持更新",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
							},
						},
						"value": schema.StringAttribute{
							Computed:    true,
							Description: "标签值。支持更新",
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 32),
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunNetTagss) Read(ctx context.Context, req datasource.ReadRequest, response *datasource.ReadResponse) {
	var config CtyunNetTagssConfig
	response.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		msg := "regionId不能为空"
		response.Diagnostics.AddError(msg, msg)
		return
	}
	params := &ctvpc.CtvpcQueryLabelsByResourceRequest{
		RegionID:     regionId,
		ResourceType: config.ResourceType.ValueString(),
		ResourceID:   config.ResourceID.ValueString(),
	}
	if !config.PageNo.IsNull() {
		params.PageNumber = int32(config.PageNo.ValueInt64())
	}

	if !config.PageSize.IsNull() {
		pageSize := int32(config.PageSize.ValueInt64())
		params.PageSize = pageSize
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcQueryLabelsByResourceApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var resultTags []business.Tag
	for _, apiTag := range resp.ReturnObj.Results {
		tag := business.Tag{
			LabelID:    types.StringPointerValue(apiTag.LabelID),
			LabelKey:   types.StringPointerValue(apiTag.LabelKey),
			LabelValue: types.StringPointerValue(apiTag.LabelValue),
		}
		resultTags = append(resultTags, tag)
	}
	tags, diags := types.SetValueFrom(ctx, utils.StructToTFObjectTypes(business.Tag{}), resultTags)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return
	}
	config.Tags = tags
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunNetTagss) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.TagsService = business.NewTagsService(meta)
}

type CtyunNetTagssConfig struct {
	ID           types.String `tfsdk:"id"`
	RegionID     types.String `tfsdk:"region_id"`     //区域id
	ResourceID   types.String `tfsdk:"resource_id"`   //需要创建 NAT 网关的 VPC 的 ID
	ResourceType types.String `tfsdk:"resource_type"` //需要创建 NAT 网关的 VPC 的 ID
	PageNo       types.Int64  `tfsdk:"page_no"`
	PageSize     types.Int64  `tfsdk:"page_size"`
	Tags         types.Set    `tfsdk:"tags"`
}

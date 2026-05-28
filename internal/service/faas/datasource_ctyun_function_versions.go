package faas

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/cf"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunFunctionVersions{}
	_ datasource.DataSourceWithConfigure = &CtyunFunctionVersions{}
)

func NewCtyunFunctionVersions() datasource.DataSource {
	return &CtyunFunctionVersions{}
}

type CtyunFunctionVersions struct {
	meta *common.CtyunMetadata
}

func (c *CtyunFunctionVersions) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function_versions"
}

func (c *CtyunFunctionVersions) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询函数版本列表", "函数计算（FaaS）", "https://www.ctyun.cn/document/10006234/10825950"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据源唯一标识",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID，如果不填则默认使用 provider ctyun 中的 region_id 或环境变量中的 CTYUN_REGION_ID",
			},
			"function_name": schema.StringAttribute{
				Required:    true,
				Description: "函数名称",
			},
			"page_index": schema.Int32Attribute{
				Optional:    true,
				Description: "页码，默认为 1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页大小，默认为 10",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"search": schema.StringAttribute{
				Optional:    true,
				Description: "模糊查询的关键字，默认为空",
			},
			"versions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "版本列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"version_id": schema.StringAttribute{
							Computed:    true,
							Description: "版本 ID",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "版本描述",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
						"binding_alias": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "版本关联的别名列表",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunFunctionVersions) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunFunctionVersions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CtyunFunctionVersionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(data.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		resp.Diagnostics.AddError("region_id 不能为空！", "region_id 不能为空！")
		return
	}
	data.RegionID = types.StringValue(regionId)
	listReq := &cf.CfListFunctionVersionsRequest{
		FunctionName: data.FunctionName.ValueString(),
		RegionId:     data.RegionID.ValueString(),
		PageIndex:    data.PageIndex.ValueInt32Pointer(),
		PageSize:     data.PageSize.ValueInt32Pointer(),
	}

	if !data.Search.IsNull() {
		search := data.Search.ValueString()
		listReq.Search = &search
	}

	listResp, err := c.meta.Apis.SdkCtCfApis.CfListFunctionVersionsApi.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		resp.Diagnostics.AddError("查询函数版本列表失败", err.Error())
		return
	}

	if *listResp.StatusCode == common.ErrorStatusCode {
		resp.Diagnostics.AddError("查询函数版本列表失败", fmt.Sprintf("API 返回错误：%s", *listResp.Message))
		return
	}

	if listResp.ReturnObj == nil {
		resp.Diagnostics.AddError("查询函数版本列表失败", "API 返回空结果")
		return
	}

	var versions []FunctionVersionModel
	for _, version := range listResp.ReturnObj.Data {
		versionModel := FunctionVersionModel{
			VersionID:   types.StringPointerValue(version.VersionId),
			Description: types.StringPointerValue(version.Description),
			CreateTime:  types.StringPointerValue(version.CreateTime),
			UpdateTime:  types.StringPointerValue(version.UpdateTime),
		}

		if version.BindingAlias != nil {
			bindingAliases := make([]types.String, len(version.BindingAlias))
			for i, alias := range version.BindingAlias {
				bindingAliases[i] = types.StringPointerValue(alias)
			}
			versionModel.BindingAlias = bindingAliases
		}

		versions = append(versions, versionModel)
	}

	data.Versions = versions
	data.ID = types.StringValue(fmt.Sprintf("function-versions-%s-%s-%d", data.RegionID.ValueString(), data.FunctionName.ValueString(), data.PageIndex.ValueInt32()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type FunctionVersionModel struct {
	VersionID    types.String   `tfsdk:"version_id"`
	Description  types.String   `tfsdk:"description"`
	CreateTime   types.String   `tfsdk:"create_time"`
	UpdateTime   types.String   `tfsdk:"update_time"`
	BindingAlias []types.String `tfsdk:"binding_alias"`
}

type CtyunFunctionVersionsDataSourceModel struct {
	ID           types.String           `tfsdk:"id"`
	RegionID     types.String           `tfsdk:"region_id"`
	FunctionName types.String           `tfsdk:"function_name"`
	PageIndex    types.Int32            `tfsdk:"page_index"`
	PageSize     types.Int32            `tfsdk:"page_size"`
	Search       types.String           `tfsdk:"search"`
	Versions     []FunctionVersionModel `tfsdk:"versions"`
}

func convertInt32ToString(i int32) string {
	return strconv.Itoa(int(i))
}

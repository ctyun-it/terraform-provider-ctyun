package faas

import (
	"context"
	"fmt"
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
	_ datasource.DataSource              = &CtyunFunctionTriggers{}
	_ datasource.DataSourceWithConfigure = &CtyunFunctionTriggers{}
)

func NewCtyunFunctionTriggers() datasource.DataSource {
	return &CtyunFunctionTriggers{}
}

type CtyunFunctionTriggers struct {
	meta *common.CtyunMetadata
}

func (c *CtyunFunctionTriggers) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function_triggers"
}

func (c *CtyunFunctionTriggers) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询函数触发器列表", "函数工作流（FunctionGraph）", "https://www.ctyun.cn/document/10355289"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据源唯一标识",
			},
			"function_name": schema.StringAttribute{
				Required:    true,
				Description: "函数名称",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池 ID，如果不填则默认使用 provider ctyun 中的 region_id 或环境变量中的 CTYUN_REGION_ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目 ID，如果不填则默认使用 provider ctyun 中的 project_id 或环境变量中的 CTYUN_PROJECT_ID",
			},
			"trigger_name": schema.StringAttribute{
				Optional:    true,
				Description: "触发器名称，支持模糊搜索",
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
			"version": schema.StringAttribute{
				Optional:    true,
				Description: "版本或别名",
			},
			"triggers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "触发器列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"trigger_name": schema.StringAttribute{
							Computed:    true,
							Description: "触发器名称",
						},
						"trigger_type": schema.StringAttribute{
							Computed:    true,
							Description: "触发器类型",
						},
						"trigger_config": schema.StringAttribute{
							Computed:    true,
							Description: "触发器配置，JSON 字符串",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "版本或别名",
						},
						"status": schema.Int32Attribute{
							Computed:    true,
							Description: "触发器状态。1：启用；2：禁用；3：系统禁用",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
						"url_internet": schema.StringAttribute{
							Computed:    true,
							Description: "外网 URL",
						},
						"url_intranet": schema.StringAttribute{
							Computed:    true,
							Description: "内网 URL",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunFunctionTriggers) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunFunctionTriggers) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CtyunFunctionTriggersDataSourceModel
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

	listReq := &cf.CfListTriggersV2024Request{
		FunctionName: data.FunctionName.ValueString(),
		RegionId:     data.RegionID.ValueString(),
	}

	if !data.TriggerName.IsNull() {
		triggerName := data.TriggerName.ValueString()
		listReq.TriggerName = &triggerName
	}

	if !data.PageIndex.IsNull() {
		pageIndex := data.PageIndex.ValueInt32()
		listReq.PageIndex = &pageIndex
	}

	if !data.PageSize.IsNull() {
		pageSize := data.PageSize.ValueInt32()
		listReq.PageSize = &pageSize
	}

	if !data.Version.IsNull() {
		version := data.Version.ValueString()
		listReq.Version = &version
	}

	listResp, err := c.meta.Apis.SdkCtCfApis.CfListTriggersV2024Api.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		resp.Diagnostics.AddError("查询触发器列表失败", err.Error())
		return
	}

	if *listResp.StatusCode == common.ErrorStatusCode {
		resp.Diagnostics.AddError("查询触发器列表失败", fmt.Sprintf("API 返回错误：%s", *listResp.Message))
		return
	}

	if listResp.ReturnObj == nil {
		resp.Diagnostics.AddError("查询触发器列表失败", "API 返回空结果")
		return
	}

	// 转换 API 响应数据到 Terraform 模型
	var triggers []TriggerModel
	for _, trigger := range listResp.ReturnObj.Data {
		triggerModel := TriggerModel{
			TriggerName:   types.StringPointerValue(trigger.TriggerName),
			TriggerType:   types.StringPointerValue(trigger.TriggerType),
			TriggerConfig: types.StringPointerValue(trigger.TriggerConfig),
			Version:       types.StringPointerValue(trigger.Version),
			Status:        types.Int32PointerValue(trigger.Status),
			CreatedAt:     types.StringPointerValue(trigger.CreatedAt),
			UpdatedAt:     types.StringPointerValue(trigger.UpdatedAt),
			UrlInternet:   types.StringPointerValue(trigger.UrlInternet),
			UrlIntranet:   types.StringPointerValue(trigger.UrlIntranet),
		}

		triggers = append(triggers, triggerModel)
	}

	data.Triggers = triggers
	data.ID = types.StringValue(fmt.Sprintf("triggers-%s-%s", data.FunctionName.ValueString(), data.RegionID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type TriggerModel struct {
	TriggerName   types.String `tfsdk:"trigger_name"`
	TriggerType   types.String `tfsdk:"trigger_type"`
	TriggerConfig types.String `tfsdk:"trigger_config"`
	Version       types.String `tfsdk:"version"`
	Status        types.Int32  `tfsdk:"status"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	UrlInternet   types.String `tfsdk:"url_internet"`
	UrlIntranet   types.String `tfsdk:"url_intranet"`
}

type CtyunFunctionTriggersDataSourceModel struct {
	ID           types.String   `tfsdk:"id"`
	FunctionName types.String   `tfsdk:"function_name"`
	RegionID     types.String   `tfsdk:"region_id"`
	ProjectID    types.String   `tfsdk:"project_id"`
	TriggerName  types.String   `tfsdk:"trigger_name"`
	PageIndex    types.Int32    `tfsdk:"page_index"`
	PageSize     types.Int32    `tfsdk:"page_size"`
	Version      types.String   `tfsdk:"version"`
	Triggers     []TriggerModel `tfsdk:"triggers"`
}

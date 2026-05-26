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
	_ datasource.DataSource              = &CtyunFunctions{}
	_ datasource.DataSourceWithConfigure = &CtyunFunctions{}
)

func NewCtyunFunctions() datasource.DataSource {
	return &CtyunFunctions{}
}

type CtyunFunctions struct {
	meta *common.CtyunMetadata
}

func (c *CtyunFunctions) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_functions"
}

func (c *CtyunFunctions) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询函数列表", "函数工作流（FunctionGraph）", "https://www.ctyun.cn/document/10355289"),
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
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目 ID，如果不填则默认使用 provider ctyun 中的 project_id 或环境变量中的 CTYUN_PROJECT_ID",
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
			"functions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "函数列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"function_id": schema.StringAttribute{
							Computed:    true,
							Description: "函数 ID",
						},
						"function_name": schema.StringAttribute{
							Computed:    true,
							Description: "函数名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "函数描述",
						},
						"create_type": schema.Int32Attribute{
							Computed:    true,
							Description: "创建类型 1:内置运行时 2:自定义运行时 3:自定义镜像",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "函数状态",
						},
						"runtime": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "运行时配置",
							Attributes: map[string]schema.Attribute{
								"runtime": schema.StringAttribute{
									Computed:    true,
									Description: "运行时类型",
								},
								"handler": schema.StringAttribute{
									Computed:    true,
									Description: "函数执行入口",
								},
								"execute_timeout": schema.Int32Attribute{
									Computed:    true,
									Description: "执行超时时间",
								},
							},
						},
						"container": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "容器配置",
							Attributes: map[string]schema.Attribute{
								"time_zone": schema.StringAttribute{
									Computed:    true,
									Description: "时区",
								},
								"memory_size": schema.Int32Attribute{
									Computed:    true,
									Description: "内存规格 (Mb)",
								},
								"cpu": schema.Float64Attribute{
									Computed:    true,
									Description: "CPU 规格 (vCPU)",
								},
								"disk_size": schema.Int32Attribute{
									Computed:    true,
									Description: "磁盘规格 (Mb)",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *CtyunFunctions) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunFunctions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CtyunFunctionsDataSourceModel
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

	listReq := &cf.CfListFunctionsV2024Request{
		RegionId:  data.RegionID.ValueString(),
		PageIndex: data.PageIndex.ValueInt32Pointer(),
		PageSize:  data.PageSize.ValueInt32Pointer(),
	}

	if !data.Search.IsNull() {
		search := data.Search.ValueString()
		listReq.Search = &search
	}

	listResp, err := c.meta.Apis.SdkCtCfApis.CfListFunctionsV2024Api.Do(ctx, c.meta.SdkCredential, listReq)
	if err != nil {
		resp.Diagnostics.AddError("查询函数列表失败", err.Error())
		return
	}

	if *listResp.StatusCode == common.ErrorStatusCode {
		resp.Diagnostics.AddError("查询函数列表失败", fmt.Sprintf("API 返回错误：%s", *listResp.Message))
		return
	}

	if listResp.ReturnObj == nil {
		resp.Diagnostics.AddError("查询函数列表失败", "API 返回空结果")
		return
	}

	// 转换 API 响应数据到 Terraform 模型
	var functions []FunctionModel
	for _, function := range listResp.ReturnObj.Data {
		functionModel := FunctionModel{
			FunctionID:   types.StringPointerValue(function.FunctionId),
			FunctionName: types.StringPointerValue(function.FunctionName),
			Description:  types.StringPointerValue(function.Description),
			CreateType:   types.Int32PointerValue(function.CreateType),
		}

		if function.DeployInfo != nil {
			functionModel.Status = types.StringPointerValue(function.DeployInfo.Status)
		}

		if function.Runtime != nil {
			functionModel.Runtime = &RuntimeModel{
				Runtime:        types.StringPointerValue(function.Runtime.Runtime),
				Handler:        types.StringPointerValue(function.Runtime.Handler),
				ExecuteTimeout: types.Int32PointerValue(function.Runtime.ExecuteTimeout),
			}
		}

		if function.Container != nil {
			functionModel.Container = &ContainerModel{
				TimeZone:   types.StringPointerValue(function.Container.TimeZone),
				MemorySize: types.Int32PointerValue(function.Container.MemorySize),
				Cpu:        types.Float64PointerValue(function.Container.Cpu),
				DiskSize:   types.Int32PointerValue(function.Container.DiskSize),
			}
		}

		functions = append(functions, functionModel)
	}

	data.Functions = functions
	data.ID = types.StringValue(fmt.Sprintf("functions-%s-%d", data.RegionID.ValueString(), data.PageIndex.ValueInt32()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type RuntimeModel struct {
	Runtime        types.String `tfsdk:"runtime"`
	Handler        types.String `tfsdk:"handler"`
	ExecuteTimeout types.Int32  `tfsdk:"execute_timeout"`
}

type ContainerModel struct {
	TimeZone   types.String  `tfsdk:"time_zone"`
	MemorySize types.Int32   `tfsdk:"memory_size"`
	Cpu        types.Float64 `tfsdk:"cpu"`
	DiskSize   types.Int32   `tfsdk:"disk_size"`
}

type FunctionModel struct {
	FunctionID   types.String    `tfsdk:"function_id"`
	FunctionName types.String    `tfsdk:"function_name"`
	Description  types.String    `tfsdk:"description"`
	CreateType   types.Int32     `tfsdk:"create_type"`
	Status       types.String    `tfsdk:"status"`
	Runtime      *RuntimeModel   `tfsdk:"runtime"`
	Container    *ContainerModel `tfsdk:"container"`
}

type CtyunFunctionsDataSourceModel struct {
	ID        types.String    `tfsdk:"id"`
	RegionID  types.String    `tfsdk:"region_id"`
	ProjectID types.String    `tfsdk:"project_id"`
	PageIndex types.Int32     `tfsdk:"page_index"`
	PageSize  types.Int32     `tfsdk:"page_size"`
	Search    types.String    `tfsdk:"search"`
	Functions []FunctionModel `tfsdk:"functions"`
}

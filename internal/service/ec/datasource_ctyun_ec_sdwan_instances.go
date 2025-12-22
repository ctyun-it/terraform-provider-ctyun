package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunEcSdwanInstances{}
	_ datasource.DataSourceWithConfigure = &CtyunEcSdwanInstances{}
)

func NewCtyunEcSdwanInstances() datasource.DataSource {
	return &CtyunEcSdwanInstances{}
}

type CtyunEcSdwanInstances struct {
	meta *common.CtyunMetadata
}

type CtyunEcSdwanInstancesConfig struct {
	ID           types.String               `tfsdk:"id"`
	EcID         types.String               `tfsdk:"ec_id"`
	SdwanID      types.String               `tfsdk:"sdwan_id"`
	CgwID        types.String               `tfsdk:"cgw_id"`
	InstanceID   types.String               `tfsdk:"instance_id"`
	QueryContent types.String               `tfsdk:"query_content"`
	IsAuth       types.Int64                `tfsdk:"is_auth"`
	Instances    []CtyunEcSdwanInstanceInfo `tfsdk:"instances"`
}

type CtyunEcSdwanInstanceInfo struct {
	EcID              types.String `tfsdk:"ec_id"`
	CgwID             types.String `tfsdk:"cgw_id"`
	CgwName           types.String `tfsdk:"cgw_name"`
	DcID              types.String `tfsdk:"dc_id"`
	SdwanID           types.String `tfsdk:"sdwan_id"`
	SdwanName         types.String `tfsdk:"sdwan_name"`
	InstanceID        types.String `tfsdk:"instance_id"`
	DefaultRtbID      types.String `tfsdk:"default_rtb_id"`
	DefaultRtbName    types.String `tfsdk:"default_rtb_name"`
	CIDR              types.List   `tfsdk:"cidr"`
	V6CIDR            types.List   `tfsdk:"v6_cidr"`
	CreateDate        types.String `tfsdk:"create_date"`
	Status            types.String `tfsdk:"status"`
	Weights           types.Int64  `tfsdk:"weights"`
	RedundantType     types.Int64  `tfsdk:"redundant_type"`
	RedundantInstUUID types.String `tfsdk:"redundant_inst_uuid"`
	RedundantInstName types.String `tfsdk:"redundant_inst_name"`
	RedundantInstType types.String `tfsdk:"redundant_inst_type"`
	RedundantInstID   types.String `tfsdk:"redundant_inst_id"`
	IsAuth            types.Int64  `tfsdk:"is_auth"`
	RouteLearn        types.Int64  `tfsdk:"route_learn"`
	RouteSync         types.Int64  `tfsdk:"route_sync"`
}

func (c *CtyunEcSdwanInstances) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ec_sdwan_instances"
}

func (c *CtyunEcSdwanInstances) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038220`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "数据源ID",
			},
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sdwan_id": schema.StringAttribute{
				Optional:    true,
				Description: "sdwan ID",
			},
			"cgw_id": schema.StringAttribute{
				Optional:    true,
				Description: "云网关ID",
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "网络实例ID",
			},
			"query_content": schema.StringAttribute{
				Optional:    true,
				Description: "模糊查询，支持sdwanID，网络实例ID，云网关ID三个属性",
			},
			"is_auth": schema.Int64Attribute{
				Optional:    true,
				Description: "是否是跨账号实例，取值包括：0:本账号，1:跨账号",
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"instances": schema.ListNestedAttribute{
				Computed:    true,
				Description: "返回的sdwan网络实例列表",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ec_id": schema.StringAttribute{
							Computed:    true,
							Description: "云间高速ID",
						},
						"cgw_id": schema.StringAttribute{
							Computed:    true,
							Description: "云网关ID",
						},
						"cgw_name": schema.StringAttribute{
							Computed:    true,
							Description: "云网关名称",
						},
						"dc_id": schema.StringAttribute{
							Computed:    true,
							Description: "云网关资源池ID",
						},
						"sdwan_id": schema.StringAttribute{
							Computed:    true,
							Description: "sdwan ID",
						},
						"sdwan_name": schema.StringAttribute{
							Computed:    true,
							Description: "sdwan名称",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "网络实例的ID",
						},
						"default_rtb_id": schema.StringAttribute{
							Computed:    true,
							Description: "默认路由表ID",
						},
						"default_rtb_name": schema.StringAttribute{
							Computed:    true,
							Description: "默认路由表名称",
						},
						"cidr": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "v4 CIDR列表",
						},
						"v6_cidr": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "v6 CIDR列表",
						},
						"create_date": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态描述，取值包括：creating:加载中, running:已连接, removing:卸载中, flushing:路由待更新, error:失败",
						},
						"weights": schema.Int64Attribute{
							Computed:    true,
							Description: "权重，sdwan默认60",
						},
						"redundant_type": schema.Int64Attribute{
							Computed:    true,
							Description: "冗余类型，取值包括：1：主备, 2：负载, 0：无",
						},
						"redundant_inst_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "冗余负载实例的UUID",
						},
						"redundant_inst_name": schema.StringAttribute{
							Computed:    true,
							Description: "冗余负载实例的名称",
						},
						"redundant_inst_type": schema.StringAttribute{
							Computed:    true,
							Description: "冗余负载实例的类型，取值包括：2: 云专线, 3: sdwan, 4: vpn",
						},
						"redundant_inst_id": schema.StringAttribute{
							Computed:    true,
							Description: "冗余负载侧ID，即当前为cdaID, sdwanID，vpnID",
						},
						"is_auth": schema.Int64Attribute{
							Computed:    true,
							Description: "是否是跨账号实例，取值包括：0:本账号, 1:跨账号",
						},
						"route_learn": schema.Int64Attribute{
							Computed:    true,
							Description: "路由学习开关，开启后云网关自动学习网络实例路由，取值范围: 1:学习, 0:不学习, 默认学习",
						},
						"route_sync": schema.Int64Attribute{
							Computed:    true,
							Description: "路由同步开关，开启后云网关路由自动同步到网络实例，取值范围: 1:同步, 0:不同步, 默认同步",
						},
					},
				},
			},
		},
	}
}

func (c *CtyunEcSdwanInstances) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	meta := req.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunEcSdwanInstances) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunEcSdwanInstancesConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 构造请求参数
	request := &ec.EcEcListSDWANInstanceRequest{
		EcID: plan.EcID.ValueString(),
	}

	if !plan.SdwanID.IsNull() {
		request.SdwanID = plan.SdwanID.ValueStringPointer()
	}

	if !plan.CgwID.IsNull() {
		request.CgwID = plan.CgwID.ValueStringPointer()
	}

	if !plan.InstanceID.IsNull() {
		request.InstanceID = plan.InstanceID.ValueStringPointer()
	}

	if !plan.QueryContent.IsNull() {
		request.QueryContent = plan.QueryContent.ValueStringPointer()
	}

	if !plan.IsAuth.IsNull() {
		isAuth := int32(plan.IsAuth.ValueInt64())
		request.IsAuth = &isAuth
	}

	response, err := c.meta.Apis.SdkEcApis.EcEcListSDWANInstanceApi.Do(ctx, c.meta.SdkCredential, request)
	if err != nil {
		return
	} else if response == nil {
		err = fmt.Errorf("API return error. StatusCode is nil")
		return
	} else if *response.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *response.Message)
		return
	} else if response.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 处理返回结果
	var instances []CtyunEcSdwanInstanceInfo
	if response.ReturnObj.Results != nil {
		for _, result := range response.ReturnObj.Results {
			instance := CtyunEcSdwanInstanceInfo{}

			if result.EcID != nil {
				instance.EcID = types.StringValue(*result.EcID)
			}

			if result.CgwID != nil {
				instance.CgwID = types.StringValue(*result.CgwID)
			}

			if result.CgwName != nil {
				instance.CgwName = types.StringValue(*result.CgwName)
			}

			if result.DcID != nil {
				instance.DcID = types.StringValue(*result.DcID)
			}

			if result.SdwanID != nil {
				instance.SdwanID = types.StringValue(*result.SdwanID)
			}

			if result.SdwanName != nil {
				instance.SdwanName = types.StringValue(*result.SdwanName)
			}

			if result.InstanceID != nil {
				instance.InstanceID = types.StringValue(*result.InstanceID)
			}

			if result.DefaultRtbID != nil {
				instance.DefaultRtbID = types.StringValue(*result.DefaultRtbID)
			}

			if result.DefaultRtbName != nil {
				instance.DefaultRtbName = types.StringValue(*result.DefaultRtbName)
			}

			if result.CIDR != nil {
				cidrValues := make([]string, len(result.CIDR))
				for i, cidr := range result.CIDR {
					if cidr != nil {
						cidrValues[i] = *cidr
					}
				}
				instance.CIDR, _ = types.ListValueFrom(ctx, types.StringType, cidrValues)
			}

			if result.V6CIDR != nil {
				v6CidrValues := make([]string, len(result.V6CIDR))
				for i, cidr := range result.V6CIDR {
					if cidr != nil {
						v6CidrValues[i] = *cidr
					}
				}
				instance.V6CIDR, _ = types.ListValueFrom(ctx, types.StringType, v6CidrValues)
			}

			if result.CreateDate != nil {
				instance.CreateDate = types.StringValue(*result.CreateDate)
			}

			if result.Status != nil {
				instance.Status = types.StringValue(*result.Status)
			}

			if result.Weights != nil {
				instance.Weights = types.Int64Value(int64(*result.Weights))
			}

			if result.RedundantType != nil {
				instance.RedundantType = types.Int64Value(int64(*result.RedundantType))
			}

			if result.RedundantInstUUID != nil {
				instance.RedundantInstUUID = types.StringValue(*result.RedundantInstUUID)
			}

			if result.RedundantInstName != nil {
				instance.RedundantInstName = types.StringValue(*result.RedundantInstName)
			}

			if result.RedundantInstType != nil {
				instance.RedundantInstType = types.StringValue(*result.RedundantInstType)
			}

			if result.RedundantInstID != nil {
				instance.RedundantInstID = types.StringValue(*result.RedundantInstID)
			}

			if result.IsAuth != nil {
				instance.IsAuth = types.Int64Value(int64(*result.IsAuth))
			}

			if result.RouteLearn != nil {
				instance.RouteLearn = types.Int64Value(int64(*result.RouteLearn))
			}

			if result.RouteSync != nil {
				instance.RouteSync = types.Int64Value(int64(*result.RouteSync))
			}

			instances = append(instances, instance)
		}
	}

	plan.Instances = instances
	plan.ID = types.StringValue("sdwan_instances")

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

package elb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &CtyunElbTargetGroup{}
	_ resource.ResourceWithConfigure   = &CtyunElbTargetGroup{}
	_ resource.ResourceWithImportState = &CtyunElbTargetGroup{}
)

type CtyunElbTargetGroup struct {
	meta *common.CtyunMetadata
}

func (c *CtyunElbTargetGroup) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func NewCtyunElbTargetGroup() resource.Resource {
	return &CtyunElbTargetGroup{}
}

func (c *CtyunElbTargetGroup) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_target_group"
}

func (c *CtyunElbTargetGroup) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10155289**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池Id，默认使用provider ctyun总region_id 或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Description: "支持 TCP / UDP / HTTP / HTTPS, 该字段不支持更新。当protocol=HTTP/HTTPS时，target_group.session_sticky_mode仅支持INSERT/REWRITE",
				Validators: []validator.String{
					stringvalidator.OneOf(business.ListenerProtocols...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "名称，唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 32),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "描述，支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:'{},./;'[,]·！@#￥%……&*（） —— -+={},，支持更新",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					validator2.Desc(),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "需要创建后端主机组的 VPC 的 ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"health_check_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "需要关联的健康检查Id，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"algorithm": schema.StringAttribute{
				Required:    true,
				Description: "调度算法。取值范围：rr（轮询）、wrr（带权重轮询）、lc（最少连接）、sh（源IP哈希），支持更新",
				Validators: []validator.String{
					stringvalidator.OneOf(business.TargetGroupAlgorithms...),
				},
			},
			"proxy_protocol": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "1 开启，0 关闭，只有protocol=tcp的时候,可填写（关闭/开启proxy_protocol），其他协议默认关闭。不支持更改",
				Default:     int32default.StaticInt32(0),
				Validators: []validator.Int32{
					int32validator.Between(0, 1),
					validator2.AlsoRequiresEqualInt32(
						path.MatchRoot("protocol"),
						types.StringValue(business.ListenerProtocolTCP),
					),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"session_sticky_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "会话保持模式，支持取值：CLOSE（关闭）、INSERT（插入）、REWRITE（重写）。当 algorithm 为 lc / sh 时，sessionStickyMode无需填写，默认为 CLOSE，支持更新",
				Default:     stringdefault.StaticString("CLOSE"),
				Validators: []validator.String{
					stringvalidator.OneOf(business.TargetGroupSessionStickyModes...),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("algorithm"),
						types.StringValue(business.TargetGroupAlgorithmLC),
						types.StringValue(business.TargetGroupAlgorithmSH),
					),
				},
			},
			"cookie_expire": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "cookie过期时间。session_sticky_mode = INSERT模式必填，支持更新",
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("session_sticky_mode"),
						types.StringValue(business.TargetGroupSessionStickyModeINSERT),
					),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("session_sticky_mode"),
						types.StringValue(business.TargetGroupSessionStickyModeREWRITE),
						types.StringValue(business.TargetGroupSessionStickyModeSourceIP),
						types.StringValue(business.TargetGroupSessionStickyModeCLOSE),
					),
				},
			},
			"rewrite_cookie_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "cookie重写名称，REWRITE模式必填，支持更新",
				Validators: []validator.String{
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("session_sticky_mode"),
						types.StringValue(business.TargetGroupSessionStickyModeREWRITE),
					),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("session_sticky_mode"),
						types.StringValue(business.TargetGroupSessionStickyModeCLOSE),
						types.StringValue(business.TargetGroupSessionStickyModeSourceIP),
						types.StringValue(business.TargetGroupSessionStickyModeINSERT),
					),
				},
			},
			"source_ip_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "源IP会话保持超时时间。SOURCE_IP模式必填，支持更新",
				Validators: []validator.Int64{
					validator2.AlsoRequiresEqualInt64(
						path.MatchRoot("session_sticky_mode"),
						types.StringValue(business.TargetGroupSessionStickyModeSourceIP),
					),
					validator2.ConflictsWithEqualInt64(
						path.MatchRoot("session_sticky_mode"),
						types.StringValue(business.TargetGroupSessionStickyModeCLOSE),
						types.StringValue(business.TargetGroupSessionStickyModeREWRITE),
						types.StringValue(business.TargetGroupSessionStickyModeINSERT),
					),
				},
			},
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "后端服务组ID",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "状态: ACTIVE / DOWN",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
			"updated_time": schema.StringAttribute{
				Computed:    true,
				Description: "更新时间，为UTC格式",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
		},
	}
}

func (c *CtyunElbTargetGroup) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunElbTargetGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 开始创建

	err = c.createTargetGroup(ctx, &plan)
	if err != nil {
		return
	}

	// 创建后，反查详情，补充plan
	err = c.getAndMergeTargetGroup(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunElbTargetGroup) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunElbTargetGroupConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeTargetGroup(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "不存在") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (c *CtyunElbTargetGroup) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置
	var plan CtyunElbTargetGroupConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunElbTargetGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.updateTargetGroupInfo(ctx, &state, &plan)
	if err != nil {
		return
	}
	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMergeTargetGroup(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (c *CtyunElbTargetGroup) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state CtyunElbTargetGroupConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &ctelb.CtelbDeleteTargetGroupRequest{
		ClientToken:   uuid.NewString(),
		RegionID:      state.RegionID.ValueString(),
		TargetGroupID: state.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbDeleteTargetGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}
	return
}

func (c *CtyunElbTargetGroup) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunElbTargetGroup) createTargetGroup(ctx context.Context, plan *CtyunElbTargetGroupConfig) (err error) {
	if plan.RegionID.IsNull() {
		err = fmt.Errorf("regionId 为空！")
		return
	}

	params := &ctelb.CtelbCreateTargetGroupRequest{
		ClientToken:   uuid.NewString(),
		RegionID:      plan.RegionID.ValueString(),
		Name:          plan.Name.ValueString(),
		VpcID:         plan.VpcID.ValueString(),
		Algorithm:     plan.Algorithm.ValueString(),
		SessionSticky: nil,
	}
	if !plan.Protocol.IsNull() {
		params.Protocol = plan.Protocol.ValueString()
	}
	if !plan.HealthCheckID.IsNull() {
		params.HealthCheckID = plan.HealthCheckID.ValueString()
	}
	if !plan.ProxyProtocol.IsNull() {
		if !c.validProxyProtocol(plan.Protocol, plan.ProxyProtocol) {
			err = fmt.Errorf("ProxyProtocol取值有误，只有protocol=TCP时，ProxyProtocol才可开启")
			return
		}
		params.ProxyProtocol = plan.ProxyProtocol.ValueInt32()
	}
	sessionSticky := &ctelb.CtelbCreateTargetGroupSessionStickyRequest{SessionStickyMode: "CLOSE"}
	if plan.SessionStickyMode.ValueString() != "" {
		// 若protocol=HTTPS或HTTPS时， session_sticky_mode 仅支持INSERT/REWRITE
		if plan.Protocol.ValueString() != "" && plan.Protocol.ValueString() == business.ListenerProtocolHTTP || plan.Protocol.ValueString() == business.ListenerProtocolHTTPS {
			if plan.SessionStickyMode.ValueString() != business.TargetGroupSessionStickyModeINSERT && plan.SessionStickyMode.ValueString() != business.TargetGroupSessionStickyModeREWRITE {
				err = errors.New("protocol=HTTPS或HTTPS时， session_sticky_mode 仅支持INSERT/REWRITE")
				return err
			}

		}
		sessionSticky.SessionStickyMode = plan.SessionStickyMode.ValueString()
		if plan.Algorithm.ValueString() == business.TargetGroupAlgorithmLC || plan.Algorithm.ValueString() == business.TargetGroupAlgorithmSH {
			//当 algorithm 为 lc / sh 时，sessionStickyMode 必须为 CLOSE
			if plan.SessionStickyMode.ValueString() != business.TargetGroupSessionStickyModeCLOSE {
				err = fmt.Errorf("当 algorithm 为 lc / sh 时，sessionStickyMode 必须为 CLOSE")
				return
			}
		}
		if !plan.CookieExpire.IsNull() && !plan.CookieExpire.IsUnknown() {
			sessionSticky.CookieExpire = int32(plan.CookieExpire.ValueInt64())
		}
		if !plan.RewriteCookieName.IsNull() && !plan.RewriteCookieName.IsUnknown() {
			sessionSticky.RewriteCookieName = plan.RewriteCookieName.ValueString()
		}
		if !plan.SourceIpTimeout.IsNull() && !plan.SourceIpTimeout.IsUnknown() {
			sessionSticky.SourceIpTimeout = int32(plan.SourceIpTimeout.ValueInt64())
		}
	}
	params.SessionSticky = sessionSticky

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbCreateTargetGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}

	// 创建成功后，保存后端服务组id
	if len(resp.ReturnObj) == 1 {
		plan.ID = types.StringValue(resp.ReturnObj[0].ID)
	} else {
		err = fmt.Errorf("返回id数量有误！")
	}

	return
}
func (c *CtyunElbTargetGroup) updateTargetGroupInfo(ctx context.Context, state *CtyunElbTargetGroupConfig, plan *CtyunElbTargetGroupConfig) (err error) {

	params := &ctelb.CtelbUpdateTargetGroupRequest{
		ClientToken:   uuid.NewString(),
		RegionID:      state.RegionID.ValueString(),
		ProjectID:     state.ProjectID.ValueString(),
		ID:            state.ID.ValueString(),
		TargetGroupID: state.ID.ValueString(),
		Name:          state.Name.ValueString(),
		HealthCheckID: state.HealthCheckID.ValueString(),
		Algorithm:     state.Algorithm.ValueString(),
		ProxyProtocol: state.ProxyProtocol.ValueInt32(),
		SessionSticky: &ctelb.CtelbUpdateTargetGroupSessionStickyRequest{
			SessionStickyMode: state.SessionStickyMode.ValueString(),
			CookieExpire:      int32(state.CookieExpire.ValueInt64()),
			RewriteCookieName: state.RewriteCookieName.ValueString(),
			SourceIpTimeout:   int32(state.SourceIpTimeout.ValueInt64()),
		},
	}

	if !state.ProjectID.IsNull() {
		params.ProjectID = plan.ProjectID.ValueString()
	}
	if !plan.Name.Equal(state.Name) {
		params.Name = plan.Name.ValueString()
	}
	if !plan.HealthCheckID.IsNull() && !plan.HealthCheckID.Equal(state.HealthCheckID) {
		params.HealthCheckID = plan.HealthCheckID.ValueString()
	}
	if !plan.Algorithm.IsNull() && !plan.Algorithm.Equal(state.Algorithm) {
		params.Algorithm = plan.Algorithm.ValueString()
	}

	if !plan.SessionStickyMode.IsNull() {
		sessionSticky := &ctelb.CtelbUpdateTargetGroupSessionStickyRequest{}
		sessionSticky.SessionStickyMode = plan.SessionStickyMode.ValueString()
		if !plan.CookieExpire.IsNull() {
			sessionSticky.CookieExpire = int32(plan.CookieExpire.ValueInt64())
		}
		if !plan.RewriteCookieName.IsNull() {
			sessionSticky.RewriteCookieName = plan.RewriteCookieName.ValueString()
		}
		if !plan.SourceIpTimeout.IsNull() {
			sessionSticky.SourceIpTimeout = int32(plan.SourceIpTimeout.ValueInt64())
		}
		params.SessionSticky = sessionSticky
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbUpdateTargetGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}

	return
}

// validProxyProtocol 验证ProxyProtocol参数时候合理，当protocol=tcp时，proxyProtocol可以为0或者1；当protocol=其他值时，proxyProtocol=0
func (c *CtyunElbTargetGroup) validProxyProtocol(protocol types.String, proxyProtocol types.Int32) bool {
	if proxyProtocol.ValueInt32() == 1 && protocol.ValueString() != business.ListenerProtocolTCP {
		return false
	}

	return true
}

func (c *CtyunElbTargetGroup) getAndMergeTargetGroup(ctx context.Context, plan *CtyunElbTargetGroupConfig) (err error) {
	// 定义查看后端主机组详情请求参数
	params := &ctelb.CtelbShowTargetGroupRequest{
		RegionID:      plan.RegionID.ValueString(),
		TargetGroupID: plan.ID.ValueString(),
	}
	// 请求查看后端主机组详情接口
	resp, err := c.meta.Apis.SdkCtElbApis.CtelbShowTargetGroupApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		return
	}

	if len(resp.ReturnObj) != 1 {
		err = fmt.Errorf("TargetGroupID与后端主机组关系不是1：1对应！，通过TargetGroupID查询出多条信息。")
	}
	// 解析详情接口返回值
	returnObj := resp.ReturnObj[0]
	plan.Status = types.StringValue(returnObj.Status)
	plan.CreatedTime = types.StringValue(returnObj.CreatedTime)
	plan.UpdatedTime = types.StringValue(returnObj.UpdatedTime)
	plan.Description = types.StringValue(returnObj.Description)
	plan.Algorithm = types.StringValue(returnObj.Algorithm)
	plan.Name = types.StringValue(returnObj.Name)
	plan.SessionStickyMode = types.StringValue(returnObj.SessionSticky.SessionStickyMode)
	plan.CookieExpire = types.Int64Value(int64(returnObj.SessionSticky.CookieExpire))
	plan.RewriteCookieName = types.StringValue(returnObj.SessionSticky.RewriteCookieName)
	plan.SourceIpTimeout = types.Int64Value(int64(returnObj.SessionSticky.SourceIpTimeout))
	plan.ProxyProtocol = types.Int32Value(returnObj.ProxyProtocol)
	plan.HealthCheckID = types.StringValue(returnObj.HealthCheckID)
	return
}

type CtyunElbTargetGroupConfig struct {
	RegionID          types.String `tfsdk:"region_id"`           //区域ID
	Protocol          types.String `tfsdk:"protocol"`            //	支持 TCP / UDP / HTTP / HTTPS
	Name              types.String `tfsdk:"name"`                //唯一。支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
	Description       types.String `tfsdk:"description"`         //支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:'{},./;'[,]·！@#￥%……&*（） —— -+={},
	VpcID             types.String `tfsdk:"vpc_id"`              //vpc ID
	HealthCheckID     types.String `tfsdk:"health_check_id"`     //健康检查ID
	Algorithm         types.String `tfsdk:"algorithm"`           //调度算法。取值范围：rr（轮询）、wrr（带权重轮询）、lc（最少连接）、sh（源IP哈希）
	ProxyProtocol     types.Int32  `tfsdk:"proxy_protocol"`      //1 开启，0 关闭
	SessionStickyMode types.String `tfsdk:"session_sticky_mode"` //会话保持模式，支持取值：CLOSE（关闭）、INSERT（插入）、REWRITE（重写），当 algorithm 为 lc / sh 时，sessionStickyMode 必须为 CLOSE
	CookieExpire      types.Int64  `tfsdk:"cookie_expire"`       //cookie过期时间。INSERT模式必填
	RewriteCookieName types.String `tfsdk:"rewrite_cookie_name"` //cookie重写名称，REWRITE模式必填
	SourceIpTimeout   types.Int64  `tfsdk:"source_ip_timeout"`   //源IP会话保持超时时间。SOURCE_IP模式必填
	ID                types.String `tfsdk:"id"`                  //后端服务组ID
	ProjectID         types.String `tfsdk:"project_id"`          //项目ID
	Status            types.String `tfsdk:"status"`              //状态: ACTIVE / DOWN
	CreatedTime       types.String `tfsdk:"created_time"`        //创建时间，为UTC格式
	UpdatedTime       types.String `tfsdk:"updated_time"`        //更新时间，为UTC格式
}

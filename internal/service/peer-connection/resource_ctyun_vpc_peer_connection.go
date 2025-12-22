package peer_connection

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &CtyunVpcPeerConnection{}
	_ resource.ResourceWithConfigure   = &CtyunVpcPeerConnection{}
	_ resource.ResourceWithImportState = &CtyunVpcPeerConnection{}
)

type CtyunVpcPeerConnection struct {
	meta          *common.CtyunMetadata
	regionService *business.RegionService
}

func (c *CtyunVpcPeerConnection) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_peer_connection"
}

func (c *CtyunVpcPeerConnection) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.regionService = business.NewRegionService(c.meta)

}

func NewCtyunVpcPeerConnection() resource.Resource {
	return &CtyunVpcPeerConnection{}
}

func (c *CtyunVpcPeerConnection) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[instanceId],[projectID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunVpcPeerConnectionConfig
	var ID, regionId, projectId, instanceId string
	// 根据分隔符数量判断是否输入了regionID,projectId
	if strings.Count(request.ID, common.ImportSeparator) == 1 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		projectId = c.meta.GetExtraIfEmpty(projectId, common.ExtraProjectId)
		err = terraform_extend.Split(request.ID, &ID, &instanceId)
		if err != nil {
			return
		}
	} else if strings.Count(request.ID, common.ImportSeparator) == 2 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &ID, &instanceId, &projectId)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &ID, &instanceId, &projectId, &regionId)
		if err != nil {
			return
		}
	}

	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}

	config.ID = types.StringValue(ID)
	config.RegionID = types.StringValue(regionId)
	if projectId != "" {
		config.ProjectID = types.StringValue(projectId)
	}
	if instanceId != "" {
		config.InstanceID = types.StringValue(instanceId)
	}
	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunVpcPeerConnection) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026760/10037873",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "对等连接id",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Description: "对等连接实例id，跨账号情况下使用，如果该字段为空，说明status=pending，需要调用ctyun_vpc_peer_connection_attch同意",
				Computed:    true,
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
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
			"name": schema.StringAttribute{
				Description: "对等连接名称，支持更新。要求：支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32。注：当status=pending时，在控制台上修改name会导致资源state丢失！",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(2, 32),
				},
			},
			"request_vpc_id": schema.StringAttribute{
				Description: "本端vpc id",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"accept_vpc_id": schema.StringAttribute{
				Description: "对端的vpc id",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.VpcValidate(),
				},
			},
			"accept_email": schema.StringAttribute{
				Description: "对端vpc账户的邮箱，当建立跨帐号的对等连接，需要对端同意。可调用ctyun_vpc_peer_connection_attch实现",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					validator2.Email(),
				},
			},
			"description": schema.StringAttribute{
				Description: "对等连接描述，支持更新",
				Optional:    true,
				Validators: []validator.String{
					validator2.Desc(),
				},
			},
			"request_vpc_name": schema.StringAttribute{
				Description: "本端的vpc名称",
				Computed:    true,
			},
			"request_vpc_cidr": schema.StringAttribute{
				Description: "本端的vpc cidr",
				Computed:    true,
			},
			"accept_vpc_name": schema.StringAttribute{
				Description: "对端的vpc名称",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"accept_vpc_cidr": schema.StringAttribute{
				Description: "对端的vpc cidr",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Description: "对等连接状态，agree(已连接)/pending(等待审核)",
				//Description: "对等连接类型：current(同一个租户) / other(不同租户)",
				Computed: true,
			},
			"user_type": schema.StringAttribute{
				Description: "对等连接类型：current(同一个租户) / other(不同租户)",
				Computed:    true,
			},
			"tags": schema.SetNestedAttribute{
				Description: "标签，支持更新",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "key，支持更新",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"value": schema.StringAttribute{
							Description: "value，支持更新",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"id": schema.StringAttribute{
							Description: "标签id",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (c *CtyunVpcPeerConnection) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunVpcPeerConnectionConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunVpcPeerConnection) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpcPeerConnectionConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "不存在") || strings.Contains(err.Error(), "not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunVpcPeerConnection) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunVpcPeerConnectionConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunVpcPeerConnectionConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &state, &plan)
	if err != nil {
		return
	}

	// 更新远端后，查询远端并同步一下本地信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunVpcPeerConnection) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunVpcPeerConnectionConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunVpcPeerConnection) create(ctx context.Context, config *CtyunVpcPeerConnectionConfig) error {
	params := &ctvpc.CtvpcCreateVpcPeerConnectionRequest{
		ClientToken:  uuid.NewString(),
		RequestVpcID: config.RequestVpcID.ValueString(),
		AcceptVpcID:  config.AcceptVpcID.ValueString(),
		Name:         config.Name.ValueString(),
		RegionID:     config.RegionID.ValueString(),
	}
	if !config.AcceptEmail.IsNull() && !config.AcceptEmail.IsUnknown() {
		params.AcceptEmail = config.AcceptEmail.ValueStringPointer()
	}
	if !config.ProjectID.IsNull() && !config.ProjectID.IsUnknown() {
		params.ProjectID = config.ProjectID.ValueString()
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		params.Description = config.Description.ValueStringPointer()
	}
	// 通过vpc_id 获取vpc_name 和 vpc_cidr
	vpcDetailResp, err := c.getVpcDetail(ctx, config.RegionID.ValueString(), config.RequestVpcID.ValueString())
	if err != nil {
		return err
	}
	vpcDetail := vpcDetailResp.ReturnObj
	params.RequestVpcCidr = *vpcDetail.CIDR
	params.RequestVpcName = *vpcDetail.Name

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateVpcPeerConnectionApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建vpc对等连接失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	// 如果是跨账号的情况，需要查询对等连接list获取uuid
	config.ID = types.StringValue(*resp.ReturnObj.InstanceID)
	config.InstanceID = types.StringValue(*resp.ReturnObj.InstanceID)
	if !config.AcceptEmail.IsNull() && !config.AcceptEmail.IsUnknown() {
		err = c.getPeerConnectionID(ctx, config)
		if err != nil {
			return err
		}
		config.ID = config.InstanceID
	}
	// 处理标签问题
	if !config.Tags.IsNull() && !config.Tags.IsUnknown() {
		var tags []CtyunVpcPeerConnectionTagsModel
		diags := config.Tags.ElementsAs(ctx, &tags, false)
		if diags.HasError() {
			err = fmt.Errorf(diags[0].Detail())
			return err
		}
		for _, tag := range tags {
			err = c.addTags(ctx, config, tag)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CtyunVpcPeerConnection) getVpcDetail(ctx context.Context, regionID string, vpcID string) (*ctvpc.CtvpcShowVpcResponse, error) {
	params := &ctvpc.CtvpcShowVpcRequest{
		RegionID: regionID,
		VpcID:    vpcID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowVpcApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取vpc详情失败(vpc_id=%s)，接口返回nil，请联系研发确认问题原因！", vpcID)
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

func (c *CtyunVpcPeerConnection) getAndMerge(ctx context.Context, config *CtyunVpcPeerConnectionConfig) error {
	resp, err := c.getPeerConnectionDetail(ctx, config)
	if err != nil {
		return err
	}
	returnObj := resp.ReturnObj
	config.Name = types.StringValue(returnObj.Name)
	config.RequestVpcID = types.StringValue(returnObj.RequestVpcID)
	config.AcceptVpcID = types.StringValue(returnObj.AcceptVpcID)
	config.RequestVpcName = types.StringValue(returnObj.RequestVpcName)
	config.RequestVpcCidr = types.StringValue(returnObj.RequestVpcCidr)
	config.AcceptVpcCidr = types.StringValue(returnObj.AcceptVpcCidr)
	config.AcceptVpcName = types.StringValue(returnObj.AcceptVpcName)
	config.UserType = types.StringValue(returnObj.UserType)
	config.Status = types.StringValue(returnObj.Status)

	// 处理tags的id
	tagList := make([]CtyunVpcPeerConnectionTagsModel, 0)
	tags, err := c.getTags(ctx, config)
	if err != nil {
		return err
	}
	for _, tagItem := range tags {
		var tag CtyunVpcPeerConnectionTagsModel
		tag.Key = types.StringValue(*tagItem.LabelKey)
		tag.Value = types.StringValue(*tagItem.LabelValue)
		tag.ID = types.StringValue(*tagItem.LabelID)
		tagList = append(tagList, tag)
	}
	if config.Tags.IsNull() {
		tagList = nil
	}
	tagListTmp, diags := types.SetValueFrom(ctx, utils.StructToTFObjectTypes(CtyunVpcPeerConnectionTagsModel{}), tagList)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return err
	}
	config.Tags = tagListTmp
	return nil
}

func (c *CtyunVpcPeerConnection) getPeerConnectionDetail(ctx context.Context, config *CtyunVpcPeerConnectionConfig) (*ctvpc.CtvpcShowVpcPeerConnectionResponse, error) {
	params := &ctvpc.CtvpcShowVpcPeerConnectionRequest{
		RegionID: config.RegionID.ValueString(),
	}
	// 邮箱不为空，表示跨账户对等连接
	if !config.AcceptEmail.IsNull() && !config.AcceptEmail.IsUnknown() {
		//1.先利用两端vpc id查询下是否已经完成连接。如果发现status=agree，则更新instance id并使用instance id
		err := c.getPeerConnectionID(ctx, config)
		if err != nil {
			return nil, err
		}
		// 如果instance id不为空，表示已经建立连接
		if !config.InstanceID.IsNull() && !config.InstanceID.IsUnknown() {
			params.InstanceID = config.InstanceID.ValueString()
		} else {
			// 如果instance id 为空，表示仍未建立连接。使用uuid
			params.InstanceID = config.ID.ValueString()
		}
	} else {
		params.InstanceID = config.ID.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowVpcPeerConnectionApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取vpc对等连接详情失败(instance_id=%s)，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp, nil
}

func (c *CtyunVpcPeerConnection) update(ctx context.Context, state *CtyunVpcPeerConnectionConfig, plan *CtyunVpcPeerConnectionConfig) error {
	params := &ctvpc.CtvpcUpdateVpcPeerConnectionAttributeRequest{
		ClientToken: uuid.NewString(),
		Name:        plan.Name.ValueStringPointer(),
		RegionID:    state.RegionID.ValueString(),
	}
	// 如果跨账号，且状态为agree时，使用instance id
	if !state.InstanceID.IsNull() && !state.InstanceID.IsUnknown() && state.Status.ValueString() == business.PeerConnectStatusAgree {
		params.InstanceID = state.InstanceID.ValueString()
	} else {
		params.InstanceID = state.ID.ValueString()
	}

	err := c.reqModifyPeerConnection(ctx, params)
	if err != nil {
		return err
	}
	state.Description = plan.Description
	// 处理标签更新
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() && !plan.Tags.Equal(state.Tags) {
		err = c.updateTags(ctx, state, plan)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CtyunVpcPeerConnection) reqModifyPeerConnection(ctx context.Context, params *ctvpc.CtvpcUpdateVpcPeerConnectionAttributeRequest) error {
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcUpdateVpcPeerConnectionAttributeApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新vpc对等连接失败(id=%s)，接口返回nil，请联系研发确认问题原因", params.InstanceID)
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	}
	return nil
}

func (c *CtyunVpcPeerConnection) delete(ctx context.Context, config CtyunVpcPeerConnectionConfig) error {
	params := &ctvpc.CtvpcDeleteVpcPeerConnectionRequest{
		ClientToken: uuid.NewString(),
		RegionID:    config.RegionID.ValueString(),
	}

	if !config.InstanceID.IsNull() && !config.InstanceID.IsUnknown() && config.Status.ValueString() == business.PeerConnectStatusAgree {
		params.InstanceID = config.InstanceID.ValueString()
	} else {
		params.InstanceID = config.ID.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteVpcPeerConnectionApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除vpc对等连接失败(id=%s)，接口返回nil，请联系研发确认问题原因", config.ID.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	}
	return nil
}

func (c *CtyunVpcPeerConnection) addTags(ctx context.Context, config *CtyunVpcPeerConnectionConfig, tag CtyunVpcPeerConnectionTagsModel) error {
	params := &ctvpc.CtvpcVpcVpcpeerBindLabelRequest{
		RegionID:   config.RegionID.ValueString(),
		LabelKey:   tag.Key.ValueString(),
		LabelValue: tag.Value.ValueString(),
	}
	// 如果跨账号，且状态为agree时，使用instance id
	if !config.InstanceID.IsNull() && !config.InstanceID.IsUnknown() && config.Status.ValueString() == business.PeerConnectStatusAgree {
		params.VpcPeerID = config.InstanceID.ValueString()
	} else {
		params.VpcPeerID = config.ID.ValueString()
	}
	err := c.LoopBindTag(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (c *CtyunVpcPeerConnection) updateTags(ctx context.Context, state *CtyunVpcPeerConnectionConfig, plan *CtyunVpcPeerConnectionConfig) error {
	var stateTags, planTags []CtyunVpcPeerConnectionTagsModel
	diags := state.Tags.ElementsAs(ctx, &stateTags, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}

	diags = plan.Tags.ElementsAs(ctx, &planTags, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}

	// 先处理增加标签的情况
	for _, planTag := range planTags {
		flag := true
		for _, stateTag := range stateTags {
			if planTag.Key.Equal(stateTag.Key) && planTag.Value.Equal(stateTag.Value) {
				flag = false
				break
			}
		}
		if flag {
			err := c.addTags(ctx, state, planTag)
			if err != nil {
				return err
			}
		}
	}
	// 处理删除标签的情况
	for _, stateTag := range stateTags {
		flag := true
		for _, planTag := range planTags {
			if planTag.Key.Equal(stateTag.Key) {
				flag = false
				break
			}
		}
		if flag {
			err := c.removeTags(ctx, state, stateTag)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func (c *CtyunVpcPeerConnection) removeTags(ctx context.Context, config *CtyunVpcPeerConnectionConfig, tag CtyunVpcPeerConnectionTagsModel) error {
	params := &ctvpc.CtvpcVpcVpcpeerUnbindLabelRequest{
		RegionID:  config.RegionID.ValueString(),
		VpcPeerID: config.ID.ValueString(),
		LabelID:   tag.ID.ValueString(),
	}

	// 如果跨账号，且状态为agree时，使用instance id
	if !config.InstanceID.IsNull() && !config.InstanceID.IsUnknown() && config.Status.ValueString() == business.PeerConnectStatusAgree {
		params.VpcPeerID = config.InstanceID.ValueString()
	} else {
		params.VpcPeerID = config.ID.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcVpcVpcpeerUnbindLabelApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("为vpc对等连接(id=%s)解绑标签失败，key=%s,value=%s，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString(), tag.Key.ValueString(), tag.Value.ValueString())
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	}
	return nil
}

func (c *CtyunVpcPeerConnection) getTags(ctx context.Context, config *CtyunVpcPeerConnectionConfig) ([]*ctvpc.CtvpcVpcVpcpeerListLabelsReturnObjResultsResponse, error) {
	params := &ctvpc.CtvpcVpcVpcpeerListLabelsRequest{
		RegionID: config.RegionID.ValueString(),
		PageSize: 50,
	}
	if !config.AcceptEmail.IsNull() && !config.AcceptEmail.IsUnknown() && config.Status.ValueString() == business.PeerConnectStatusAgree {
		params.VpcPeerID = config.InstanceID.ValueString()
	} else {
		params.VpcPeerID = config.ID.ValueString()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcVpcVpcpeerListLabelsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("查询vpc对等连接id=%s的标签列表失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj.Results, nil
}

// 当跨账号创建对等连接时，需要通过查询list来确认uuid；如果未跨账号，id为：vpr-xxxxxx
func (c *CtyunVpcPeerConnection) getPeerConnectionID(ctx context.Context, config *CtyunVpcPeerConnectionConfig) error {
	params := &ctvpc.CtvpcListVpcPeerConnectionRequest{
		PageSize: 200,
		PageNo:   1,
		RegionID: config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListVpcPeerConnectionApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("查询对等连接失败，接口返回nil。请联系研发确认问题原因！")
		return err
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}

	for _, item := range resp.ReturnObj {
		if *item.RequestVpcID == config.RequestVpcID.ValueString() && *item.AcceptVpcID == config.AcceptVpcID.ValueString() {
			config.InstanceID = types.StringValue(*item.InstanceID)
			config.Status = types.StringValue(*item.Status)
			config.UserType = types.StringValue(*item.UserType)
		}
	}
	if config.ID.IsNull() || config.ID.IsUnknown() {
		err = fmt.Errorf("未找到对等连接信息，获取id失败，请检查参数是否正确！")
		return err
	}
	return nil
}

func (c *CtyunVpcPeerConnection) checkID(idStr string) (bool, error) {
	_, err := uuid.Parse(idStr)
	if err != nil {
		// 检查前缀
		if strings.HasPrefix(idStr, "vpr-") {
			return false, nil
		}
		return false, fmt.Errorf("对等连接id格式错误，既不是uuid，也不是vpr前缀标准id！")
	}
	return true, nil
}

func (c *CtyunVpcPeerConnection) LoopBindTag(ctx context.Context, params *ctvpc.CtvpcVpcVpcpeerBindLabelRequest, loopCount ...int) error {
	var err error
	count := 61
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*10, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			resp, err2 := c.meta.Apis.SdkCtVpcApis.CtvpcVpcVpcpeerBindLabelApi.Do(ctx, c.meta.SdkCredential, params)
			if err2 != nil {
				err = err2
				return false
			} else if resp == nil {
				err = fmt.Errorf("为vpc对等连接(id=%s)绑定标签失败，key=%s,value=%s，接口返回nil，请联系研发确认问题原因！", params.VpcPeerID, params.LabelKey, params.LabelValue)
				return false
			} else if resp.StatusCode != common.NormalStatusCode {
				if strings.Contains(*resp.Description, "上报it延迟") {
					return true
				}
				err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			return false
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，对等连接的标签仍未绑定成功！")
	}
	return err
}

type CtyunVpcPeerConnectionConfig struct {
	RegionID       types.String `tfsdk:"region_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	Name           types.String `tfsdk:"name"`
	RequestVpcID   types.String `tfsdk:"request_vpc_id"`
	AcceptVpcID    types.String `tfsdk:"accept_vpc_id"`
	AcceptEmail    types.String `tfsdk:"accept_email"`
	ID             types.String `tfsdk:"id"`
	InstanceID     types.String `tfsdk:"instance_id"`
	RequestVpcName types.String `tfsdk:"request_vpc_name"`
	RequestVpcCidr types.String `tfsdk:"request_vpc_cidr"`
	AcceptVpcName  types.String `tfsdk:"accept_vpc_name"`
	AcceptVpcCidr  types.String `tfsdk:"accept_vpc_cidr"`
	Status         types.String `tfsdk:"status"`
	UserType       types.String `tfsdk:"user_type"`
	Description    types.String `tfsdk:"description"`
	Tags           types.Set    `tfsdk:"tags"`
}

type CtyunVpcPeerConnectionTagsModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
	ID    types.String `tfsdk:"id"`
}

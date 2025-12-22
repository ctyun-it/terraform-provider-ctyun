package kafka

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strings"
)

var (
	_ resource.Resource                = &ctyunKafkaUser{}
	_ resource.ResourceWithConfigure   = &ctyunKafkaUser{}
	_ resource.ResourceWithImportState = &ctyunKafkaUser{}
)

type ctyunKafkaUser struct {
	meta       *common.CtyunMetadata
	vpcService *business.VpcService
	sgService  *business.SecurityGroupService
}

func NewCtyunKafkaUser() resource.Resource {
	return &ctyunKafkaUser{}
}

func (c *ctyunKafkaUser) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_kafka_user"
}

type CtyunKafkaUserConfig struct {
	Id                 types.Int32  `tfsdk:"id"`
	UserName           types.String `tfsdk:"name"`
	InstanceId         types.String `tfsdk:"instance_id"`
	RegionId           types.String `tfsdk:"region_id"`
	Password           types.String `tfsdk:"password"`
	Description        types.String `tfsdk:"description"`
	CreateTime         types.String `tfsdk:"create_time"`
	PermissionInfo     types.Set    `tfsdk:"permission_info"`
	permissionInfoList []CtyunKafkaAclPermissionInfo
}

// ACL操作信息结构体
type CtyunKafkaAclPermissionInfo struct {
	Permission types.String `tfsdk:"permission"`
	Ip         types.String `tfsdk:"ip"`
	Operation  types.String `tfsdk:"operation"`
	Topic      types.String `tfsdk:"topic"`
}

func (c *ctyunKafkaUser) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029624/10145597`,
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Computed:      true,
				Description:   "资源唯一标识符",
				PlanModifiers: []planmodifier.Int32{int32planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "用户名称，规则如下：\n以英文字母、数字、下划线开头，且只能由英文字母、数字、句点、中划线、下划线组成。\n长度3-64。\n名称不可重复。",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 64),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9_.-]*$`),
						"必须以英文字母、数字、下划线开头，只能包含英文字母、数字、句点、中划线、下划线",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID。",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
			"password": schema.StringAttribute{
				Required:    true,
				Description: "密码，规则如下：\n长度8-26字符。\n必须同时包含大写字母、小写字母、数字和英文格式特殊符号(@%^*_+!$-=.)中的至少三种类型。\n不能有空格。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 26),
				},
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Description: "用户描述，规则如下：\n不能以+,-,@,= 特殊字符开头。\n长度不能大于200。支持更新",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 200),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^+\-@=].*$`),
						"不能以+,-,@,=特殊字符开头",
					),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:      true,
				Description:   "创建时间，为UTC格式",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

			"permission_info": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "用户ACL权限",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"permission": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "权限，ALLOW:允许，DENY:拒绝，默认：ALLOW 支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf("ALLOW", "DENY"),
							},
						},
						"ip": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "ip或网段，* 表示所有ip，默认：* 支持更新",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
						"operation": schema.StringAttribute{
							Required:    true,
							Description: "操作，READ:消费，WRITE:生产 支持更新",
							Validators: []validator.String{
								stringvalidator.OneOf("READ", "WRITE"),
							},
						},
						"topic": schema.StringAttribute{
							Required:    true,
							Description: "topic名称 支持更新",
							Validators: []validator.String{
								stringvalidator.UTF8LengthAtLeast(1),
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunKafkaUser) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunKafkaUserConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建
	err = c.create(ctx, plan)
	if err != nil {
		return
	}
	// 计算权限信息
	err = c.calcPermissionInfo(ctx, &plan)
	if err != nil {
		return
	}
	if len(plan.permissionInfoList) > 0 {
		// 更新用户ACL权限
		err = c.updateUserTopicsAcl(ctx, plan, "CREATE")
		if err != nil {
			return
		}
	}

	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunKafkaUser) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunKafkaUserConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunKafkaUser) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunKafkaUserConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunKafkaUserConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}

	// 更新用户ACL权限
	err = c.updatePermissionInfo(ctx, plan, state)
	if err != nil {
		return
	}

	state.Password = plan.Password
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunKafkaUser) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunKafkaUserConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 销毁
	err = c.destroy(ctx, state)
	if err != nil {
		return
	}
	response.State.RemoveResource(ctx)
}

func (c *ctyunKafkaUser) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.vpcService = business.NewVpcService(meta)
	c.sgService = business.NewSecurityGroupService(meta)
}

func (c *ctyunKafkaUser) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [instanceId],[userName],[password],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunKafkaUserConfig
	var instanceId, regionID, userName, password string

	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) == 2 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &instanceId, &userName, &password)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &instanceId, &userName, &password, &regionID)
		if err != nil {
			return
		}
	}

	if instanceId == "" {
		err = fmt.Errorf("instanceId不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	if userName == "" {
		err = fmt.Errorf("userName不能为空")
		return
	}
	if password == "" {
		err = fmt.Errorf("password不能为空")
		return
	}

	cfg.RegionId = types.StringValue(regionID)
	cfg.InstanceId = types.StringValue(instanceId)
	cfg.UserName = types.StringValue(userName)
	cfg.Password = types.StringValue(password)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建
func (c *ctyunKafkaUser) create(ctx context.Context, plan CtyunKafkaUserConfig) (err error) {
	params := &ctgkafka.CtgkafkaSaslUserCreateV3Request{
		RegionId:    plan.RegionId.ValueString(),
		ProdInstId:  plan.InstanceId.ValueString(),
		Username:    plan.UserName.ValueString(),
		Password:    plan.Password.ValueString(),
		Description: plan.Description.ValueString(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaSaslUserCreateV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	return
}

// update 更新
func (c *ctyunKafkaUser) update(ctx context.Context, plan, state CtyunKafkaUserConfig) (err error) {
	params := &ctgkafka.CtgkafkaUserUpdateRequest{
		RegionId:    state.RegionId.ValueString(),
		ProdInstId:  state.InstanceId.ValueString(),
		Username:    plan.UserName.ValueString(),
		NewPassword: plan.Password.ValueString(),
		Description: plan.Description.ValueString(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaUserUpdateApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		return fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj == nil {
		return common.InvalidReturnObjError
	} else if resp.ReturnObj.Data != "update success" {
		return fmt.Errorf("API return error. Data: %s", resp.ReturnObj.Data)
	}
	return
}

// destroy 销毁
func (c *ctyunKafkaUser) destroy(ctx context.Context, plan CtyunKafkaUserConfig) (err error) {
	params := &ctgkafka.CtgkafkaSaslUserDeleteV3Request{
		RegionId:   plan.RegionId.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		Username:   plan.UserName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaSaslUserDeleteV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// getAndMerge 从远端查询
func (c *ctyunKafkaUser) getAndMerge(ctx context.Context, plan *CtyunKafkaUserConfig) (err error) {
	params := &ctgkafka.CtgkafkaSaslUserQueryV3Request{
		RegionId:   plan.RegionId.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		Username:   plan.UserName.ValueString(),
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaSaslUserQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCodeString {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil || resp.ReturnObj.Data == nil {
		err = common.InvalidReturnObjError
		return
	}

	userData := resp.ReturnObj.Data[0]
	// 设置ID
	plan.Id = types.Int32Value(userData.Id)
	// 设置用户名
	plan.UserName = types.StringValue(userData.Username)
	// 设置描述 - 处理空字符串情况
	if userData.Description == "" {
		plan.Description = types.StringNull()
	} else {
		plan.Description = types.StringValue(userData.Description)
	}
	plan.CreateTime = types.StringValue(utils.ConvertToUTCZ(utils.Layout2, userData.Ctime))

	err = c.getAndMergeUserTopicsAcl(ctx, plan)
	if err != nil {
		return
	}
	return
}

// calcPermissionInfo 将types.Set类型的PermissionInfo转换为[]CtyunKafkaAclPermissionInfo
func (c *ctyunKafkaUser) calcPermissionInfo(ctx context.Context, plan *CtyunKafkaUserConfig) (err error) {
	if plan.PermissionInfo.IsUnknown() || plan.PermissionInfo.IsNull() {
		return
	}
	plan.permissionInfoList = []CtyunKafkaAclPermissionInfo{}
	diags := plan.PermissionInfo.ElementsAs(ctx, &plan.permissionInfoList, false)
	if diags.HasError() {
		err = fmt.Errorf(diags.Errors()[0].Detail())
		return
	}
	return
}

// update
func (c *ctyunKafkaUser) updatePermissionInfo(ctx context.Context, plan, state CtyunKafkaUserConfig) (err error) {
	err = c.calcPermissionInfo(ctx, &plan)
	if err != nil {
		return
	}
	err = c.calcPermissionInfo(ctx, &state)
	if err != nil {
		return
	}

	add, del := utils.DifferenceStructArray[CtyunKafkaAclPermissionInfo](plan.permissionInfoList, state.permissionInfoList)
	plan.permissionInfoList = del
	err = c.updateUserTopicsAcl(ctx, plan, "DELETE")
	if err != nil {
		return
	}
	plan.permissionInfoList = add
	err = c.updateUserTopicsAcl(ctx, plan, "CREATE")
	if err != nil {
		return
	}
	return
}
func (c *ctyunKafkaUser) updateUserTopicsAcl(ctx context.Context, plan CtyunKafkaUserConfig, eventType string) (err error) {
	// 构建ACL操作信息列表
	aclOperationInfoList := make([]*ctgkafka.CtgkafkaUpdateUserTopicsAclAclOperationInfoListRequest, 0)

	// 从plan.PermissionInfo获取ACL信息
	if len(plan.permissionInfoList) > 0 {
		// 遍历权限信息列表
		for _, permissionInfo := range plan.permissionInfoList {
			aclInfo := &ctgkafka.CtgkafkaUpdateUserTopicsAclAclOperationInfoListRequest{
				EventType: eventType,
				UserName:  plan.UserName.ValueString(),
				Operation: permissionInfo.Operation.ValueString(),
				Topic:     permissionInfo.Topic.ValueString(),
			}

			// 可选字段处理
			if !permissionInfo.Permission.IsNull() && !permissionInfo.Permission.IsUnknown() {
				aclInfo.Permission = permissionInfo.Permission.ValueString()
			}

			if !permissionInfo.Ip.IsNull() && !permissionInfo.Ip.IsUnknown() {
				aclInfo.Ip = permissionInfo.Ip.ValueString()
			}

			aclOperationInfoList = append(aclOperationInfoList, aclInfo)
		}
	}

	// 如果没有ACL信息要更新，直接返回
	if len(aclOperationInfoList) == 0 {
		return nil
	}

	params := &ctgkafka.CtgkafkaUpdateUserTopicsAclRequest{
		RegionId:             plan.RegionId.ValueString(),
		ProdInstId:           plan.InstanceId.ValueString(),
		AclOperationInfoList: aclOperationInfoList,
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaUpdateUserTopicsAclApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp.StatusCode != common.NormalStatusCodeString {
		return fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj.Data != "update success" {
		return fmt.Errorf("API return error. Data: %s", resp.ReturnObj)
	}
	return
}

func (c *ctyunKafkaUser) getAndMergeUserTopicsAcl(ctx context.Context, plan *CtyunKafkaUserConfig) (err error) {
	var topicsToQuery []string

	// 从现有的PermissionInfo中提取topics（如果有的话）
	if !plan.PermissionInfo.IsNull() && !plan.PermissionInfo.IsUnknown() {
		var permissionInfoList []CtyunKafkaAclPermissionInfo
		diags := plan.PermissionInfo.ElementsAs(ctx, &permissionInfoList, false)
		if diags.HasError() {
			// 处理错误
		} else {
			// 提取所有唯一的topic名称
			topicMap := make(map[string]bool)
			for _, info := range permissionInfoList {
				if !info.Topic.IsNull() && !info.Topic.IsUnknown() {
					topicMap[info.Topic.ValueString()] = true
				}
			}

			for topic := range topicMap {
				topicsToQuery = append(topicsToQuery, topic)
			}
		}
	}
	if len(topicsToQuery) == 0 {
		// 如果没有topics要查询，确保PermissionInfo被正确初始化
		plan.PermissionInfo = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"permission": types.StringType,
				"ip":         types.StringType,
				"operation":  types.StringType,
				"topic":      types.StringType,
			},
		})
		return
	}

	// 分别查询READ和WRITE操作的ACL信息
	readAclList, err := c.queryUserTopicsAcl(ctx, plan, "READ", topicsToQuery[0])
	if err != nil {
		return err
	}

	writeAclList, err := c.queryUserTopicsAcl(ctx, plan, "WRITE", topicsToQuery[0])
	if err != nil {
		return err
	}

	// 合并两种操作的ACL信息
	allAclList := append(readAclList, writeAclList...)

	// 定义对象属性类型（必须与schema中定义的一致）
	objectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"permission": types.StringType,
			"ip":         types.StringType,
			"operation":  types.StringType,
			"topic":      types.StringType,
		},
	}
	// 创建Object列表
	objectList := make([]attr.Value, 0, len(allAclList))
	for _, aclData := range allAclList {
		// 创建属性值映射
		attrValues := map[string]attr.Value{
			"permission": types.StringValue(aclData.Permission),
			"ip":         types.StringValue(aclData.Ip),
			"operation":  types.StringValue(aclData.Operation),
			"topic":      types.StringValue(aclData.Topic),
		}

		// 创建Object值
		object, diags := types.ObjectValue(objectType.AttrTypes, attrValues)
		if diags.HasError() {
			return fmt.Errorf("failed to create object value: %v", diags)
		}

		objectList = append(objectList, object)
	}

	// 创建Set值，使用正确的ObjectType
	setValue, diags := types.SetValue(objectType, objectList)
	if diags.HasError() {
		return fmt.Errorf("failed to convert permission info list to set: %v", diags)
	}

	plan.PermissionInfo = setValue

	return nil
}

// queryUserTopicsAcl 查询指定操作类型的用户ACL信息
func (c *ctyunKafkaUser) queryUserTopicsAcl(ctx context.Context, plan *CtyunKafkaUserConfig, operation string, topic string) ([]*ctgkafka.CtgkafkaSaslUserTopicsAclReturnObjDataResponse, error) {
	params := &ctgkafka.CtgkafkaSaslUserTopicsAclRequest{
		RegionId:   plan.RegionId.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		Username:   plan.UserName.ValueString(),
		Operation:  operation,
		Topic:      topic,
	}

	resp, err := c.meta.Apis.SdkKafkaApis.CtgkafkaSaslUserTopicsAclApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != common.NormalStatusCodeString {
		return nil, fmt.Errorf("API return error. Message: %s", resp.Message)
	} else if resp.ReturnObj == nil || resp.ReturnObj.Data == nil {
		return nil, common.InvalidReturnObjError
	}

	return resp.ReturnObj.Data, nil
}

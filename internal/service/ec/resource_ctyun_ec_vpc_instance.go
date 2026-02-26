package ec

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net"
	"strings"
	"time"
)

type CtyunExpressConnectVpcInstance struct {
	meta *common.CtyunMetadata
	name string
}

func (c *CtyunExpressConnectVpcInstance) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ec_vpc_instance"
	c.name = response.TypeName
}

func (c *CtyunExpressConnectVpcInstance) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta

}

func NewCtyunExpressConnectVpcInstance() resource.Resource {
	return &CtyunExpressConnectVpcInstance{}
}

func (c *CtyunExpressConnectVpcInstance) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := fmt.Sprintf("%s导入实例: %s 失败：%s", c.name, request.ID, err.Error())
			detail := fmt.Sprintf("导入命令：terraform import [%s].[导入配置名称] [id],[ec_id],[cgw_id]", c.name)
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunExpressConnectVpcInstanceConfig

	var ID, ecID, cgwID string

	err = terraform_extend.Split(request.ID, &ID, &ecID, &cgwID)
	if err != nil {
		return
	}

	if ID == "" {
		err = fmt.Errorf("id不能为空")
		return
	}
	if ecID == "" {
		err = fmt.Errorf("ec_id不能为空")
		return
	}
	if cgwID == "" {
		err = fmt.Errorf("cgw_id不能为空")
		return
	}

	config.ID = types.StringValue(ID)
	config.EcID = types.StringValue(ecID)
	config.CgwID = types.StringValue(cgwID)

	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunExpressConnectVpcInstance) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("管理云间高速VPC网络实例", "云间高速（标准版）（CT-EC, Express Connect Standard）", "https://www.ctyun.cn/document/10026763/10038256"),
		Attributes: map[string]schema.Attribute{
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速实例ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cgw_id": schema.StringAttribute{
				Required:    true,
				Description: "云网关实例ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rtb_id": schema.StringAttribute{
				Required:    true,
				Description: "路由表ID",
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
				Description: "资源池id,如果不填这默认使用provider ctyun中的region_id或者环境变量",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "vpc ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"exclusive_id": schema.StringAttribute{
				Optional:    true,
				Description: "专属云资源池ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"route_learn": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(1),
				Description: "路由学习开关，开启后云网关自动学习网络实例路由, 1-开启，0不开启",
				Validators: []validator.Int32{
					int32validator.OneOf(1, 0),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
					int32planmodifier.RequiresReplace(),
				},
			},
			"route_sync": schema.Int32Attribute{
				Optional:    true,
				Computed:    true,
				Description: "路由同步开关，开启后云网关路由自动同步到网络实例,1-开启同步，0-不开启同步",
				Default:     int32default.StaticInt32(1),
				Validators: []validator.Int32{
					int32validator.OneOf(1, 0),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
					int32planmodifier.RequiresReplace(),
				},
			},
			"subnets": schema.SetAttribute{
				Required:    true,
				Description: "subnet id列表 支持更新",
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "vpc网络实例id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *CtyunExpressConnectVpcInstance) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunExpressConnectVpcInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.flushVpcRoute(ctx, plan)
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

func (c *CtyunExpressConnectVpcInstance) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunExpressConnectVpcInstanceConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if errors.Is(err, common.ResourceNotExistError) {
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

func (c *CtyunExpressConnectVpcInstance) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// 读取tf文件中配置

	var plan CtyunExpressConnectVpcInstanceConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 读取state中的配置
	var state CtyunExpressConnectVpcInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.update(ctx, &state, &plan)
	if err != nil {
		return
	}
	err = c.flushVpcRoute(ctx, plan)
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

func (c *CtyunExpressConnectVpcInstance) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunExpressConnectVpcInstanceConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunExpressConnectVpcInstance) create(ctx context.Context, config *CtyunExpressConnectVpcInstanceConfig) error {

	// todo cgw 和vpc和subnet建议做个判断兜底
	params := &ec.EcEcAddVPCNetworkRequest{
		EcID:       config.EcID.ValueString(),
		CgwID:      config.CgwID.ValueString(),
		RtbID:      config.RtbID.ValueString(),
		DcID:       config.RegionID.ValueString(),
		VpcID:      config.VpcID.ValueString(),
		RouteLearn: config.RouteLearn.ValueInt32Pointer(),
		RouteSync:  config.RouteSync.ValueInt32Pointer(),
	}
	if !config.ExclusiveID.IsNull() {
		params.ExclusiveID = config.ExclusiveID.ValueStringPointer()
	}
	var subnetIds []string
	diags := config.Subnets.ElementsAs(ctx, &subnetIds, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	var subnets []*ec.EcEcAddVPCNetworkSubnetsRequest
	for _, subnetId := range subnetIds {
		var subnet ec.EcEcAddVPCNetworkSubnetsRequest
		subnet.SubnetID = subnetId
		name, ipVersion, cidr, err := c.getSubnetInfoByID(ctx, subnetId, config)
		if err != nil {
			return err
		}
		subnet.SubnetName = name
		subnet.IPVersion = ipVersion
		subnet.CIDR = cidr
		subnets = append(subnets, &subnet)
	}
	params.Subnets = subnets
	resp, err := c.meta.Apis.SdkEcApis.EcEcAddVPCNetworkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建VPC网络实例失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", *resp.ErrorCode, *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	// 通过vpc id + ec id + cwg id 去获取id
	err = c.getInstanceIdByEcVpcInfo(ctx, config)
	if err != nil {
		return err
	}
	return nil
}

// flushVpcRoute 刷新路由
func (c *CtyunExpressConnectVpcInstance) flushVpcRoute(ctx context.Context, config CtyunExpressConnectVpcInstanceConfig) error {
	params := &ec.EcEcUpdateInstanceRouteRequest{
		EcID:  config.EcID.ValueString(),
		VpcID: config.VpcID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkEcApis.EcEcUpdateInstanceRouteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("刷新路由失败")
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", *resp.ErrorCode, *resp.Message)
		return err
	}
	return nil
}

func (c *CtyunExpressConnectVpcInstance) getAndMerge(ctx context.Context, config *CtyunExpressConnectVpcInstanceConfig) error {
	params := &ec.EcEcListVPCNetworkRequest{
		EcID:       config.EcID.ValueString(),
		InstanceID: config.ID.ValueStringPointer(),
		CgwID:      config.CgwID.ValueStringPointer(),
	}
	resp, err := c.meta.Apis.SdkEcApis.EcEcListVPCNetworkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("获取VPC网络实例失败（id=%s），接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", *resp.ErrorCode, *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}

	if len(resp.ReturnObj.Results) > 1 {
		err = fmt.Errorf("通过id=%s查询，获取到多个VPC网络实例，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	}
	if len(resp.ReturnObj.Results) == 0 {
		err = common.ResourceNotExistError
		return err
	}
	result := resp.ReturnObj.Results[0]
	config.EcID = types.StringValue(*result.EcID)
	config.CgwID = types.StringValue(*result.CgwID)
	config.RtbID = types.StringValue(*result.RtbID)
	config.RegionID = types.StringValue(*result.DcID)
	config.VpcID = types.StringValue(*result.VpcID)
	if !config.ExclusiveID.IsNull() {
		config.ExclusiveID = types.StringValue(*result.ExclusiveID)
	}
	config.RouteLearn = types.Int32Value(*result.RouteLearn)
	config.RouteSync = types.Int32Value(*result.RouteSync)

	// 更新subnet
	// 更新subnet - 修改后的代码
	var subnets []string
	subnetMap := make(map[string]bool) // 用于去重的map
	for _, subnet := range result.SubnetList {
		subnetID := *subnet.SubnetID
		// 检查是否已经存在，避免重复添加
		if !subnetMap[subnetID] {
			subnetMap[subnetID] = true
			subnets = append(subnets, subnetID)
		}
	}
	subnetTmp, diags := types.SetValueFrom(ctx, types.StringType, subnets)
	if diags.HasError() {
		err = fmt.Errorf(diags[0].Detail())
		return err
	}
	config.Subnets = subnetTmp

	return nil

}

func (c *CtyunExpressConnectVpcInstance) update(ctx context.Context, state *CtyunExpressConnectVpcInstanceConfig, plan *CtyunExpressConnectVpcInstanceConfig) error {
	if plan.Subnets.Equal(state.Subnets) {
		return nil
	}
	params := &ec.EcEcUpdateVPCNetworkRequest{
		EcID:  state.EcID.ValueString(),
		VpcID: state.VpcID.ValueString(),
	}
	var subnetIds []string
	diags := plan.Subnets.ElementsAs(ctx, &subnetIds, false)
	if diags.HasError() {
		err := fmt.Errorf(diags[0].Detail())
		return err
	}
	var subnets []*ec.EcEcUpdateVPCNetworkSubnetsRequest
	for _, subnetId := range subnetIds {
		var subnet ec.EcEcUpdateVPCNetworkSubnetsRequest
		subnet.SubnetID = subnetId
		name, ipVersion, cidr, err := c.getSubnetInfoByID(ctx, subnetId, state)
		if err != nil {
			return err
		}
		subnet.SubnetName = name
		subnet.IPVersion = ipVersion
		subnet.CIDR = cidr
		subnets = append(subnets, &subnet)
	}
	params.Subnets = subnets
	resp, err := c.meta.Apis.SdkEcApis.EcEcUpdateVPCNetworkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("更新vpc网络实例(id=$%s)失败，接口返回nil，请联系研发确认问题原因！", state.ID.ValueString())
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", *resp.ErrorCode, *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}

	err = c.updateLoop(ctx, plan, len(subnets), 60)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunExpressConnectVpcInstance) delete(ctx context.Context, config CtyunExpressConnectVpcInstanceConfig) error {
	params := &ec.EcEcDeleteVPCNetworkRequest{
		VpcID: config.VpcID.ValueString(),
		EcID:  config.EcID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkEcApis.EcEcDeleteVPCNetworkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除VPC网络实例(vpc_id=%s, ec_id=%s,id=%s)失败，接口返回nil，请联系研发确认问题原因！", config.VpcID.ValueString(),
			config.EcID.ValueString(), config.ID.ValueString())
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", *resp.ErrorCode, *resp.Message)
		return err
	}
	err = c.deleteLoop(ctx, config)
	if err != nil {
		return err
	}
	return nil
}

func (c *CtyunExpressConnectVpcInstance) getSubnetInfoByID(ctx context.Context, subnetId string, config *CtyunExpressConnectVpcInstanceConfig) (string, string, string, error) {

	params := &ctvpc.CtvpcListSubnetRequest{
		RegionID: config.RegionID.ValueString(),
		VpcID:    config.VpcID.ValueStringPointer(),
		SubnetID: &subnetId,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListSubnetApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return "", "", "", err
	} else if resp.StatusCode == common.ErrorStatusCode {

		return "", "", "", fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
	} else if resp.ReturnObj == nil || len(resp.ReturnObj.Subnets) != 1 {

		return "", "", "", common.InvalidReturnObjError
	}
	subnet := resp.ReturnObj.Subnets[0] // Use the first result

	ipVersion, err2 := c.checkIpVersion(*subnet.CIDR)
	if err2 != nil {
		return "", "", "", err
	}
	return *subnet.Name, strings.ToUpper(ipVersion), *subnet.CIDR, nil // Changed from subnet.CIDR to *subnet.CIDR
}

func (c *CtyunExpressConnectVpcInstance) checkIpVersion(cidr string) (string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}

	if ipNet.IP.To4() != nil {
		return "IPv4", nil
	} else {
		return "IPv6", nil
	}

}

func (c *CtyunExpressConnectVpcInstance) getInstanceIdByEcVpcInfo(ctx context.Context, config *CtyunExpressConnectVpcInstanceConfig, loopCount ...int) error {

	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			vpcInstances, err2 := c.getVpcInstances(ctx, config, false)
			if err2 != nil {
				err = err2
				return false
			}
			for _, instance := range vpcInstances {
				if *instance.VpcID == config.VpcID.ValueString() && *instance.EcID == config.EcID.ValueString() && *instance.CgwID == config.CgwID.ValueString() {
					{
						if *instance.InstanceID == "" {
							return true
						}
						config.ID = types.StringValue(*instance.InstanceID)
						return false
					}
				}
			}
			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return fmt.Errorf("根据ec_id=%s，vpc_id=%s和cgw_id=%s查询vpc实例，并未查询到任何相关信息！", config.EcID.ValueString(), config.VpcID.ValueString(), config.CgwID.ValueString())
	}
	return err
}

func (c *CtyunExpressConnectVpcInstance) getVpcInstances(ctx context.Context, config *CtyunExpressConnectVpcInstanceConfig, needID bool) ([]*ec.EcEcListVPCNetworkReturnObjResultsResponse, error) {
	params := &ec.EcEcListVPCNetworkRequest{
		EcID:  config.EcID.ValueString(),
		CgwID: config.CgwID.ValueStringPointer(),
	}
	if needID {
		params.InstanceID = config.ID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkEcApis.EcEcListVPCNetworkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return nil, err
	} else if resp == nil {
		err = fmt.Errorf("获取VPC网络实例失败（id=%s），接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return nil, err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", *resp.ErrorCode, *resp.Message)
		return nil, err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return nil, err
	}
	return resp.ReturnObj.Results, nil
}

func (c *CtyunExpressConnectVpcInstance) deleteLoop(ctx context.Context, config CtyunExpressConnectVpcInstanceConfig, loopCount ...int) error {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*30, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			vpcInstances, err2 := c.getVpcInstances(ctx, &config, false)
			if err2 != nil {
				err = err2
				return false
			}
			if len(vpcInstances) == 0 {
				return false
			}
			for _, instance := range vpcInstances {
				if *instance.VpcID == config.VpcID.ValueString() && *instance.EcID == config.EcID.ValueString() && *instance.CgwID == config.CgwID.ValueString() {
					{
						return true
					}
				}
			}
			return false
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return fmt.Errorf("id=%s的vpc实例仍未删除成功", config.ID.ValueString())
	}
	return err
}

func (c *CtyunExpressConnectVpcInstance) updateLoop(ctx context.Context, config *CtyunExpressConnectVpcInstanceConfig, subnetNum int, loopCount ...int) error {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*10, count)
	if err != nil {
		return err
	}
	result := retryer.Start(
		func(currentTime int) bool {
			vpcInstances, err2 := c.getVpcInstances(ctx, config, false)
			if err2 != nil {
				err = err2
				return false
			}
			for _, instance := range vpcInstances {
				if *instance.VpcID == config.VpcID.ValueString() && *instance.EcID == config.EcID.ValueString() && *instance.CgwID == config.CgwID.ValueString() {
					{
						if len(instance.SubnetList) == subnetNum {
							return false
						}
					}
				}
			}

			return true
		},
	)
	if result.ReturnReason == business.ReachMaxLoopTime {
		return fmt.Errorf("id=%s的vpc实例子网仍未删除成功", config.ID.ValueString())
	}
	return err
}

type CtyunExpressConnectVpcInstanceConfig struct {
	EcID     types.String `tfsdk:"ec_id"`
	CgwID    types.String `tfsdk:"cgw_id"`
	RtbID    types.String `tfsdk:"rtb_id"`
	RegionID types.String `tfsdk:"region_id"`

	VpcID       types.String `tfsdk:"vpc_id"`
	ExclusiveID types.String `tfsdk:"exclusive_id"`
	RouteLearn  types.Int32  `tfsdk:"route_learn"`
	RouteSync   types.Int32  `tfsdk:"route_sync"`
	Subnets     types.Set    `tfsdk:"subnets"`
	ID          types.String `tfsdk:"id"`
}

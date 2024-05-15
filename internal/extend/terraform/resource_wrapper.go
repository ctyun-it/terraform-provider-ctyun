package terraform

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func WrapResource(t resource.Resource, hooks *AopAdvices) func() resource.Resource {
	wrapper := &ResourceWrapper{
		target: t,
		hooks:  hooks,
	}
	return func() resource.Resource {
		return wrapper
	}
}

type ResourceWrapper struct {
	resource.Resource
	resource.ResourceWithConfigure
	resource.ResourceWithConfigValidators
	resource.ResourceWithImportState
	resource.ResourceWithModifyPlan
	resource.ResourceWithUpgradeState
	resource.ResourceWithValidateConfig

	target resource.Resource
	hooks  *AopAdvices
}

func (r *ResourceWrapper) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	r.target.Metadata(ctx, request, response)
}

func (r *ResourceWrapper) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	r.target.Schema(ctx, request, response)
}
func (r *ResourceWrapper) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var advices []Advice[resource.CreateRequest, *resource.CreateResponse]
	advices = append(advices, r.hooks.ResourceCreateAopApi...)
	advices = append(
		advices,
		&adviceWrapper[resource.CreateRequest, *resource.CreateResponse]{
			Target: r.target.Create,
		},
	)
	chain := ResourceDecoratorChain[resource.CreateRequest, *resource.CreateResponse]{
		Chains: advices,
	}
	err := chain.Next(ctx, request, response)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
	}
}

func (r *ResourceWrapper) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var advices []Advice[resource.ReadRequest, *resource.ReadResponse]
	advices = append(advices, r.hooks.ResourceReadAopApi...)
	advices = append(
		advices,
		&adviceWrapper[resource.ReadRequest, *resource.ReadResponse]{
			Target: r.target.Read,
		},
	)
	chain := ResourceDecoratorChain[resource.ReadRequest, *resource.ReadResponse]{
		Chains: advices,
	}
	err := chain.Next(ctx, request, response)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
	}
}

func (r *ResourceWrapper) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var advices []Advice[resource.UpdateRequest, *resource.UpdateResponse]
	advices = append(advices, r.hooks.ResourceUpdateAopApi...)
	advices = append(
		advices,
		&adviceWrapper[resource.UpdateRequest, *resource.UpdateResponse]{
			Target: r.target.Update,
		},
	)
	chain := ResourceDecoratorChain[resource.UpdateRequest, *resource.UpdateResponse]{
		Chains: advices,
	}
	err := chain.Next(ctx, request, response)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
	}
}

func (r *ResourceWrapper) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var advices []Advice[resource.DeleteRequest, *resource.DeleteResponse]
	advices = append(advices, r.hooks.ResourceDeleteAopApi...)
	advices = append(
		advices,
		&adviceWrapper[resource.DeleteRequest, *resource.DeleteResponse]{
			Target: r.target.Delete,
		},
	)
	chain := ResourceDecoratorChain[resource.DeleteRequest, *resource.DeleteResponse]{
		Chains: advices,
	}
	err := chain.Next(ctx, request, response)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
	}
}

func (r *ResourceWrapper) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	configure, ok := interface{}(r.target).(resource.ResourceWithConfigure)
	if ok {
		configure.Configure(ctx, request, response)
	}
}

func (r *ResourceWrapper) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	configure, ok := interface{}(r.target).(resource.ResourceWithConfigValidators)
	if ok {
		return configure.ConfigValidators(ctx)
	}
	return nil
}

func (r *ResourceWrapper) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	configure, ok := interface{}(r.target).(resource.ResourceWithImportState)
	if ok {
		configure.ImportState(ctx, request, response)
	}
}

func (r *ResourceWrapper) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	configure, ok := interface{}(r.target).(resource.ResourceWithModifyPlan)
	if ok {
		configure.ModifyPlan(ctx, request, response)
	}
}

func (r *ResourceWrapper) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	configure, ok := interface{}(r.target).(resource.ResourceWithUpgradeState)
	if ok {
		return configure.UpgradeState(ctx)
	}
	return make(map[int64]resource.StateUpgrader)
}

func (r *ResourceWrapper) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	configure, ok := interface{}(r.target).(resource.ResourceWithValidateConfig)
	if ok {
		configure.ValidateConfig(ctx, request, response)
	}
}

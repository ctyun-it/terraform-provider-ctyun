package terraform

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func WrapDataSource(t datasource.DataSource, hooks *AopAdvices) func() datasource.DataSource {
	facade := &DatasourceWrapper{
		target: t,
		hooks:  hooks,
	}
	return func() datasource.DataSource {
		return facade
	}
}

type DatasourceWrapper struct {
	datasource.DataSource
	datasource.DataSourceWithConfigure
	datasource.DataSourceWithConfigValidators
	datasource.DataSourceWithValidateConfig

	target datasource.DataSource
	hooks  *AopAdvices
}

func (w *DatasourceWrapper) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	w.target.Metadata(ctx, request, response)
}

func (w *DatasourceWrapper) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	w.target.Schema(ctx, request, response)
}

func (w *DatasourceWrapper) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var advices []Advice[datasource.ReadRequest, *datasource.ReadResponse]
	advices = append(advices, w.hooks.DataSourceReadAopApi...)
	advices = append(
		advices,
		&adviceWrapper[datasource.ReadRequest, *datasource.ReadResponse]{
			Target: w.target.Read,
		},
	)
	chain := ResourceDecoratorChain[datasource.ReadRequest, *datasource.ReadResponse]{
		Chains: advices,
	}
	err := chain.Next(ctx, request, response)
	if err != nil {
		response.Diagnostics.AddError(err.Error(), err.Error())
	}
}

func (w *DatasourceWrapper) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	configure, ok := interface{}(w.target).(datasource.DataSourceWithConfigure)
	if ok {
		configure.Configure(ctx, request, response)
	}
}

func (w *DatasourceWrapper) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	configure, ok := interface{}(w.target).(datasource.DataSourceWithConfigValidators)
	if ok {
		return configure.ConfigValidators(ctx)
	}
	return nil
}

func (w *DatasourceWrapper) ValidateConfig(ctx context.Context, request datasource.ValidateConfigRequest, response *datasource.ValidateConfigResponse) {
	configure, ok := interface{}(w.target).(datasource.DataSourceWithValidateConfig)
	if ok {
		configure.ValidateConfig(ctx, request, response)
	}
}

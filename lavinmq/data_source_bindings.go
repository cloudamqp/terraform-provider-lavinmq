package lavinmq

import (
	"context"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &bindingsDataSource{}
	_ datasource.DataSourceWithConfigure = &bindingsDataSource{}
)

func NewBindingsDataSource() datasource.DataSource {
	return &bindingsDataSource{}
}

type bindingsDataSource struct {
	services *clientlibrary.Services
}

type bindingDataSourceModel struct {
	Source          types.String `tfsdk:"source"`
	Vhost           types.String `tfsdk:"vhost"`
	Destination     types.String `tfsdk:"destination"`
	DestinationType types.String `tfsdk:"destination_type"`
	RoutingKey      types.String `tfsdk:"routing_key"`
	PropertiesKey   types.String `tfsdk:"properties_key"`
}

type bindingsDataSourceModel struct {
	Vhost    types.String             `tfsdk:"vhost"`
	Bindings []bindingDataSourceModel `tfsdk:"bindings"`
}

func (d *bindingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bindings"
}

func (d *bindingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List bindings in a vhost.",
		Attributes: map[string]schema.Attribute{
			"vhost": schema.StringAttribute{
				Description: "The vhost to list bindings from.",
				Optional:    true,
			},
			"bindings": schema.ListNestedAttribute{
				Description: "List of bindings in the vhost.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source": schema.StringAttribute{
							Description: "The source exchange name.",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "The vhost the binding is located in.",
							Computed:    true,
						},
						"destination": schema.StringAttribute{
							Description: "The destination queue or exchange name.",
							Computed:    true,
						},
						"destination_type": schema.StringAttribute{
							Description: "The destination type: 'queue' or 'exchange'.",
							Computed:    true,
						},
						"routing_key": schema.StringAttribute{
							Description: "The routing key for the binding.",
							Computed:    true,
						},
						"properties_key": schema.StringAttribute{
							Description: "Unique properties key for this binding.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *bindingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.services = req.ProviderData.(*clientlibrary.Services)
}

func (d *bindingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config bindingsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bindings, err := d.services.Bindings.List(ctx, config.Vhost.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve bindings", err.Error())
		return
	}
	if len(bindings) == 0 {
		tflog.Warn(ctx, "No bindings found")
	}

	var state bindingsDataSourceModel
	state.Vhost = config.Vhost
	state.Bindings = []bindingDataSourceModel{}

	for _, binding := range bindings {
		state.Bindings = append(state.Bindings, bindingDataSourceModel{
			Source:          types.StringValue(binding.Source),
			Vhost:           types.StringValue(binding.Vhost),
			Destination:     types.StringValue(binding.Destination),
			DestinationType: types.StringValue(binding.DestinationType),
			RoutingKey:      types.StringValue(binding.RoutingKey),
			PropertiesKey:   types.StringValue(binding.PropertiesKey),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

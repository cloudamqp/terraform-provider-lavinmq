package lavinmq

import (
	"context"
	"encoding/json"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &federationUpstreamSetsDataSource{}
	_ datasource.DataSourceWithConfigure = &federationUpstreamSetsDataSource{}
)

func NewFederationUpstreamSetsDataSource() datasource.DataSource {
	return &federationUpstreamSetsDataSource{}
}

type federationUpstreamSetsDataSource struct {
	services *clientlibrary.Services
}

type federationUpstreamSetsDataSourceModel struct {
	Vhost                  types.String                           `tfsdk:"vhost"`
	FederationUpstreamSets []federationUpstreamSetDataSourceModel `tfsdk:"federation_upstream_sets"`
}

type federationUpstreamSetDataSourceModel struct {
	Name      types.String `tfsdk:"name"`
	Vhost     types.String `tfsdk:"vhost"`
	Upstreams types.List   `tfsdk:"upstreams"`
}

func (d *federationUpstreamSetsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_federation_upstream_sets"
}

func (d *federationUpstreamSetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vhost": schema.StringAttribute{
				Description: "The vhost to list federation upstream sets from. If not specified, lists all federation upstream sets.",
				Optional:    true,
			},
			"federation_upstream_sets": schema.ListNestedAttribute{
				Description: "List of federation upstream sets.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the federation upstream set.",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "Virtual host where the federation upstream set is defined.",
							Computed:    true,
						},
						"upstreams": schema.ListAttribute{
							Description: "List of upstream names in this set.",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *federationUpstreamSetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.services = req.ProviderData.(*clientlibrary.Services)
}

func (d *federationUpstreamSetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config federationUpstreamSetsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parameters, err := d.services.Parameters.List(ctx, "federation-upstream-set", config.Vhost.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve federation upstream sets", err.Error())
		return
	}
	if len(parameters) == 0 {
		tflog.Warn(ctx, "No federation upstream sets found")
	}

	var state federationUpstreamSetsDataSourceModel
	state.Vhost = config.Vhost
	state.FederationUpstreamSets = []federationUpstreamSetDataSourceModel{}

	for _, param := range parameters {
		valueBytes, err := json.Marshal(param.Value)
		if err != nil {
			resp.Diagnostics.AddError("Failed to marshal federation upstream set value", err.Error())
			return
		}

		var upstreamItems []clientlibrary.FederationUpstreamSetItem
		if err := json.Unmarshal(valueBytes, &upstreamItems); err != nil {
			resp.Diagnostics.AddError("Failed to unmarshal federation upstream set value", err.Error())
			return
		}

		upstreams := make([]string, len(upstreamItems))
		for i, item := range upstreamItems {
			upstreams[i] = item.Upstream
		}

		upstreamsList, diags := types.ListValueFrom(ctx, types.StringType, upstreams)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		state.FederationUpstreamSets = append(state.FederationUpstreamSets, federationUpstreamSetDataSourceModel{
			Name:      types.StringValue(param.Name),
			Vhost:     types.StringValue(param.Vhost),
			Upstreams: upstreamsList,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

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
	_ datasource.DataSource              = &policiesDataSource{}
	_ datasource.DataSourceWithConfigure = &policiesDataSource{}
)

func NewPoliciesDataSource() datasource.DataSource {
	return &policiesDataSource{}
}

type policiesDataSource struct {
	services *clientlibrary.Services
}

type policiesDataSourceModel struct {
	Policies []policyDataSourceModel `tfsdk:"policies"`
}

type policyDataSourceModel struct {
	Name     types.String `tfsdk:"name"`
	Vhost    types.String `tfsdk:"vhost"`
	Pattern  types.String `tfsdk:"pattern"`
	Priority types.Int64  `tfsdk:"priority"`
	ApplyTo  types.String `tfsdk:"apply_to"`
	// TODO Definition
}

func (d *policiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policies"
}

func (d *policiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"policies": schema.ListNestedAttribute{
				Description: "List of policies.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the policy.",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "Virtual host where the policy is applied.",
							Computed:    true,
						},
						"pattern": schema.StringAttribute{
							Description: "Regular expression pattern that matches the names of exchanges or queues to which the policy applies.",
							Computed:    true,
						},
						"priority": schema.Int64Attribute{
							Description: "Policy priority. Higher numbers indicate higher priority.",
							Computed:    true,
						},
						"apply_to": schema.StringAttribute{
							Description: "What the policy applies to: 'all', 'exchanges', or 'queues'.",
							Computed:    true,
						},
						// TODO: definition
					},
				},
			},
		},
	}
}

func (d *policiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.services = req.ProviderData.(*clientlibrary.Services)
}

func (d *policiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state policiesDataSourceModel

	policies, err := d.services.Policies.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve policies", err.Error())
		return
	}

	if policies == nil {
		tflog.Warn(ctx, "no policies found")
		return
	}

	for _, policy := range policies {
		state.Policies = append(state.Policies, policyDataSourceModel{
			Name:     types.StringValue(policy.Name),
			Vhost:    types.StringValue(policy.Vhost),
			Pattern:  types.StringValue(policy.Pattern),
			Priority: types.Int64Value(int64(policy.Priority)),
			ApplyTo:  types.StringValue(policy.ApplyTo),
			// TODO: Definition
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

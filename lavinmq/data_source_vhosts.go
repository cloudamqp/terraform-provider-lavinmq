package lavinmq

import (
	"context"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &vhostsDataSource{}

func NewVhostDataSource() datasource.DataSource {
	return &vhostsDataSource{}
}

type vhostsDataSource struct {
	services *clientlibrary.Services
}

type vhostDataSourceModel struct {
	Name types.String `tfsdk:"name"`
}

type vhostListDataSourceModel struct {
	Vhosts []vhostDataSourceModel `tfsdk:"vhosts"`
}

func (d *vhostsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vhosts"
}

func (d *vhostsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vhosts": schema.ListNestedAttribute{
				Description: "List of vhosts.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the managed vhost.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *vhostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.services = req.ProviderData.(*clientlibrary.Services)
}

func (d *vhostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vhostListDataSourceModel

	vhosts, err := d.services.Vhosts.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve vhosts", err.Error())
		return
	}

	if vhosts == nil {
		tflog.Warn(ctx, "no vhost found")
		return
	}

	for _, vhost := range vhosts {
		state.Vhosts = append(state.Vhosts, vhostDataSourceModel{
			Name: types.StringValue(vhost.Name),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

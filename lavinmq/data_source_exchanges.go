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
	_ datasource.DataSource              = &exchangesDataSource{}
	_ datasource.DataSourceWithConfigure = &exchangesDataSource{}
)

func NewExchangesDataSource() datasource.DataSource {
	return &exchangesDataSource{}
}

type exchangesDataSource struct {
	services *clientlibrary.Services
}

type exchangesDataSourceModel struct {
	Exchanges []exchangeDataSourceModel `tfsdk:"exchanges"`
}

type exchangeDataSourceModel struct {
	Name       types.String `tfsdk:"name"`
	Vhost      types.String `tfsdk:"vhost"`
	Type       types.String `tfsdk:"type"`
	AutoDelete types.Bool   `tfsdk:"auto_delete"`
	Durable    types.Bool   `tfsdk:"durable"`
}

func (d *exchangesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exchanges"
}

func (d *exchangesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all exchanges.",
		Attributes: map[string]schema.Attribute{
			"exchanges": schema.ListNestedAttribute{
				Description: "List of exchanges.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the exchange.",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "The vhost the exchange is located in.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The exchange type (direct, fanout, topic, headers).",
							Computed:    true,
						},
						"auto_delete": schema.BoolAttribute{
							Description: "Whether the exchange is automatically deleted when no longer used.",
							Computed:    true,
						},
						"durable": schema.BoolAttribute{
							Description: "Whether the exchange should survive a broker restart.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *exchangesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.services = req.ProviderData.(*clientlibrary.Services)
}

func (d *exchangesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state exchangesDataSourceModel

	exchanges, err := d.services.Exchanges.List(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve exchanges", err.Error())
		return
	}

	if exchanges == nil {
		tflog.Warn(ctx, "no exchanges found")
		return
	}

	for _, exchange := range exchanges {
		state.Exchanges = append(state.Exchanges, exchangeDataSourceModel{
			Name:       types.StringValue(exchange.Name),
			Vhost:      types.StringValue(exchange.Vhost),
			Type:       types.StringValue(exchange.Type),
			AutoDelete: types.BoolValue(exchange.AutoDelete),
			Durable:    types.BoolValue(exchange.Durable),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

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
	_ datasource.DataSource              = &shovelsDataSource{}
	_ datasource.DataSourceWithConfigure = &shovelsDataSource{}
)

func NewShovelsDataSource() datasource.DataSource {
	return &shovelsDataSource{}
}

type shovelsDataSource struct {
	services *clientlibrary.Services
}

type shovelsDataSourceModel struct {
	Vhost   types.String            `tfsdk:"vhost"`
	Shovels []shovelDataSourceModel `tfsdk:"shovels"`
}

type shovelDataSourceModel struct {
	Name             types.String `tfsdk:"name"`
	Vhost            types.String `tfsdk:"vhost"`
	SrcURI           types.String `tfsdk:"src_uri"`
	DestURI          types.String `tfsdk:"dest_uri"`
	SrcQueue         types.String `tfsdk:"src_queue"`
	SrcExchange      types.String `tfsdk:"src_exchange"`
	SrcExchangeKey   types.String `tfsdk:"src_exchange_key"`
	DestQueue        types.String `tfsdk:"dest_queue"`
	DestExchange     types.String `tfsdk:"dest_exchange"`
	DestExchangeKey  types.String `tfsdk:"dest_exchange_key"`
	SrcPrefetchCount types.Int64  `tfsdk:"src_prefetch_count"`
	SrcDeleteAfter   types.String `tfsdk:"src_delete_after"`
	ReconnectDelay   types.Int64  `tfsdk:"reconnect_delay"`
	AckMode          types.String `tfsdk:"ack_mode"`
}

func (d *shovelsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shovels"
}

func (d *shovelsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vhost": schema.StringAttribute{
				Description: "The vhost to list shovels from. If not specified, lists all shovels.",
				Optional:    true,
			},
			"shovels": schema.ListNestedAttribute{
				Description: "List of shovels.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the shovel.",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "Virtual host where the shovel is defined.",
							Computed:    true,
						},
						"src_uri": schema.StringAttribute{
							Description: "Source AMQP URI.",
							Computed:    true,
						},
						"dest_uri": schema.StringAttribute{
							Description: "Destination AMQP URI.",
							Computed:    true,
						},
						"src_queue": schema.StringAttribute{
							Description: "Source queue name.",
							Computed:    true,
						},
						"src_exchange": schema.StringAttribute{
							Description: "Source exchange name.",
							Computed:    true,
						},
						"src_exchange_key": schema.StringAttribute{
							Description: "Source exchange routing key.",
							Computed:    true,
						},
						"dest_queue": schema.StringAttribute{
							Description: "Destination queue name.",
							Computed:    true,
						},
						"dest_exchange": schema.StringAttribute{
							Description: "Destination exchange name.",
							Computed:    true,
						},
						"dest_exchange_key": schema.StringAttribute{
							Description: "Destination exchange routing key.",
							Computed:    true,
						},
						"src_prefetch_count": schema.Int64Attribute{
							Description: "Source prefetch count.",
							Computed:    true,
						},
						"src_delete_after": schema.StringAttribute{
							Description: "When to delete messages from source.",
							Computed:    true,
						},
						"reconnect_delay": schema.Int64Attribute{
							Description: "Reconnect delay in seconds.",
							Computed:    true,
						},
						"ack_mode": schema.StringAttribute{
							Description: "Acknowledgement mode.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *shovelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.services = req.ProviderData.(*clientlibrary.Services)
}

func (d *shovelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config shovelsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parameters, err := d.services.Parameters.List(ctx, "shovel", config.Vhost.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve shovels", err.Error())
		return
	}
	if len(parameters) == 0 {
		tflog.Warn(ctx, "No shovels found")
	}

	var state shovelsDataSourceModel
	state.Vhost = config.Vhost
	state.Shovels = []shovelDataSourceModel{}

	for _, param := range parameters {
		valueBytes, err := json.Marshal(param.Value)
		if err != nil {
			resp.Diagnostics.AddError("Failed to marshal shovel value", err.Error())
			return
		}

		var shovelValue clientlibrary.ShovelValue
		if err := json.Unmarshal(valueBytes, &shovelValue); err != nil {
			resp.Diagnostics.AddError("Failed to unmarshal shovel value", err.Error())
			return
		}

		state.Shovels = append(state.Shovels, shovelDataSourceModel{
			Name:             types.StringValue(param.Name),
			Vhost:            types.StringValue(param.Vhost),
			SrcURI:           types.StringValue(shovelValue.SrcURI),
			DestURI:          types.StringValue(shovelValue.DestURI),
			SrcQueue:         types.StringValue(shovelValue.SrcQueue),
			SrcExchange:      types.StringValue(shovelValue.SrcExchange),
			SrcExchangeKey:   types.StringValue(shovelValue.SrcExchangeKey),
			DestQueue:        types.StringValue(shovelValue.DestQueue),
			DestExchange:     types.StringValue(shovelValue.DestExchange),
			DestExchangeKey:  types.StringValue(shovelValue.DestExchangeKey),
			SrcPrefetchCount: types.Int64Value(shovelValue.SrcPrefetchCount),
			SrcDeleteAfter:   types.StringValue(shovelValue.SrcDeleteAfter),
			ReconnectDelay:   types.Int64Value(shovelValue.ReconnectDelay),
			AckMode:          types.StringValue(shovelValue.AckMode),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

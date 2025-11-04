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
	_ datasource.DataSource              = &federationUpstreamsDataSource{}
	_ datasource.DataSourceWithConfigure = &federationUpstreamsDataSource{}
)

func NewFederationUpstreamsDataSource() datasource.DataSource {
	return &federationUpstreamsDataSource{}
}

type federationUpstreamsDataSource struct {
	services *clientlibrary.Services
}

type federationUpstreamsDataSourceModel struct {
	Vhost               types.String                        `tfsdk:"vhost"`
	FederationUpstreams []federationUpstreamDataSourceModel `tfsdk:"federation_upstreams"`
}

type federationUpstreamDataSourceModel struct {
	Name           types.String `tfsdk:"name"`
	Vhost          types.String `tfsdk:"vhost"`
	URI            types.String `tfsdk:"uri"`
	PrefetchCount  types.Int64  `tfsdk:"prefetch_count"`
	ReconnectDelay types.Int64  `tfsdk:"reconnect_delay"`
	AckMode        types.String `tfsdk:"ack_mode"`
	Exchange       types.String `tfsdk:"exchange"`
	MaxHops        types.Int64  `tfsdk:"max_hops"`
	Expires        types.Int64  `tfsdk:"expires"`
	MessageTTL     types.Int64  `tfsdk:"message_ttl"`
	Queue          types.String `tfsdk:"queue"`
	ConsumerTag    types.String `tfsdk:"consumer_tag"`
}

func (d *federationUpstreamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_federation_upstreams"
}

func (d *federationUpstreamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vhost": schema.StringAttribute{
				Description: "The vhost to list federation upstreams from. If not specified, lists all federation upstreams.",
				Optional:    true,
			},
			"federation_upstreams": schema.ListNestedAttribute{
				Description: "List of federation upstreams.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the federation upstream.",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "Virtual host where the federation upstream is defined.",
							Computed:    true,
						},
						"uri": schema.StringAttribute{
							Description: "AMQP URI of the upstream broker.",
							Computed:    true,
						},
						"prefetch_count": schema.Int64Attribute{
							Description: "Number of messages to prefetch from upstream.",
							Computed:    true,
						},
						"reconnect_delay": schema.Int64Attribute{
							Description: "Delay in seconds before reconnecting.",
							Computed:    true,
						},
						"ack_mode": schema.StringAttribute{
							Description: "Acknowledgement mode.",
							Computed:    true,
						},
						"exchange": schema.StringAttribute{
							Description: "Name of upstream exchange to federate from.",
							Computed:    true,
						},
						"max_hops": schema.Int64Attribute{
							Description: "Maximum number of federation hops.",
							Computed:    true,
						},
						"expires": schema.Int64Attribute{
							Description: "Expiry time in milliseconds.",
							Computed:    true,
						},
						"message_ttl": schema.Int64Attribute{
							Description: "Message TTL in milliseconds.",
							Computed:    true,
						},
						"queue": schema.StringAttribute{
							Description: "Name of upstream queue to federate from.",
							Computed:    true,
						},
						"consumer_tag": schema.StringAttribute{
							Description: "Consumer tag for the federation link.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *federationUpstreamsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.services = req.ProviderData.(*clientlibrary.Services)
}

func (d *federationUpstreamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config federationUpstreamsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parameters, err := d.services.Parameters.List(ctx, "federation-upstream", config.Vhost.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve federation upstreams", err.Error())
		return
	}
	if len(parameters) == 0 {
		tflog.Warn(ctx, "No federation upstreams found")
	}

	var state federationUpstreamsDataSourceModel
	state.Vhost = config.Vhost
	state.FederationUpstreams = []federationUpstreamDataSourceModel{}

	for _, param := range parameters {
		valueBytes, err := json.Marshal(param.Value)
		if err != nil {
			resp.Diagnostics.AddError("Failed to marshal federation upstream value", err.Error())
			return
		}

		var federationValue clientlibrary.FederationUpstreamValue
		if err := json.Unmarshal(valueBytes, &federationValue); err != nil {
			resp.Diagnostics.AddError("Failed to unmarshal federation upstream value", err.Error())
			return
		}

		state.FederationUpstreams = append(state.FederationUpstreams, federationUpstreamDataSourceModel{
			Name:           types.StringValue(param.Name),
			Vhost:          types.StringValue(param.Vhost),
			URI:            types.StringValue(federationValue.URI),
			PrefetchCount:  types.Int64Value(federationValue.PrefetchCount),
			ReconnectDelay: types.Int64Value(federationValue.ReconnectDelay),
			AckMode:        types.StringValue(federationValue.AckMode),
			Exchange:       types.StringValue(federationValue.Exchange),
			MaxHops:        types.Int64Value(federationValue.MaxHops),
			Expires:        types.Int64Value(federationValue.Expires),
			MessageTTL:     types.Int64Value(federationValue.MessageTTL),
			Queue:          types.StringValue(federationValue.Queue),
			ConsumerTag:    types.StringValue(federationValue.ConsumerTag),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

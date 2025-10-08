package lavinmq

import (
	"context"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &queueDataSource{}

func NewQueueDataSource() datasource.DataSource {
	return &queueDataSource{}
}

type queueDataSource struct {
	services *clientlibrary.Services
}

type queueDataSourceModel struct {
	Name       types.String `tfsdk:"name"`
	Vhost      types.String `tfsdk:"vhost"`
	AutoDelete types.Bool   `tfsdk:"auto_delete"`
	Durable    types.Bool   `tfsdk:"durable"`
}

type queuesDataSourceModel struct {
	Vhost  types.String           `tfsdk:"vhost"`
	Queues []queueDataSourceModel `tfsdk:"queues"`
}

func (d *queueDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_queues"
}

func (d *queueDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List queues in a vhost.",
		Attributes: map[string]schema.Attribute{
			"vhost": schema.StringAttribute{
				Description: "The vhost to list queues from.",
				Optional:    true,
			},
			"queues": schema.ListNestedAttribute{
				Description: "List of queues in the vhost.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the queue.",
							Computed:    true,
						},
						"vhost": schema.StringAttribute{
							Description: "The vhost the queue is located in.",
							Computed:    true,
						},
						"auto_delete": schema.BoolAttribute{
							Description: "Whether the queue is automatically deleted when no longer used.",
							Computed:    true,
						},
						"durable": schema.BoolAttribute{
							Description: "Whether the queue should survive a broker restart.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *queueDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.services = req.ProviderData.(*clientlibrary.Services)
}

func (d *queueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config queuesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	queues, err := d.services.Queues.List(ctx, config.Vhost.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve queues", err.Error())
		return
	}

	if queues == nil {
		tflog.Warn(ctx, "no queues found")
		return
	}

	var state queuesDataSourceModel
	state.Vhost = config.Vhost

	for _, queue := range queues {
		state.Queues = append(state.Queues, queueDataSourceModel{
			Name:       types.StringValue(queue.Name),
			Vhost:      types.StringValue(queue.Vhost),
			AutoDelete: types.BoolValue(queue.AutoDelete),
			Durable:    types.BoolValue(queue.Durable),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

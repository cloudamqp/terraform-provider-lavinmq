package lavinmq

import (
	"context"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &publishMessageResource{}
	_ resource.ResourceWithConfigure = &publishMessageResource{}
)

func NewPublishMessageResource() resource.Resource {
	return &publishMessageResource{}
}

type publishMessageResource struct {
	services *clientlibrary.Services
}

type publishMessageResourceModel struct {
	Vhost           types.String `tfsdk:"vhost"`
	Exchange        types.String `tfsdk:"exchange"`
	RoutingKey      types.String `tfsdk:"routing_key"`
	Payload         types.String `tfsdk:"payload"`
	PayloadEncoding types.String `tfsdk:"payload_encoding"`
	Properties      types.Map    `tfsdk:"properties"`
}

func (r *publishMessageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_publish_message"
}

func (r *publishMessageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Publishes a message to an exchange. This is a one-time action resource.",
		Attributes: map[string]schema.Attribute{
			"vhost": schema.StringAttribute{
				Description: "The vhost containing the exchange.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"exchange": schema.StringAttribute{
				Description: "The exchange to publish the message to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"routing_key": schema.StringAttribute{
				Description: "The routing key for the message.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"payload": schema.StringAttribute{
				Description: "The message payload.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"payload_encoding": schema.StringAttribute{
				Description: "The encoding of the payload (e.g., 'string', 'base64'). Defaults to 'string'.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"properties": schema.MapAttribute{
				Description: "Message properties (headers, delivery mode, etc).",
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *publishMessageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	services, ok := req.ProviderData.(*clientlibrary.Services)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *clientlibrary.Services type for provider data.",
		)
		return
	}

	r.services = services
}

func (r *publishMessageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan publishMessageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	properties := make(map[string]any)
	if !plan.Properties.IsNull() {
		propertiesMap := make(map[string]string)
		resp.Diagnostics.Append(plan.Properties.ElementsAs(ctx, &propertiesMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for k, v := range propertiesMap {
			properties[k] = v
		}
	}

	request := clientlibrary.PublishRequest{
		RoutingKey: plan.RoutingKey.ValueString(),
		Payload:    plan.Payload.ValueString(),
		Properties: properties,
	}

	if plan.PayloadEncoding.IsNull() {
		request.PayloadEncoding = "string"
	}

	err := r.services.Messages.Publish(ctx, plan.Vhost.ValueString(), plan.Exchange.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Error publishing message", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *publishMessageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// This resource does not implement the Read function
}

func (r *publishMessageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This resource does not implement the Update function
}

func (r *publishMessageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource does not implement the Delete function
}

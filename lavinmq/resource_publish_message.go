package lavinmq

import (
	"context"
	"math/big"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	Vhost                 types.String  `tfsdk:"vhost"`
	Exchange              types.String  `tfsdk:"exchange"`
	RoutingKey            types.String  `tfsdk:"routing_key"`
	Payload               types.String  `tfsdk:"payload"`
	PayloadEncoding       types.String  `tfsdk:"payload_encoding"`
	Properties            types.Dynamic `tfsdk:"properties"`
	PublishMessageCounter types.Int64   `tfsdk:"publish_message_counter"`
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
			},
			"exchange": schema.StringAttribute{
				Description: "The exchange to publish the message to.",
				Required:    true,
			},
			"routing_key": schema.StringAttribute{
				Description: "The routing key for the message.",
				Required:    true,
			},
			"payload": schema.StringAttribute{
				Description: "The message payload.",
				Required:    true,
			},
			"payload_encoding": schema.StringAttribute{
				Description: "The encoding of the payload (e.g., 'string', 'base64'). Defaults to 'string'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("string"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"properties": schema.DynamicAttribute{
				Description: "Message properties (headers, delivery mode, etc).",
				Optional:    true,
			},
			"publish_message_counter": schema.Int64Attribute{
				Description: "A counter that can be used to trigger a resource update.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
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

	request, diags := r.populateRequest(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
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
	var plan publishMessageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, diags := r.populateRequest(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
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

func (r *publishMessageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource does not implement the Delete function
}

func (r *publishMessageResource) populateRequest(plan publishMessageResourceModel) (clientlibrary.PublishRequest, diag.Diagnostics) {
	propertiesMap := make(map[string]any)
	if !plan.Properties.IsNull() && !plan.Properties.IsUnknown() {
		switch v := plan.Properties.UnderlyingValue().(type) {
		case types.Object:
			for key, value := range v.Attributes() {
				switch val := value.(type) {
				case types.String:
					propertiesMap[key] = val.ValueString()
				case types.Bool:
					propertiesMap[key] = val.ValueBool()
				case types.Number:
					if bigFloat := val.ValueBigFloat(); bigFloat != nil {
						if intVal, accuracy := bigFloat.Int64(); accuracy == big.Exact {
							propertiesMap[key] = intVal
						} else if floatVal, accuracy := bigFloat.Float64(); accuracy == big.Exact {
							propertiesMap[key] = floatVal
						}
					}
				}
			}
		}
	}

	// Default properties values
	if _, ok := propertiesMap["delivery_mode"]; !ok {
		propertiesMap["delivery_mode"] = 2 // persistent
	}
	if _, ok := propertiesMap["content_type"]; !ok {
		propertiesMap["content_type"] = "application/json"
	}

	request := clientlibrary.PublishRequest{
		RoutingKey:      plan.RoutingKey.ValueString(),
		Payload:         plan.Payload.ValueString(),
		PayloadEncoding: plan.PayloadEncoding.ValueString(),
		Properties:      propertiesMap,
	}

	return request, nil
}

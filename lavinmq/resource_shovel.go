package lavinmq

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &shovelResource{}
	_ resource.ResourceWithConfigure   = &shovelResource{}
	_ resource.ResourceWithImportState = &shovelResource{}
)

func NewShovelResource() resource.Resource {
	return &shovelResource{}
}

type shovelResource struct {
	services *clientlibrary.Services
}

type shovelResourceModel struct {
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

func (r *shovelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shovel"
}

func (r *shovelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a shovel for message forwarding between queues and exchanges.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the shovel.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vhost": schema.StringAttribute{
				Description: "Virtual host where the shovel is defined.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"src_uri": schema.StringAttribute{
				Description: "Source AMQP URI (where messages are consumed from).",
				Required:    true,
			},
			"dest_uri": schema.StringAttribute{
				Description: "Destination AMQP URI (where messages are published to).",
				Required:    true,
			},
			"src_queue": schema.StringAttribute{
				Description: "Name of source queue to consume from. Either src_queue or src_exchange must be specified.",
				Optional:    true,
			},
			"src_exchange": schema.StringAttribute{
				Description: "Name of source exchange to consume from. Either src_queue or src_exchange must be specified.",
				Optional:    true,
			},
			"src_exchange_key": schema.StringAttribute{
				Description: "Routing key for source exchange binding (only used with src_exchange).",
				Optional:    true,
			},
			"dest_queue": schema.StringAttribute{
				Description: "Name of destination queue to publish to. Either dest_queue or dest_exchange must be specified.",
				Optional:    true,
			},
			"dest_exchange": schema.StringAttribute{
				Description: "Name of destination exchange to publish to. Either dest_queue or dest_exchange must be specified.",
				Optional:    true,
			},
			"dest_exchange_key": schema.StringAttribute{
				Description: "Routing key for destination exchange publishing (only used with dest_exchange).",
				Optional:    true,
			},
			"src_prefetch_count": schema.Int64Attribute{
				Description: "Number of messages to prefetch from source.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1000),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"src_delete_after": schema.StringAttribute{
				Description: "When to delete messages from source: 'never' or 'queue-length'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("never"),
				Validators: []validator.String{
					stringvalidator.OneOf("never", "queue-length"),
				},
			},
			"reconnect_delay": schema.Int64Attribute{
				Description: "Delay in seconds before reconnecting after connection failure.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(5),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"ack_mode": schema.StringAttribute{
				Description: "When to acknowledge messages from source: 'on-confirm', 'on-publish', or 'no-ack'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("on-confirm"),
				Validators: []validator.String{
					stringvalidator.OneOf("on-confirm", "on-publish", "no-ack"),
				},
			},
		},
	}
}

func (r *shovelResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.services = req.ProviderData.(*clientlibrary.Services)
}

func (r *shovelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan shovelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateShovelSource(plan); err != nil {
		resp.Diagnostics.AddError("Invalid shovel source configuration", err.Error())
		return
	}

	if err := validateShovelDestination(plan); err != nil {
		resp.Diagnostics.AddError("Invalid shovel destination configuration", err.Error())
		return
	}

	shovelValue := clientlibrary.ShovelValue{
		SrcURI:           plan.SrcURI.ValueString(),
		DestURI:          plan.DestURI.ValueString(),
		SrcQueue:         plan.SrcQueue.ValueString(),
		SrcExchange:      plan.SrcExchange.ValueString(),
		SrcExchangeKey:   plan.SrcExchangeKey.ValueString(),
		DestQueue:        plan.DestQueue.ValueString(),
		DestExchange:     plan.DestExchange.ValueString(),
		DestExchangeKey:  plan.DestExchangeKey.ValueString(),
		SrcPrefetchCount: plan.SrcPrefetchCount.ValueInt64(),
		SrcDeleteAfter:   plan.SrcDeleteAfter.ValueString(),
		ReconnectDelay:   plan.ReconnectDelay.ValueInt64(),
		AckMode:          plan.AckMode.ValueString(),
	}

	createReq := clientlibrary.ParameterRequest{
		Value: shovelValue,
	}

	err := r.services.Parameters.CreateOrUpdate(ctx, "shovel", plan.Vhost.ValueString(), plan.Name.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating shovel", err.Error())
		return
	}

	parameter, err := r.services.Parameters.Get(ctx, "shovel", plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read shovel data", err.Error())
		return
	}

	if err := updateStateFromParameter(ctx, &plan, parameter); err != nil {
		resp.Diagnostics.AddError("Failed to update state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "create diag failed")
	}
}

func (r *shovelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state shovelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parameter, err := r.services.Parameters.Get(ctx, "shovel", state.Vhost.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read shovel data", err.Error())
		return
	}
	if parameter == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := updateStateFromParameter(ctx, &state, parameter); err != nil {
		resp.Diagnostics.AddError("Failed to update state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *shovelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan shovelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateShovelSource(plan); err != nil {
		resp.Diagnostics.AddError("Invalid shovel source configuration", err.Error())
		return
	}

	if err := validateShovelDestination(plan); err != nil {
		resp.Diagnostics.AddError("Invalid shovel destination configuration", err.Error())
		return
	}

	shovelValue := clientlibrary.ShovelValue{
		SrcURI:           plan.SrcURI.ValueString(),
		DestURI:          plan.DestURI.ValueString(),
		SrcQueue:         plan.SrcQueue.ValueString(),
		SrcExchange:      plan.SrcExchange.ValueString(),
		SrcExchangeKey:   plan.SrcExchangeKey.ValueString(),
		DestQueue:        plan.DestQueue.ValueString(),
		DestExchange:     plan.DestExchange.ValueString(),
		DestExchangeKey:  plan.DestExchangeKey.ValueString(),
		SrcPrefetchCount: plan.SrcPrefetchCount.ValueInt64(),
		SrcDeleteAfter:   plan.SrcDeleteAfter.ValueString(),
		ReconnectDelay:   plan.ReconnectDelay.ValueInt64(),
		AckMode:          plan.AckMode.ValueString(),
	}

	updateReq := clientlibrary.ParameterRequest{
		Value: shovelValue,
	}

	err := r.services.Parameters.CreateOrUpdate(ctx, "shovel", plan.Vhost.ValueString(), plan.Name.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating shovel", err.Error())
		return
	}

	parameter, err := r.services.Parameters.Get(ctx, "shovel", plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read shovel data", err.Error())
		return
	}

	if err := updateStateFromParameter(ctx, &plan, parameter); err != nil {
		resp.Diagnostics.AddError("Failed to update state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *shovelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan shovelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.Parameters.Delete(ctx, "shovel", plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting shovel", err.Error())
		return
	}
}

func (r *shovelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "@")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: vhost@shovel_name",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vhost"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

func validateShovelSource(model shovelResourceModel) error {
	hasSrcQueue := !model.SrcQueue.IsNull() && !model.SrcQueue.IsUnknown() && model.SrcQueue.ValueString() != ""
	hasSrcExchange := !model.SrcExchange.IsNull() && !model.SrcExchange.IsUnknown() && model.SrcExchange.ValueString() != ""

	if hasSrcQueue && hasSrcExchange {
		return &ValidationError{"Cannot specify both src_queue and src_exchange"}
	}

	if !hasSrcQueue && !hasSrcExchange {
		return &ValidationError{"Must specify either src_queue or src_exchange"}
	}

	return nil
}

func validateShovelDestination(model shovelResourceModel) error {
	hasDestQueue := !model.DestQueue.IsNull() && !model.DestQueue.IsUnknown() && model.DestQueue.ValueString() != ""
	hasDestExchange := !model.DestExchange.IsNull() && !model.DestExchange.IsUnknown() && model.DestExchange.ValueString() != ""

	if hasDestQueue && hasDestExchange {
		return &ValidationError{"Cannot specify both dest_queue and dest_exchange"}
	}

	if !hasDestQueue && !hasDestExchange {
		return &ValidationError{"Must specify either dest_queue or dest_exchange"}
	}

	return nil
}

func updateStateFromParameter(ctx context.Context, state *shovelResourceModel, parameter *clientlibrary.ParameterResponse) error {
	state.Name = types.StringValue(parameter.Name)
	state.Vhost = types.StringValue(parameter.Vhost)

	valueBytes, err := json.Marshal(parameter.Value)
	if err != nil {
		return err
	}

	var shovelValue clientlibrary.ShovelValue
	if err := json.Unmarshal(valueBytes, &shovelValue); err != nil {
		return err
	}

	state.SrcURI = types.StringValue(shovelValue.SrcURI)
	state.DestURI = types.StringValue(shovelValue.DestURI)

	if shovelValue.SrcQueue != "" {
		state.SrcQueue = types.StringValue(shovelValue.SrcQueue)
	} else {
		state.SrcQueue = types.StringNull()
	}

	if shovelValue.SrcExchange != "" {
		state.SrcExchange = types.StringValue(shovelValue.SrcExchange)
	} else {
		state.SrcExchange = types.StringNull()
	}

	if shovelValue.SrcExchangeKey != "" {
		state.SrcExchangeKey = types.StringValue(shovelValue.SrcExchangeKey)
	} else {
		state.SrcExchangeKey = types.StringNull()
	}

	if shovelValue.DestQueue != "" {
		state.DestQueue = types.StringValue(shovelValue.DestQueue)
	} else {
		state.DestQueue = types.StringNull()
	}

	if shovelValue.DestExchange != "" {
		state.DestExchange = types.StringValue(shovelValue.DestExchange)
	} else {
		state.DestExchange = types.StringNull()
	}

	if shovelValue.DestExchangeKey != "" {
		state.DestExchangeKey = types.StringValue(shovelValue.DestExchangeKey)
	} else {
		state.DestExchangeKey = types.StringNull()
	}

	state.SrcPrefetchCount = types.Int64Value(shovelValue.SrcPrefetchCount)

	if shovelValue.SrcDeleteAfter != "" {
		state.SrcDeleteAfter = types.StringValue(shovelValue.SrcDeleteAfter)
	} else {
		state.SrcDeleteAfter = types.StringValue("never")
	}

	state.ReconnectDelay = types.Int64Value(shovelValue.ReconnectDelay)

	if shovelValue.AckMode != "" {
		state.AckMode = types.StringValue(shovelValue.AckMode)
	} else {
		state.AckMode = types.StringValue("on-confirm")
	}

	return nil
}

type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

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
	_ resource.Resource                = &federationUpstreamResource{}
	_ resource.ResourceWithConfigure   = &federationUpstreamResource{}
	_ resource.ResourceWithImportState = &federationUpstreamResource{}
)

func NewFederationUpstreamResource() resource.Resource {
	return &federationUpstreamResource{}
}

type federationUpstreamResource struct {
	services *clientlibrary.Services
}

type federationUpstreamResourceModel struct {
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

func (r *federationUpstreamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_federation_upstream"
}

func (r *federationUpstreamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a federation upstream for replicating exchanges and queues from remote brokers.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the federation upstream.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vhost": schema.StringAttribute{
				Description: "Virtual host where the federation upstream is defined.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uri": schema.StringAttribute{
				Description: "AMQP URI of the upstream broker.",
				Required:    true,
			},
			"prefetch_count": schema.Int64Attribute{
				Description: "Number of messages to prefetch from upstream.",
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
				Description: "When to acknowledge messages from upstream: 'on-confirm', 'on-publish', or 'no-ack'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("on-confirm"),
				Validators: []validator.String{
					stringvalidator.OneOf("on-confirm", "on-publish", "no-ack"),
				},
			},
			"exchange": schema.StringAttribute{
				Description: "Name of upstream exchange to federate from.",
				Optional:    true,
			},
			"max_hops": schema.Int64Attribute{
				Description: "Maximum number of federation hops for a message.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"expires": schema.Int64Attribute{
				Description: "Expiry time in milliseconds for the federation link.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"message_ttl": schema.Int64Attribute{
				Description: "TTL in milliseconds for messages in the federation link.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"queue": schema.StringAttribute{
				Description: "Name of upstream queue to federate from (for queue federation).",
				Optional:    true,
			},
			"consumer_tag": schema.StringAttribute{
				Description: "Consumer tag for the federation link.",
				Optional:    true,
			},
		},
	}
}

func (r *federationUpstreamResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.services = req.ProviderData.(*clientlibrary.Services)
}

func (r *federationUpstreamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan federationUpstreamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	federationValue := clientlibrary.FederationUpstreamValue{
		URI:            plan.URI.ValueString(),
		PrefetchCount:  plan.PrefetchCount.ValueInt64(),
		ReconnectDelay: plan.ReconnectDelay.ValueInt64(),
		AckMode:        plan.AckMode.ValueString(),
		Exchange:       plan.Exchange.ValueString(),
		MaxHops:        plan.MaxHops.ValueInt64(),
		Expires:        plan.Expires.ValueInt64(),
		MessageTTL:     plan.MessageTTL.ValueInt64(),
		Queue:          plan.Queue.ValueString(),
		ConsumerTag:    plan.ConsumerTag.ValueString(),
	}

	createReq := clientlibrary.ParameterRequest{
		Value: federationValue,
	}

	err := r.services.Parameters.CreateOrUpdate(ctx, "federation-upstream", plan.Vhost.ValueString(), plan.Name.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating federation upstream", err.Error())
		return
	}

	parameter, err := r.services.Parameters.Get(ctx, "federation-upstream", plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read federation upstream data", err.Error())
		return
	}

	if err := updateFederationUpstreamStateFromParameter(ctx, &plan, parameter); err != nil {
		resp.Diagnostics.AddError("Failed to update state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "create diag failed")
	}
}

func (r *federationUpstreamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state federationUpstreamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parameter, err := r.services.Parameters.Get(ctx, "federation-upstream", state.Vhost.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read federation upstream data", err.Error())
		return
	}
	if parameter == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := updateFederationUpstreamStateFromParameter(ctx, &state, parameter); err != nil {
		resp.Diagnostics.AddError("Failed to update state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *federationUpstreamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan federationUpstreamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	federationValue := clientlibrary.FederationUpstreamValue{
		URI:            plan.URI.ValueString(),
		PrefetchCount:  plan.PrefetchCount.ValueInt64(),
		ReconnectDelay: plan.ReconnectDelay.ValueInt64(),
		AckMode:        plan.AckMode.ValueString(),
		Exchange:       plan.Exchange.ValueString(),
		MaxHops:        plan.MaxHops.ValueInt64(),
		Expires:        plan.Expires.ValueInt64(),
		MessageTTL:     plan.MessageTTL.ValueInt64(),
		Queue:          plan.Queue.ValueString(),
		ConsumerTag:    plan.ConsumerTag.ValueString(),
	}

	updateReq := clientlibrary.ParameterRequest{
		Value: federationValue,
	}

	err := r.services.Parameters.CreateOrUpdate(ctx, "federation-upstream", plan.Vhost.ValueString(), plan.Name.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating federation upstream", err.Error())
		return
	}

	parameter, err := r.services.Parameters.Get(ctx, "federation-upstream", plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read federation upstream data", err.Error())
		return
	}

	if err := updateFederationUpstreamStateFromParameter(ctx, &plan, parameter); err != nil {
		resp.Diagnostics.AddError("Failed to update state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *federationUpstreamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan federationUpstreamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.Parameters.Delete(ctx, "federation-upstream", plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting federation upstream", err.Error())
		return
	}
}

func (r *federationUpstreamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "@")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: vhost@federation_upstream_name",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vhost"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

func updateFederationUpstreamStateFromParameter(ctx context.Context, state *federationUpstreamResourceModel, parameter *clientlibrary.ParameterResponse) error {
	state.Name = types.StringValue(parameter.Name)
	state.Vhost = types.StringValue(parameter.Vhost)

	valueBytes, err := json.Marshal(parameter.Value)
	if err != nil {
		return err
	}

	var federationValue clientlibrary.FederationUpstreamValue
	if err := json.Unmarshal(valueBytes, &federationValue); err != nil {
		return err
	}

	state.URI = types.StringValue(federationValue.URI)
	state.PrefetchCount = types.Int64Value(federationValue.PrefetchCount)
	state.ReconnectDelay = types.Int64Value(federationValue.ReconnectDelay)

	if federationValue.AckMode != "" {
		state.AckMode = types.StringValue(federationValue.AckMode)
	} else {
		state.AckMode = types.StringValue("on-confirm")
	}

	if federationValue.Exchange != "" {
		state.Exchange = types.StringValue(federationValue.Exchange)
	} else {
		state.Exchange = types.StringNull()
	}

	state.MaxHops = types.Int64Value(federationValue.MaxHops)

	if federationValue.Expires > 0 {
		state.Expires = types.Int64Value(federationValue.Expires)
	} else {
		state.Expires = types.Int64Null()
	}

	if federationValue.MessageTTL > 0 {
		state.MessageTTL = types.Int64Value(federationValue.MessageTTL)
	} else {
		state.MessageTTL = types.Int64Null()
	}

	if federationValue.Queue != "" {
		state.Queue = types.StringValue(federationValue.Queue)
	} else {
		state.Queue = types.StringNull()
	}

	if federationValue.ConsumerTag != "" {
		state.ConsumerTag = types.StringValue(federationValue.ConsumerTag)
	} else {
		state.ConsumerTag = types.StringNull()
	}

	return nil
}

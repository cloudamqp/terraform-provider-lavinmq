package lavinmq

import (
	"context"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

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
	Name    types.String `tfsdk:"name"`
	Vhost   types.String `tfsdk:"vhost"`
	SrcUri  types.String `tfsdk:"src_uri"`
	DestUri types.String `tfsdk:"dest_uri"`
	// Optional parameters
	SrcQueue         types.String `tfsdk:"src_queue"`
	SrcExchange      types.String `tfsdk:"src_exchange"`
	SrcExchangeKey   types.String `tfsdk:"src_exchange_key"`
	SrcPrefetchCount types.Int64  `tfsdk:"src_prefetch_count"`
	SrcDelayAfter    types.String `tfsdk:"src_delay_after"`
	DestQueue        types.String `tfsdk:"dest_queue"`
	DestExchange     types.String `tfsdk:"dest_exchange"`
	DestExchangeKey  types.String `tfsdk:"dest_exchange_key"`
	ReconnectDelay   types.Int64  `tfsdk:"reconnect_delay"`
	AckMode          types.String `tfsdk:"ack_mode"`
}

func (r *shovelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shovel"
}

func (r *shovelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a shovel.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the managed shovel.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vhost": schema.StringAttribute{
				Description: "The vhost the shovel is located in.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"src_uri": schema.StringAttribute{
				Description: "The source URI for the shovel.",
				Required:    true,
			},
			"dest_uri": schema.StringAttribute{
				Description: "The destination URI for the shovel.",
				Required:    true,
			},
			"src_queue": schema.StringAttribute{
				Description: "The source queue name.",
				Optional:    true,
			},
			"src_exchange": schema.StringAttribute{
				Description: "The source exchange name.",
				Optional:    true,
			},
			"src_exchange_key": schema.StringAttribute{
				Description: "The source exchange key.",
				Optional:    true,
			},
			"src_prefetch_count": schema.Int64Attribute{
				Description: "The source prefetch count.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"src_delay_after": schema.StringAttribute{
				Description: "The source delay after.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("never", "queue-length"),
				},
			},
			"dest_queue": schema.StringAttribute{
				Description: "The destination queue name.",
				Optional:    true,
			},
			"dest_exchange": schema.StringAttribute{
				Description: "The destination exchange name.",
				Optional:    true,
			},
			"dest_exchange_key": schema.StringAttribute{
				Description: "The destination exchange key.",
				Optional:    true,
			},
			"reconnect_delay": schema.Int64Attribute{
				Description: "The reconnect delay in seconds.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"ack_mode": schema.StringAttribute{
				Description: "The acknowledgment mode.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("on-confirm", "no-ack", "on-publish"),
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
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := r.populateShovelRequest(&plan)
	err := r.services.ShovelParameters.CreateOrUpdate(ctx, plan.Vhost.ValueString(), plan.Name.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating shovel",
			"Could not create shovel, unexpected error: "+err.Error(),
		)
		return
	}

	// Read the newly created resource
	var state shovelResourceModel
	shovelResp, err := r.services.ShovelParameters.Get(ctx, plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading shovel",
			"Could not read shovel after creation, unexpected error: "+err.Error(),
		)
		return
	}
	if shovelResp == nil {
		resp.Diagnostics.AddError(
			"Error reading shovel",
			"Could not read shovel after creation, shovel not found",
		)
		return
	}

	r.populateShovelResourceModel(&state, shovelResp)
	resp.State.Set(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *shovelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state shovelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	shovelResp, err := r.services.ShovelParameters.Get(ctx, state.Vhost.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading shovel",
			"Could not read shovel, unexpected error: "+err.Error(),
		)
		return
	}
	if shovelResp == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	r.populateShovelResourceModel(&state, shovelResp)
	resp.State.Set(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *shovelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan shovelResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := r.populateShovelRequest(&plan)
	err := r.services.ShovelParameters.CreateOrUpdate(ctx, plan.Vhost.ValueString(), plan.Name.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating shovel",
			"Could not update shovel, unexpected error: "+err.Error(),
		)
		return
	}

	// Read the updated resource
	var state shovelResourceModel
	shovelResp, err := r.services.ShovelParameters.Get(ctx, plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading shovel",
			"Could not read shovel after update, unexpected error: "+err.Error(),
		)
		return
	}
	if shovelResp == nil {
		resp.Diagnostics.AddError(
			"Error reading shovel",
			"Could not read shovel after update, shovel not found",
		)
		return
	}

	r.populateShovelResourceModel(&state, shovelResp)
	resp.State.Set(ctx, &state)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *shovelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state shovelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.ShovelParameters.Delete(ctx, state.Vhost.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting shovel",
			"Could not delete shovel, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *shovelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import ID is expected to be in the format <vhost>@<name>
	importIDParts := strings.SplitN(req.ID, "@", 2)
	if len(importIDParts) != 2 {
		resp.Diagnostics.AddError(
			"Error importing shovel",
			"Invalid import ID format. Expected format: <vhost>@<name>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vhost"), importIDParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), importIDParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *shovelResource) populateShovelResourceModel(state *shovelResourceModel, resp *clientlibrary.ShovelParametersResponse) {
	state.Name = types.StringValue(resp.Name)
	state.Vhost = types.StringValue(resp.Vhost)
	state.SrcUri = types.StringValue(resp.Value.SrcUri)
	state.DestUri = types.StringValue(resp.Value.DestUri)

	// Optional parameters
	if resp.Value.SrcQueue != nil {
		state.SrcQueue = types.StringValue(*resp.Value.SrcQueue)
	}
	if resp.Value.SrcExchange != nil {
		state.SrcExchange = types.StringValue(*resp.Value.SrcExchange)
	}
	if resp.Value.SrcExchangeKey != nil {
		state.SrcExchangeKey = types.StringValue(*resp.Value.SrcExchangeKey)
	}
	if resp.Value.SrcPrefetchCount != nil {
		state.SrcPrefetchCount = types.Int64Value(*resp.Value.SrcPrefetchCount)
	}
	if resp.Value.SrcDelayAfter != nil {
		state.SrcDelayAfter = types.StringValue(*resp.Value.SrcDelayAfter)
	}
	if resp.Value.DestQueue != nil {
		state.DestQueue = types.StringValue(*resp.Value.DestQueue)
	}
	if resp.Value.DestExchange != nil {
		state.DestExchange = types.StringValue(*resp.Value.DestExchange)
	}
	if resp.Value.DestExchangeKey != nil {
		state.DestExchangeKey = types.StringValue(*resp.Value.DestExchangeKey)
	}
	if resp.Value.ReconnectDelay != nil {
		state.ReconnectDelay = types.Int64Value(*resp.Value.ReconnectDelay)
	}
	if resp.Value.AckMode != nil {
		state.AckMode = types.StringValue(*resp.Value.AckMode)
	}
}

func (r *shovelResource) populateShovelRequest(plan *shovelResourceModel) clientlibrary.ShovelParametersRequest {
	shovel := clientlibrary.ShovelParametersObject{
		SrcUri:  plan.SrcUri.ValueString(),
		DestUri: plan.DestUri.ValueString(),
	}

	// Optional parameters
	if !plan.SrcQueue.IsNull() {
		shovel.SrcQueue = plan.SrcQueue.ValueStringPointer()
	}
	if !plan.SrcExchange.IsNull() {
		shovel.SrcExchange = plan.SrcExchange.ValueStringPointer()
	}
	if !plan.SrcExchangeKey.IsNull() {
		shovel.SrcExchangeKey = plan.SrcExchangeKey.ValueStringPointer()
	}
	if !plan.SrcPrefetchCount.IsNull() {
		shovel.SrcPrefetchCount = plan.SrcPrefetchCount.ValueInt64Pointer()
	}
	if !plan.SrcDelayAfter.IsNull() {
		shovel.SrcDelayAfter = plan.SrcDelayAfter.ValueStringPointer()
	}
	if !plan.DestQueue.IsNull() {
		shovel.DestQueue = plan.DestQueue.ValueStringPointer()
	}
	if !plan.DestExchange.IsNull() {
		shovel.DestExchange = plan.DestExchange.ValueStringPointer()
	}
	if !plan.DestExchangeKey.IsNull() {
		shovel.DestExchangeKey = plan.DestExchangeKey.ValueStringPointer()
	}
	if !plan.ReconnectDelay.IsNull() {
		shovel.ReconnectDelay = plan.ReconnectDelay.ValueInt64Pointer()
	}
	if !plan.AckMode.IsNull() {
		shovel.AckMode = plan.AckMode.ValueStringPointer()
	}

	request := clientlibrary.ShovelParametersRequest{
		Value: shovel,
	}

	return request
}

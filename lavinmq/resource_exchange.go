package lavinmq

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &exchangeResource{}
	_ resource.ResourceWithConfigure   = &exchangeResource{}
	_ resource.ResourceWithImportState = &exchangeResource{}
)

// NewExchangeResource is a helper function to simplify the provider implementation.
func NewExchangeResource() resource.Resource {
	return &exchangeResource{}
}

// exchangeResource is the resource implementation.
type exchangeResource struct {
	services *clientlibrary.Services
}

// exchangeResourceModel is the
type exchangeResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Vhost      types.String `tfsdk:"vhost"`
	Type       types.String `tfsdk:"type"`
	AutoDelete types.Bool   `tfsdk:"auto_delete"`
	Durable    types.Bool   `tfsdk:"durable"`
	Internal   types.Bool   `tfsdk:"internal"`
}

// Metadata returns the data source type name.
func (r *exchangeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exchange"
}

// Schema defines the schema for the resource.
func (r *exchangeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage an exchange.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the managed exchange.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vhost": schema.StringAttribute{
				Description: "The vhost the exchange is located in.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The exchange type (direct, fanout, topic, headers).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auto_delete": schema.BoolAttribute{
				Description: "Whether the exchange is automatically deleted when no longer used.",
				Optional:    true,
				Computed:    true,
				// Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"durable": schema.BoolAttribute{
				Description: "Whether the exchange should survive a broker restart.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"internal": schema.BoolAttribute{
				Description: "Whether the exchange is internal (cannot be published to directly).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured services to the resource.
func (r *exchangeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.services = req.ProviderData.(*clientlibrary.Services)
}

func (r *exchangeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIDParts := strings.Split(req.ID, ",")

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vhost"), importIDParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), importIDParts[1])...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *exchangeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan exchangeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var request clientlibrary.ExchangeRequest
	request.Type = plan.Type.ValueString()
	if !plan.AutoDelete.IsUnknown() {
		request.AutoDelete = plan.AutoDelete.ValueBoolPointer()
	}
	if !plan.Durable.IsUnknown() {
		request.Durable = plan.Durable.ValueBoolPointer()
	}
	if !plan.Internal.IsUnknown() {
		request.Internal = plan.Internal.ValueBoolPointer()
	}

	err := r.services.Exchanges.CreateOrUpdate(ctx, plan.Vhost.ValueString(), plan.Name.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating exchange", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.Vhost.ValueString(), plan.Name.ValueString()))
	tflog.Info(ctx, "Created exchange", map[string]any{"id": plan.ID.ValueString()})

	// Read back the exchange to get the actual state from the server
	exchange, err := r.services.Exchanges.Get(ctx, plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading exchange after creation", err.Error())
		return
	}

	// Update the plan with actual values from the server
	plan.Type = types.StringValue(exchange.Type)
	plan.AutoDelete = types.BoolValue(exchange.AutoDelete)
	plan.Durable = types.BoolValue(exchange.Durable)
	plan.Internal = types.BoolValue(exchange.Internal)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// TODO: Check so import handles default values
func (r *exchangeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state exchangeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	exchange, err := r.services.Exchanges.Get(ctx, state.Vhost.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading exchange", err.Error())
		return
	}

	state.Type = types.StringValue(exchange.Type)
	state.AutoDelete = types.BoolValue(exchange.AutoDelete)
	state.Durable = types.BoolValue(exchange.Durable)
	state.Internal = types.BoolValue(exchange.Internal)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *exchangeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This resource does not implement the Update function
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *exchangeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state exchangeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.Exchanges.Delete(ctx, state.Vhost.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting exchange", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

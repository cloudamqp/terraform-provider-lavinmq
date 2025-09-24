package lavinmq

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vhostResource{}
	_ resource.ResourceWithConfigure   = &vhostResource{}
	_ resource.ResourceWithImportState = &vhostResource{}
)

// NewVhostResource is a helper function to simplify the provider implementation.
func NewVhostResource() resource.Resource {
	return &vhostResource{}
}

// vhostResource is the resource implementation.
type vhostResource struct {
	services *clientlibrary.Services
}

// vhostResourceModel is the
type vhostResourceModel struct {
	Name           types.String `tfsdk:"name"`
	MaxConnections types.Int64  `tfsdk:"max_connections"`
	MaxQueues      types.Int64  `tfsdk:"max_queues"`
}

// Metadata returns the data source type name.
func (r *vhostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vhost"
}

// Schema defines the schema for the resource.
func (r *vhostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a vhost.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the managed vhost.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"max_connections": schema.Int64Attribute{
				Description: "Limit the number of connections for the vhost.",
				Optional:    true,
				Default:     nil,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_queues": schema.Int64Attribute{
				Description: "Limit the number of queues for the vhost.",
				Optional:    true,
				Default:     nil,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured services to the resource.
func (r *vhostResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.services = req.ProviderData.(*clientlibrary.Services)
}

// Create creates the resource and sets the initial Terraform state.
func (r *vhostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vhostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.Vhosts.CreateOrUpdate(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	var limits clientlibrary.VhostLimits
	updateLimits := false
	if !plan.MaxConnections.IsNull() {
		limits.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
		updateLimits = true
	}
	if !plan.MaxQueues.IsNull() {
		limits.MaxQueues = plan.MaxQueues.ValueInt64Pointer()
		updateLimits = true
	}
	if updateLimits {
		err := r.services.VhostLimits.Update(ctx, plan.Name.ValueString(), limits)
		if err != nil {
			resp.Diagnostics.AddError("Error setting limits", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "created diag failed")
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vhostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vhostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Name.IsUnknown() {
		tflog.Info(ctx, fmt.Sprintf("import resource with name identifier %s", state.Name))
	}

	vhost, err := r.services.Vhosts.Get(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read vhost data", err.Error())
		return
	}

	state.Name = types.StringValue(vhost.Name)

	limits, err := r.services.VhostLimits.Get(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read limits data", err.Error())
		return
	}

	if limits.Value.MaxConnections == nil {
		state.MaxConnections = types.Int64Null()

	} else {
		state.MaxConnections = types.Int64PointerValue(limits.Value.MaxConnections)
	}

	if limits.Value.MaxQueues == nil {
		state.MaxQueues = types.Int64Null()
	} else {
		state.MaxQueues = types.Int64PointerValue(limits.Value.MaxQueues)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vhostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vhostResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var limits clientlibrary.VhostLimits
	if plan.MaxConnections.IsNull() {
		limits.MaxConnections = nil
	} else {
		limits.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if plan.MaxQueues.IsNull() {
		limits.MaxQueues = nil
	} else {
		limits.MaxQueues = plan.MaxQueues.ValueInt64Pointer()
	}
	err := r.services.VhostLimits.Update(ctx, plan.Name.ValueString(), limits)
	if err != nil {
		resp.Diagnostics.AddError("Error setting limits", err.Error())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "update diag failed")
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vhostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan vhostResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.Vhosts.Delete(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
		return
	}
}

func (r *vhostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import resource by name argument
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

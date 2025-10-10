package lavinmq

import (
	"context"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &permissionResource{}
	_ resource.ResourceWithConfigure   = &permissionResource{}
	_ resource.ResourceWithImportState = &permissionResource{}
)

// NewPermissionResource is a helper function to simplify the provider implementation.
func NewPermissionResource() resource.Resource {
	return &permissionResource{}
}

// permissionResource is the resource implementation.
type permissionResource struct {
	services *clientlibrary.Services
}

type permissionResourceModel struct {
	Vhost     types.String `tfsdk:"vhost"`
	User      types.String `tfsdk:"user"`
	Configure types.String `tfsdk:"configure"`
	Read      types.String `tfsdk:"read"`
	Write     types.String `tfsdk:"write"`
}

// Metadata returns the resource type name.
func (r *permissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

// Schema defines the schema for the resource.
func (r *permissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage user permissions for a vhost.",
		Attributes: map[string]schema.Attribute{
			"vhost": schema.StringAttribute{
				Description: "Virtual host where the permission is applied.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user": schema.StringAttribute{
				Description: "Name of the user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"configure": schema.StringAttribute{
				Description: "Regular expression pattern for configure permissions.",
				Required:    true,
			},
			"read": schema.StringAttribute{
				Description: "Regular expression pattern for read permissions.",
				Required:    true,
			},
			"write": schema.StringAttribute{
				Description: "Regular expression pattern for write permissions.",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured services to the resource.
func (r *permissionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.services = req.ProviderData.(*clientlibrary.Services)
}

// Create creates the resource and sets the initial Terraform state.
func (r *permissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan permissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := clientlibrary.PermissionRequest{
		Configure: plan.Configure.ValueString(),
		Read:      plan.Read.ValueString(),
		Write:     plan.Write.ValueString(),
	}

	err := r.services.Permissions.CreateOrUpdate(ctx, plan.Vhost.ValueString(), plan.User.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating permission", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "create diag failed")
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *permissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state permissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := r.services.Permissions.Get(ctx, state.Vhost.ValueString(), state.User.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read permission data", err.Error())
		return
	}
	if permission == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Vhost = types.StringValue(permission.Vhost)
	state.User = types.StringValue(permission.User)
	state.Configure = types.StringValue(permission.Configure)
	state.Read = types.StringValue(permission.Read)
	state.Write = types.StringValue(permission.Write)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *permissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan permissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := clientlibrary.PermissionRequest{
		Configure: plan.Configure.ValueString(),
		Read:      plan.Read.ValueString(),
		Write:     plan.Write.ValueString(),
	}

	err := r.services.Permissions.CreateOrUpdate(ctx, plan.Vhost.ValueString(), plan.User.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating permission", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *permissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan permissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.Permissions.Delete(ctx, plan.Vhost.ValueString(), plan.User.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting permission", err.Error())
		return
	}
}

func (r *permissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import resource by vhost@user (e.g., "/@my-user" or "my-vhost@my-user")
	parts := strings.Split(req.ID, "@")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: vhost@user",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vhost"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user"), parts[1])...)
}

package lavinmq

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq/converters"
	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	client *clientlibrary.Client
}

type userResourceModel struct {
	Name         types.String `tfsdk:"name"`
	Password     types.String `tfsdk:"password"`
	PasswordHash types.String `tfsdk:"password_hash"`
	Tags         types.List   `tfsdk:"tags"`
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a user.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the managed user.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Description: "Password of the managed user.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_hash": schema.StringAttribute{
				Description: "Hashed version of the password.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.ListAttribute{
				Description: "List of tags associated with the user.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.OneOf(
						"administrator",
						"monitoring",
						"management",
						"policymaker",
						"impersonator",
					)),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*clientlibrary.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var createReq clientlibrary.UserRequest
	if !plan.Password.IsNull() {
		createReq.Password = utils.Pointer(plan.Password.ValueString())
	}
	if !plan.PasswordHash.IsUnknown() {
		createReq.PasswordHash = utils.Pointer(plan.PasswordHash.ValueString())
	}
	if !plan.Tags.IsUnknown() {
		var tags []string
		diags := plan.Tags.ElementsAs(ctx, &tags, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		createReq.Tags = strings.Join(tags, ",")
	}

	err := r.client.Users.CreateOrUpdate(ctx, plan.Name.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	// Read out computed values
	user, err := r.client.Users.Get(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user state", err.Error())
		return
	}

	if plan.PasswordHash.IsUnknown() {
		plan.PasswordHash = types.StringValue(*user.PasswordHash)
	}

	if user.Tags != "" {
		tags := strings.Split(user.Tags, ",")
		plan.Tags, _ = types.ListValue(types.StringType, converters.StringsToAttrValues(tags))
	} else {
		plan.Tags, _ = types.ListValue(types.StringType, []attr.Value{})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "created diag failed")
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Name.IsUnknown() {
		tflog.Info(ctx, fmt.Sprintf("import resource with name identifier %s", state.Name))
	}

	user, err := r.client.Users.Get(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user data", err.Error())
		return
	}
	if user == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(user.Name)
	if user.PasswordHash != nil {
		state.PasswordHash = types.StringValue(*user.PasswordHash)
	}

	if user.Tags != "" {
		tags := strings.Split(user.Tags, ",")
		state.Tags, _ = types.ListValue(types.StringType, converters.StringsToAttrValues(tags))
	} else {
		state.Tags, _ = types.ListValue(types.StringType, []attr.Value{})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	var state userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateReq clientlibrary.UserRequest
	if !plan.Password.IsNull() && !state.Password.Equal(plan.Password) {
		updateReq.Password = utils.Pointer(plan.Password.ValueString())
	}
	if !plan.PasswordHash.IsUnknown() && !state.PasswordHash.Equal(plan.PasswordHash) {
		updateReq.PasswordHash = utils.Pointer(plan.PasswordHash.ValueString())
	}
	if !plan.Tags.IsNull() && !state.Tags.Equal(plan.Tags) {
		var tags []string
		diags := plan.Tags.ElementsAs(ctx, &tags, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		updateReq.Tags = strings.Join(tags, ",")
	}

	err := r.client.Users.CreateOrUpdate(ctx, plan.Name.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	// Read out computed values
	user, err := r.client.Users.Get(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user state", err.Error())
		return
	}

	if plan.PasswordHash.IsUnknown() || plan.PasswordHash.IsNull() {
		plan.PasswordHash = types.StringValue(*user.PasswordHash)
	}
	if plan.Tags.IsUnknown() || plan.Tags.IsNull() {
		if user.Tags != "" {
			tags := strings.Split(user.Tags, ",")
			plan.Tags, _ = types.ListValue(types.StringType, converters.StringsToAttrValues(tags))
		} else {
			plan.Tags, _ = types.ListValue(types.StringType, []attr.Value{})
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "update diag failed")
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Users.Delete(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import resource by name argument
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

package lavinmq

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"
	"github.com/cloudamqp/terraform-provider-lavinmq/lavinmq/converters"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	services *clientlibrary.Services
}

type userResourceModel struct {
	Name            types.String `tfsdk:"name"`
	Password        types.String `tfsdk:"password"`
	PasswordVersion types.Int64  `tfsdk:"password_version"`
	// PasswordHash *userPasswordHashModel `tfsdk:"password_hash"`
	Tags types.List `tfsdk:"tags"`
}

type userPasswordHashModel struct {
	Value     types.String `tfsdk:"value"`
	Algorithm types.String `tfsdk:"algorithm"`
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
				WriteOnly:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				// Validators: []validator.String{
				// 	stringvalidator.ConflictsWith(path.MatchRoot("password_hash")),
				// },
			},
			"password_version": schema.Int64Attribute{
				Description: "Version of write only password or password hash.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// "password_hash": schema.SingleNestedAttribute{
			// 	Description: "Hashed version of the password.",
			// 	Optional:    true,
			// 	WriteOnly:   true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"value": schema.StringAttribute{
			// 			Description: "The hashed password value.",
			// 			Required:    true,
			// 			WriteOnly:   true,
			// 			Validators: []validator.String{
			// 				stringvalidator.ConflictsWith(path.MatchRoot("password")),
			// 			},
			// 		},
			// 		"algorithm": schema.StringAttribute{
			// 			Description: "The hashing algorithm used.",
			// 			Optional:    true,
			// 			WriteOnly:   true,
			// 			Validators: []validator.String{
			// 				stringvalidator.OneOfCaseInsensitive("sha256", "sha512", "bcrypt", "MD5"),
			// 			},
			// 		},
			// 	},
			// },
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

// Configure adds the provider configured services to the resource.
func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.services = req.ProviderData.(*clientlibrary.Services)
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config userResourceModel
	var plan userResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var request clientlibrary.UserRequest

	// Password or PasswordHash must be set
	// if plan.Password.IsNull() && plan.PasswordHash() {
	// 	resp.Diagnostics.AddError(
	// 		"Missing required attribute",
	// 		"Either 'password' or 'password_hash' must be specified to create a user.",
	// 	)
	// 	return
	// }

	// tflog.Info(ctx, fmt.Sprintf("config password_hash: %+v", config.PasswordHash))
	// tflog.Info(ctx, fmt.Sprintf("config password_hash.value: %s", config.PasswordHash.Value.ValueString()))
	// tflog.Info(ctx, fmt.Sprintf("config password_hash.algorithm: %s", config.PasswordHash.Algorithm.ValueString()))

	// if config.PasswordHash != nil {
	// 	request.PasswordHash = config.PasswordHash.Value.ValueString()
	// 	resp.Private.SetKey(ctx, "password_hash", []byte(config.PasswordHash.Value.ValueString()))
	// 	if config.PasswordHash.Algorithm.IsNull() {
	// 		request.HashingAlgorithm = "sha256"
	// 	} else {
	// 		request.HashingAlgorithm = config.PasswordHash.Algorithm.ValueString()
	// 	}
	// }

	if !config.Password.IsNull() {
		request.Password = config.Password.ValueString()
		// resp.Private.SetKey(ctx, "password", []byte(config.Password.ValueString()))
	}

	if !plan.Tags.IsNull() {
		var tags []string
		diags := plan.Tags.ElementsAs(ctx, &tags, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		request.Tags = strings.Join(tags, ",")
	}

	tflog.Info(ctx, fmt.Sprintf("creating user with name: %s and values: %+v", plan.Name.ValueString(), request))

	err := r.services.Users.CreateOrUpdate(ctx, plan.Name.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	// if plan.PasswordHash.IsNull() {
	// 	// Read out computed values
	// 	user, err := r.services.Users.Get(ctx, plan.Name.ValueString())
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Failed to read user state", err.Error())
	// 		return
	// 	}

	// 	plan.PasswordHash = types.StringValue(user.PasswordHash)
	// 	plan.HashingAlgorithm = types.StringValue(r.filterHashingAlgorithm(user.HashingAlgorithm))

	// 	if plan.Tags.IsNull() {
	// 		if user.Tags != "" {
	// 			tags := strings.Split(user.Tags, ",")
	// 			plan.Tags, _ = types.ListValue(types.StringType, converters.StringsToAttrValues(tags))
	// 		} else {
	// 			plan.Tags, _ = types.ListValue(types.StringType, []attr.Value{})
	// 		}
	// 	}
	// }

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

	user, err := r.services.Users.Get(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user data", err.Error())
		return
	}
	if user == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(user.Name)

	// oldPasswordHash, _ := req.Private.GetKey(ctx, "password_hash")
	// if string(oldPasswordHash) != user.PasswordHash {
	// 	tflog.Info(ctx, "password hash changed outside of terraform, updating state")
	// 	resp.Private.SetKey(ctx, "password_hash", []byte(user.PasswordHash))
	// }

	// state.PasswordHash = &userPasswordHashModel{
	// 	Value:     types.StringValue(string(oldPasswordHash)),
	// 	Algorithm: types.StringValue(r.filterHashingAlgorithm(user.HashingAlgorithm)),
	// }
	// if user.PasswordHash != nil {
	// 	state.PasswordHash = types.StringValue(*user.PasswordHash)
	// }

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
	var config userResourceModel
	var plan userResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var request clientlibrary.UserRequest

	if !config.Password.IsNull() {
		request.Password = config.Password.ValueString()
	}

	// if !plan.PasswordHash.IsUnknown() && !state.PasswordHash.Equal(plan.PasswordHash) {
	// 	updateReq.PasswordHash = utils.Pointer(plan.PasswordHash.ValueString())
	// }

	if !plan.Tags.IsNull() {
		var tags []string
		diags := plan.Tags.ElementsAs(ctx, &tags, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		request.Tags = strings.Join(tags, ",")
	}

	err := r.services.Users.CreateOrUpdate(ctx, plan.Name.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	// // Read out computed values
	// user, err := r.services.Users.Get(ctx, plan.Name.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError("Failed to read user state", err.Error())
	// 	return
	// }

	// // if plan.PasswordHash.IsUnknown() || plan.PasswordHash.IsNull() {
	// // 	plan.PasswordHash = types.StringValue(*user.PasswordHash)
	// // }
	// if plan.Tags.IsUnknown() || plan.Tags.IsNull() {
	// 	if user.Tags != "" {
	// 		tags := strings.Split(user.Tags, ",")
	// 		plan.Tags, _ = types.ListValue(types.StringType, converters.StringsToAttrValues(tags))
	// 	} else {
	// 		plan.Tags, _ = types.ListValue(types.StringType, []attr.Value{})
	// 	}
	// }

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

	err := r.services.Users.Delete(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import resource by name argument
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *userResource) filterHashingAlgorithm(algorithm string) string {
	parts := strings.Split(algorithm, "_")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return algorithm
}

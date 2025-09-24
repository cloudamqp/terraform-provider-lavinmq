package lavinmq

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-lavinmq/clientlibrary"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &policyResource{}
	_ resource.ResourceWithConfigure   = &policyResource{}
	_ resource.ResourceWithImportState = &policyResource{}
)

// NewPolicyResource is a helper function to simplify the provider implementation.
func NewPolicyResource() resource.Resource {
	return &policyResource{}
}

// policyResource is the resource implementation.
type policyResource struct {
	services *clientlibrary.Services
}

type policyResourceModel struct {
	Name       types.String `tfsdk:"name"`
	Vhost      types.String `tfsdk:"vhost"`
	Pattern    types.String `tfsdk:"pattern"`
	Definition types.Map    `tfsdk:"definition"`
	Priority   types.Int64  `tfsdk:"priority"`
	ApplyTo    types.String `tfsdk:"apply_to"`
}

// Metadata returns the resource type name.
func (r *policyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

// Schema defines the schema for the resource.
func (r *policyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a policy.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the policy.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vhost": schema.StringAttribute{
				Description: "Virtual host where the policy is applied.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pattern": schema.StringAttribute{
				Description: "Regular expression pattern that matches the names of exchanges or queues to which the policy applies.",
				Required:    true,
			},
			"definition": schema.MapAttribute{
				Description: "Policy definition as a map of key-value pairs.",
				Required:    true,
				ElementType: types.StringType,
			},
			"priority": schema.Int64Attribute{
				Description: "Policy priority. Higher numbers indicate higher priority.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"apply_to": schema.StringAttribute{
				Description: "What the policy applies to: 'all', 'exchanges', or 'queues'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("all"),
				Validators: []validator.String{
					stringvalidator.OneOf("all", "exchanges", "queues"),
				},
			},
		},
	}
}

// Configure adds the provider configured services to the resource.
func (r *policyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.services = req.ProviderData.(*clientlibrary.Services)
}

// Create creates the resource and sets the initial Terraform state.
func (r *policyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the definition map to a map[string]any
	definitionMap := make(map[string]any)
	if !plan.Definition.IsNull() && !plan.Definition.IsUnknown() {
		for key, value := range plan.Definition.Elements() {
			if strValue, ok := value.(types.String); ok {
				definitionMap[key] = convertPolicyValue(key, strValue.ValueString())
			}
		}
	}

	createReq := clientlibrary.PolicyRequest{
		Pattern:    plan.Pattern.ValueString(),
		Definition: definitionMap,
		Priority:   plan.Priority.ValueInt64(),
		ApplyTo:    plan.ApplyTo.ValueString(),
	}

	err := r.services.Policies.CreateOrUpdate(ctx, plan.Vhost.ValueString(), plan.Name.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating policy", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "create diag failed")
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *policyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.services.Policies.Get(ctx, state.Vhost.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read policy data", err.Error())
		return
	}
	if policy == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(policy.Name)
	state.Vhost = types.StringValue(policy.Vhost)
	state.Pattern = types.StringValue(policy.Pattern)
	state.Priority = types.Int64Value(int64(policy.Priority))
	state.ApplyTo = types.StringValue(policy.ApplyTo)

	// Convert definition to types.Map
	if len(policy.Definition) > 0 {
		elements := make(map[string]attr.Value)
		for key, value := range policy.Definition {
			elements[key] = types.StringValue(fmt.Sprintf("%v", value))
		}
		definitionMap, diags := types.MapValue(types.StringType, elements)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		state.Definition = definitionMap
	} else {
		emptyElements := make(map[string]attr.Value)
		emptyMap, _ := types.MapValue(types.StringType, emptyElements)
		state.Definition = emptyMap
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *policyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan policyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the definition map to a map[string]any
	definitionMap := make(map[string]any)
	if !plan.Definition.IsNull() && !plan.Definition.IsUnknown() {
		for key, value := range plan.Definition.Elements() {
			if strValue, ok := value.(types.String); ok {
				definitionMap[key] = convertPolicyValue(key, strValue.ValueString())
			}
		}
	}

	updateReq := clientlibrary.PolicyRequest{
		Pattern:    plan.Pattern.ValueString(),
		Definition: definitionMap,
		Priority:   plan.Priority.ValueInt64(),
		ApplyTo:    plan.ApplyTo.ValueString(),
	}

	err := r.services.Policies.CreateOrUpdate(ctx, plan.Vhost.ValueString(), plan.Name.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating policy", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *policyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan policyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.services.Policies.Delete(ctx, plan.Vhost.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting policy", err.Error())
		return
	}
}

func (r *policyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import resource by vhost@name (e.g., "my-vhost@my-policy" or "/@my-policy")
	parts := strings.Split(req.ID, "@")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: vhost@policy_name",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vhost"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

// convertPolicyValue converts string values to appropriate types for policy definitions
func convertPolicyValue(key, value string) any {
	// Numeric policy keys that should be converted to integers
	numericKeys := map[string]bool{
		"message-ttl":       true,
		"max-length":        true,
		"max-length-bytes":  true,
		"expires":           true,
		"priority":          true,
		"ha-params":         true,
		"overflow-length":   true,
		"delivery-limit":    true,
	}
	
	if numericKeys[key] {
		// Try to parse as integer
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	
	// For non-numeric keys or if parsing fails, return as string
	return value
}

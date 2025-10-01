package lavinmq

import (
	"context"
	"math/big"
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
	Name       types.String  `tfsdk:"name"`
	Vhost      types.String  `tfsdk:"vhost"`
	Pattern    types.String  `tfsdk:"pattern"`
	Definition types.Dynamic `tfsdk:"definition"`
	Priority   types.Int64   `tfsdk:"priority"`
	ApplyTo    types.String  `tfsdk:"apply_to"`
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
			"definition": schema.DynamicAttribute{
				Description: "Policy definition as a map of key-value pairs.",
				Required:    true,
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

	definitionMap := make(map[string]any)
	if !plan.Definition.IsNull() && !plan.Definition.IsUnknown() {
		underlyingValue := plan.Definition.UnderlyingValue()

		if mapValue, ok := underlyingValue.(types.Map); ok {
			for key, value := range mapValue.Elements() {
				if dynamicValue, ok := value.(types.Dynamic); ok {
					innerValue := dynamicValue.UnderlyingValue()
					switch v := innerValue.(type) {
					case types.String:
						definitionMap[key] = v.ValueString()
					case types.Int64:
						definitionMap[key] = v.ValueInt64()
					case types.Float64:
						definitionMap[key] = v.ValueFloat64()
					case types.Bool:
						definitionMap[key] = v.ValueBool()
					}
				}
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

	elements := make(map[string]attr.Value)
	for key, value := range policy.Definition {
		switch v := value.(type) {
		case int64:
			elements[key] = types.DynamicValue(types.Int64Value(v))
		case float64:
			elements[key] = types.DynamicValue(types.Float64Value(v))
		case bool:
			elements[key] = types.DynamicValue(types.BoolValue(v))
		case string:
			elements[key] = types.DynamicValue(types.StringValue(v))
		}
	}
	definitionMap, diags := types.MapValue(types.DynamicType, elements)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	state.Definition = types.DynamicValue(definitionMap)

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
		underlyingValue := plan.Definition.UnderlyingValue()

		if mapValue, ok := underlyingValue.(types.Map); ok {
			for key, value := range mapValue.Elements() {
				if dynamicValue, ok := value.(types.Dynamic); ok {
					innerValue := dynamicValue.UnderlyingValue()
					switch v := innerValue.(type) {
					case types.String:
						definitionMap[key] = v.ValueString()
					case types.Int64:
						definitionMap[key] = v.ValueInt64()
					case types.Float64:
						definitionMap[key] = v.ValueFloat64()
					case types.Bool:
						definitionMap[key] = v.ValueBool()
					case types.Number:
						if bigFloat := v.ValueBigFloat(); bigFloat != nil {
							if intVal, accuracy := bigFloat.Int64(); accuracy == big.Exact {
								definitionMap[key] = intVal
							} else if floatVal, accuracy := bigFloat.Float64(); accuracy == big.Exact {
								definitionMap[key] = floatVal
							} else {
								definitionMap[key] = bigFloat
							}
						}
					}
				}
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

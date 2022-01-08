package provider

import (
	"context"
	"fmt"
	"github.com/chris922/terraform-provider-openhab/internal/api"
	"github.com/chris922/terraform-provider-openhab/internal/provider/util"
	"github.com/chris922/terraform-provider-openhab/internal/provider/validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ItemResourceType struct{}

func (t ItemResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "OpenHAB Item",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Resource ID",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Item name",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
					tfsdk.RequiresReplace(),
				},
			},
			"type": {
				MarkdownDescription: "Item type",
				Required:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
				Validators: []tfsdk.AttributeValidator{
					validator.ItemTypeValidator(),
				},
			},
			"label": {
				MarkdownDescription: "Item label",
				Required:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"category": {
				MarkdownDescription: "Item category (often used as the icon)",
				Optional:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"tags": {
				MarkdownDescription: "Item tags",
				Optional:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.ListType{ElemType: types.StringType},
			},
			"group_names": {
				MarkdownDescription: "Item groups",
				Optional:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
				Type: types.ListType{ElemType: types.StringType},
			},
		},
	}, nil
}

func (t ItemResourceType) NewResource(_ context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := ConvertProviderType(in)

	return itemResource{
		client: provider.Client,
	}, diags
}

type itemResourceData struct {
	Id types.String `tfsdk:"id"`

	// required
	Name  types.String `tfsdk:"name"`
	Label types.String `tfsdk:"label"`
	Type  types.String `tfsdk:"type"`

	// optional
	Category   types.String `tfsdk:"category"`
	Tags       types.List   `tfsdk:"tags"`
	GroupNames types.List   `tfsdk:"group_names"`
}

type itemResource struct {
	client *api.Client
}

func (r itemResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data itemResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := api.AddOrUpdateItemInRegistryJSONRequestBody{
		Name:  util.TypeToString(data.Name),
		Label: util.TypeToString(data.Label),
		Type:  util.TypeToString(data.Type),

		Category:   util.TypeToString(data.Category),
		Tags:       util.TypeToStringArray(data.Tags),
		GroupNames: util.TypeToStringArray(data.GroupNames),
	}
	apiResp, err := r.client.AddOrUpdateItemInRegistry(ctx, data.Name.Value,
		&api.AddOrUpdateItemInRegistryParams{}, body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	if apiResp.StatusCode == 200 {
		resp.Diagnostics.AddWarning("Create Item Warning",
			fmt.Sprintf("Item %s was not created, but updated", data.Name.Value))
	} else if apiResp.StatusCode != 201 {
		resp.Diagnostics.AddError("Create Item Error",
			fmt.Sprintf("Unable to create item, got status: %s", apiResp.Status))
		return
	}

	apiRespObj := &api.EnrichedItemDTO{}
	err = api.ReadResponseBody(apiResp, apiRespObj)
	if err != nil {
		resp.Diagnostics.AddError("Create Item Error",
			fmt.Sprintf("Unable to read response of create item action, got error: %s", err))
		return
	}

	// store enriched item to resource
	enrichedItemToData(data, apiRespObj)

	tflog.Trace(ctx, "created an Item resource", "name", data.Name.Value)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r itemResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data itemResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.GetItemByName(ctx, data.Name.Value,
		&api.GetItemByNameParams{})
	if err != nil {
		resp.Diagnostics.AddError("Read Item Error",
			fmt.Sprintf("Unable to read item, got error: %s", err))
		return
	}

	if apiResp.StatusCode == 404 {
		tflog.Debug(ctx, "Item not found, will be removed from state", "name", data.Name.Value)

		resp.State.RemoveResource(ctx)
		return
	}
	if apiResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Read Item Error",
			fmt.Sprintf("Unknown error reading item, got status: %s", apiResp.Status))
		return
	}

	apiRespObj := &api.EnrichedItemDTO{}
	err = api.ReadResponseBody(apiResp, apiRespObj)
	if err != nil {
		resp.Diagnostics.AddError("Read Item Error",
			fmt.Sprintf("Unable to read response of read item action, got error: %s", err))
		return
	}

	// store enriched item to resource
	enrichedItemToData(data, apiRespObj)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r itemResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data itemResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := api.AddOrUpdateItemInRegistryJSONRequestBody{
		Name:  &data.Name.Value,
		Label: &data.Label.Value,
		Type:  &data.Type.Value,

		Category:   util.TypeToString(data.Category),
		Tags:       util.TypeToStringArray(data.Tags),
		GroupNames: util.TypeToStringArray(data.GroupNames),
	}
	apiResp, err := r.client.AddOrUpdateItemInRegistry(ctx, data.Name.Value,
		&api.AddOrUpdateItemInRegistryParams{}, body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	if apiResp.StatusCode == 201 {
		resp.Diagnostics.AddWarning("Update Item Warning",
			fmt.Sprintf("Item %s was not updated, but created", data.Name.Value))
	} else if apiResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Update Item Error",
			fmt.Sprintf("Unable to update item, got status: %s", apiResp.Status))
		return
	}

	apiRespObj := &api.EnrichedItemDTO{}
	err = api.ReadResponseBody(apiResp, apiRespObj)
	if err != nil {
		resp.Diagnostics.AddError("Update Item Error",
			fmt.Sprintf("Unable to read response of update item action, got error: %s", err))
		return
	}

	// store enriched item to resource
	enrichedItemToData(data, apiRespObj)

	tflog.Trace(ctx, "updated an Item resource", "name", data.Name.Value)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r itemResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data itemResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.RemoveItemFromRegistry(ctx, data.Name.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create client, got error: %s", err))
		return
	}

	if apiResp.StatusCode == 404 {
		tflog.Debug(ctx, "Planned to remove an item, but it was already removed", "name", data.Name.Value)
	} else if apiResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Delete Item Error",
			fmt.Sprintf("Unable to delete item, got status: %s", apiResp.Status))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r itemResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("name"), req, resp)
}

func enrichedItemToData(data itemResourceData, apiRespObj *api.EnrichedItemDTO) {
	data.Id = util.StringToType(apiRespObj.Name)
	data.Name = util.StringToType(apiRespObj.Name)
	data.Label = util.StringToType(apiRespObj.Label)
	data.Type = util.StringToType(apiRespObj.Type)

	data.Category = util.StringToType(apiRespObj.Category)
	data.Tags = util.StringArrayToType(apiRespObj.Tags)
	data.GroupNames = util.StringArrayToType(apiRespObj.GroupNames)
}

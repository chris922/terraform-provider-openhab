package provider

import (
	"context"
	"fmt"
	"github.com/chris922/terraform-provider-openhab/internal/api"
	"github.com/chris922/terraform-provider-openhab/internal/provider/util"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type LinkResourceType struct{}

func (t LinkResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "OpenHAB Link between an Item and a Channel.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Resource ID",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"item_name": {
				MarkdownDescription: "Item name",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
					tfsdk.RequiresReplace(),
				},
			},
			"channel_uid": {
				MarkdownDescription: "Channel UID",
				Required:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
					tfsdk.RequiresReplace(),
				},
				Type: types.StringType,
			},
			"configuration": {
				MarkdownDescription: "Link configuration",
				Optional:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
					tfsdk.RequiresReplace(),
				},
				// TODO: can the values also be e.g. a number?
				Type: types.MapType{ElemType: types.StringType},
			},
		},
	}, nil
}

func (t LinkResourceType) NewResource(_ context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := ConvertProviderType(in)

	return linkResource{
		client: provider.Client,
	}, diags
}

type linkResourceData struct {
	Id types.String `tfsdk:"id"`

	// required
	ItemName   types.String `tfsdk:"item_name"`
	ChannelUid types.String `tfsdk:"channel_uid"`

	// optional
	Configuration types.Map `tfsdk:"configuration"`
}

type linkResource struct {
	client *api.Client
}

func (r linkResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data linkResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := api.LinkItemToChannelJSONRequestBody{
		ItemName:   &data.ItemName.Value,
		ChannelUID: &data.ChannelUid.Value,

		Configuration: util.TypeToStringMap(data.Configuration),
	}
	apiResp, err := r.client.LinkItemToChannel(ctx, data.ItemName.Value, data.ChannelUid.Value, body)
	if err != nil {
		resp.Diagnostics.AddError("Create Link Error",
			fmt.Sprintf("Unable to create link, got error: %s", err))
		return
	}

	if apiResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Create Link Error",
			fmt.Sprintf("Unable to create link, got status: %s", apiResp.Status))
		return
	}

	// generate ID out of item name + channel uid
	data.Id = types.String{Value: data.ItemName.Value + "-" + data.ChannelUid.Value}

	tflog.Trace(ctx, "created a Link resource", "item_name", data.ItemName.Value,
		"channel_uid", data.ChannelUid.Value)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r linkResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data linkResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.GetItemLink(ctx, data.ItemName.Value, data.ChannelUid.Value)
	if err != nil {
		resp.Diagnostics.AddError("Read Link Error",
			fmt.Sprintf("Unable to read link, got error: %s", err))
		return
	}

	if apiResp.StatusCode == 404 {
		tflog.Debug(ctx, "Link not found, will be removed from state", "item_name", data.ItemName.Value,
			"channel_uid", data.ChannelUid.Value)

		resp.State.RemoveResource(ctx)
		return
	}
	if apiResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Read Link Error",
			fmt.Sprintf("Unknown error reading link, got status: %s", apiResp.Status))
		return
	}

	apiRespObj := &api.EnrichedItemChannelLinkDTO{}
	err = api.ReadResponseBody(apiResp, apiRespObj)
	if err != nil {
		resp.Diagnostics.AddError("Read Link Error",
			fmt.Sprintf("Unable to read response of read link action, got error: %s", err))
		return
	}

	// store enriched link to resource
	enrichedItemChannelLinkToData(data, apiRespObj)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r linkResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data linkResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Update Link Error", "Updating links is not supported")
}

func (r linkResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data linkResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.UnlinkItemFromChannel(ctx, data.ItemName.Value, data.ChannelUid.Value)
	if err != nil {
		resp.Diagnostics.AddError("Delete Link Error",
			fmt.Sprintf("Unable to delete link, got error: %s", err))
		return
	}

	if apiResp.StatusCode == 404 {
		tflog.Debug(ctx, "Planned to remove a link, but it was already removed",
			"item_name", data.ItemName.Value, "channel_uid", data.ChannelUid.Value)
	} else if apiResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Delete Link Error",
			fmt.Sprintf("Unable to delete link, got status: %s", apiResp.Status))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r linkResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

func generateLinkResourceId(itemName string, channelUid string) string {
	return itemName + "-" + channelUid
}

func enrichedItemChannelLinkToData(data linkResourceData, apiRespObj *api.EnrichedItemChannelLinkDTO) {
	id := generateLinkResourceId(*apiRespObj.ItemName, *apiRespObj.ChannelUID)

	data.Id = util.StringToType(&id)
	data.ItemName = util.StringToType(apiRespObj.ItemName)
	data.ChannelUid = util.StringToType(apiRespObj.ChannelUID)
	data.Configuration = util.StringMapToType(apiRespObj.Configuration)
}

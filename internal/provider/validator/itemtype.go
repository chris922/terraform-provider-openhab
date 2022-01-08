package validator

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type itemTypeValidator struct {
	tfsdk.AttributeValidator
}

func ItemTypeValidator() *itemTypeValidator {
	return &itemTypeValidator{}
}

func (v itemTypeValidator) Description(ctx context.Context) string {
	return "Ensures a given type is valid."
}

func (v itemTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

/**
Available Item types are:
	Color	Color information (RGB)	OnOff, IncreaseDecrease, Percent, HSB
	Contact	Status of contacts, e.g. door/window contacts. Does not accept commands, only status updates.	OpenClosed
	DateTime	Stores date and time	-
	Dimmer	Percentage value for dimmers	OnOff, IncreaseDecrease, Percent
	Group	Item to nest other items / collect them in groups	-
	Image	Binary data of an image	-
	Location	GPS coordinates	Point
	Number	Values in number format	Decimal
	Player	Allows control of players (e.g. audio players)	PlayPause, NextPrevious, RewindFastforward
	Rollershutter	Roller shutter Item, typically used for blinds	UpDown, StopMove, Percent
	String	Stores texts	String
	Switch	Switch Item, used for anything that needs to be switched ON and OFF	OnOff

See: https://www.openhab.org/docs/configuration/items.html#type
*/

func (v itemTypeValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	fullValue := request.AttributeConfig.(types.String).Value
	valueParts := strings.Split(fullValue, ":")

	baseType := valueParts[0]
	switch baseType {
	case "Color":
	case "Contact":
	case "DateTime":
	case "Dimmer":
	case "Group":
	case "Image":
	case "Location":
	case "Number":
	case "Player":
	case "Rollershutter":
	case "String":
	case "Switch":
		break
	default:
		response.Diagnostics.AddAttributeError(request.AttributePath, "Unknown type", fmt.Sprintf("Given type '%s' is unknown.", fullValue))
	}
}

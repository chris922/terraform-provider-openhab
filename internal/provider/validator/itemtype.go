package validator

import (
	"context"
	"fmt"
	"github.com/chris922/terraform-provider-openhab/internal/provider/util"
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

See: https://www.openhab.org/docs/concepts/items.html
*/

// see https://www.openhab.org/docs/concepts/units-of-measurement.html#list-of-units
var imperialUnits = []string{
	"Pressure",
	"Temperature",
	"Speed",
	"Length",
}

// see https://www.openhab.org/docs/concepts/units-of-measurement.html#list-of-units
var siUnits = []string{
	"Acceleration",
	"AmountOfSubstance",
	"Angle",
	"Area",
	"ArealDensity",
	"CatalyticActivity",
	"DataAmount",
	"DataTransferRate",
	"Density",
	"Dimensionless",
	"ElectricPotential",
	"ElectricCapacitance",
	"ElectricCharge",
	"ElectricConductance",
	"ElectricConductivity",
	"ElectricCurrent",
	"ElectricInductance",
	"ElectricResistance",
	"Energy",
	"Force",
	"Frequency",
	"Illuminance",
	"Intensity",
	"Length",
	"LuminousFlux",
	"LuminousIntensity",
	"MagneticFlux",
	"MagneticFluxDensity",
	"Mass",
	"Power",
	"Pressure",
	"Radioactivity",
	"RadiationDoseAbsorbed",
	"RadiationDoseEffective",
	"SolidAngle",
	"Speed",
	"Temperature",
	"Time",
	"Volume",
	"VolumetricFlowRate",
}

func (v itemTypeValidator) Validate(_ context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
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
	case "Player":
	case "Rollershutter":
	case "String":
	case "Switch":
		break
	case "Number":
		if len(valueParts) > 1 {
			dimension := valueParts[1]
			if !util.StringArrayContains(imperialUnits, dimension) && !util.StringArrayContains(siUnits, dimension) {
				response.Diagnostics.AddAttributeError(request.AttributePath,
					"Unknown dimension used for Number type",
					fmt.Sprintf("Given dimension '%s' is unknown.", dimension))
			}
		}
		break
	default:
		response.Diagnostics.AddAttributeError(request.AttributePath, "Unknown type", fmt.Sprintf("Given type '%s' is unknown.", fullValue))
	}
}

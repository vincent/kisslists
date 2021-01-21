package pkg

type Theme struct {
	Background        string
	Title             string
	BoxGradientStart  string
	BoxGradientEnd    string
	LineGradientStart string
	LineGradientEnd   string
	InputCheck        string
	InputChecked      string
	ItemBackground    string
	ItemText          string
	ItemCheck         string
}

type Themes map[string]Theme

var KnownThemes = Themes{

	"default": Theme{
		Background:        "#ddd",
		Title:             "darkcyan",
		BoxGradientStart:  "#27FDC7",
		BoxGradientEnd:    "#0FC0F5",
		LineGradientStart: "#0FC0F5",
		LineGradientEnd:   "#27FDC7",
		InputCheck:        "#447b7f",
		InputChecked:      "#135156",
		ItemBackground:    "#FFF",
		ItemText:          "#135156",
		ItemCheck:         "#27FDC7",
	},

	"original": Theme{
		Background:        "#ddd",
		Title:             "darkcyan",
		BoxGradientStart:  "#eaf7cf",
		BoxGradientEnd:    "#ebefbf",
		LineGradientStart: "#ceb5a7",
		LineGradientEnd:   "#ceb5a7",
		InputCheck:        "#92828d",
		InputChecked:      "#adaabf",
		ItemBackground:    "#FFF",
		ItemText:          "#135156",
		ItemCheck:         "#27FDC7",
	},

	// use https://coolors.co
	"playful": asTheme([]string{"264653", "2a9d8f", "e9c46a", "f4a261", "e76f51"}),
	"ocean":   asTheme([]string{"03045e", "023e8a", "0077b6", "0096c7", "00b4d8", "48cae4", "90e0ef", "ade8f4", "caf0f8"}),
	"nature":  asTheme([]string{"ffe8d6", "b7b7a4", "a5a58d", "6b705c", "cb997e", "ddbea9"}),
}

func asTheme(colors []string) Theme {
	return Theme{
		Background:        "#" + colors[1],
		Title:             "#" + colors[0],
		BoxGradientStart:  "#" + colors[2],
		BoxGradientEnd:    "#" + colors[0],
		LineGradientStart: "#" + colors[1],
		InputCheck:        "#" + colors[3],
		InputChecked:      "#" + colors[0],
		ItemBackground:    "#" + colors[4],
		ItemText:          "#" + colors[0],
		ItemCheck:         "#" + colors[1],
	}
}

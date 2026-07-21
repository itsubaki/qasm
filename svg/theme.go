package svg

type Theme struct {
	Background    string
	Wire          string
	WireLabel     string
	Text          string
	GateFill      string
	GateStroke    string
	MeasureFill   string
	MeasureStroke string
	SubFill       string
	SubStroke     string
	BarrierStroke string
}

var Paper = Theme{
	Background:    "#ffffff",
	Wire:          "#000000",
	WireLabel:     "#000000",
	Text:          "#000000",
	GateFill:      "#ffffff",
	GateStroke:    "#000000",
	MeasureFill:   "#ffffff",
	MeasureStroke: "#000000",
	SubFill:       "#ffffff",
	SubStroke:     "#000000",
	BarrierStroke: "#000000",
}

var Cyber = Theme{
	Background:    "#0d1117",
	Wire:          "#4b5563",
	WireLabel:     "#4b5563",
	Text:          "#e5e7eb",
	GateFill:      "#1f2937",
	GateStroke:    "#0ea5e9",
	MeasureFill:   "#1f2937",
	MeasureStroke: "#10b981",
	SubFill:       "#1f2937",
	SubStroke:     "#8b5cf6",
	BarrierStroke: "#f59e0b",
}

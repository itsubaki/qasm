package svg

import (
	"fmt"
	"sort"
	"strings"
)

type Config struct {
	WireGap    int
	WireStartX int
	WireStartY int
	OpWidth    int
	OpHeight   int
	OpRX       int
	FontSize   int
	Theme      Theme
}

var DefaultConfig = Config{
	WireGap:    56,
	WireStartX: 80,
	WireStartY: 42,
	OpWidth:    36,
	OpHeight:   36,
	OpRX:       8,
	FontSize:   13,
	Theme:      Paper,
}

func Render(layout *Layout, config Config) string {
	// the size of the SVG canvas
	width := config.WireStartX + len(layout.Layers)*config.WireGap + config.WireGap/2
	height := config.WireStartY + len(layout.Wires)*config.WireGap

	// svg
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`,
		width,
		height,
		width,
		height,
	)

	fmt.Fprintf(&b, `<rect x="0" y="0" width="%d" height="%d" fill="%s"/>`,
		width,
		height,
		config.Theme.Background,
	)

	// style
	b.WriteString(`<style>`)
	fmt.Fprintf(&b, `.gate-label { font-family: ui-monospace, monospace; font-size: %dpx; font-weight: 600; }`, config.FontSize)
	fmt.Fprintf(&b, `.wire-label { font-family: ui-monospace, monospace; font-size: %dpx; font-weight: 500; }`, config.FontSize)
	b.WriteString(`</style>`)

	// wires
	for i, w := range layout.Wires {
		y := config.WireStartY + i*config.WireGap
		fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2" />`,
			config.WireStartX,
			y,
			width,
			y,
			config.Theme.Wire,
		)

		fmt.Fprintf(&b, `<text x="%d" y="%d" text-anchor="end" fill="%s" class="wire-label">%s</text>`,
			config.WireStartX-8,
			y+5,
			config.Theme.WireLabel,
			w.Name,
		)
	}

	// ops
	x := config.WireStartX + config.WireGap/2
	for _, layer := range layout.Layers {
		for _, op := range layer.Ops {
			switch o := op.(type) {
			case *Gate:
				// wires
				for _, c := range o.Control {
					for _, t := range o.Target {
						cy := config.WireStartY + c*config.WireGap
						ty := config.WireStartY + t*config.WireGap

						fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2" />`,
							x+config.OpWidth/2,
							cy,
							x+config.OpWidth/2,
							ty,
							config.Theme.GateStroke,
						)

						fmt.Fprintf(&b, `<circle cx="%d" cy="%d" r="6" fill="%s" />`,
							x+config.OpWidth/2,
							cy,
							config.Theme.GateStroke,
						)
					}
				}

				// operation box
				for _, t := range o.Target {
					y := config.WireStartY + t*config.WireGap

					fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="%d" rx="%d" fill="%s" stroke="%s" stroke-width="2"/>`,
						x,
						y-config.OpHeight/2,
						config.OpWidth,
						config.OpHeight,
						config.OpRX,
						config.Theme.GateFill,
						config.Theme.GateStroke,
					)

					fmt.Fprintf(&b, `<text x="%d" y="%d" text-anchor="middle" fill="%s" class="gate-label">%s</text>`,
						x+config.OpWidth/2,
						y+config.OpHeight/2-13,
						config.Theme.Text,
						o.Name,
					)
				}
			case *Subroutine:
				// operation box
				minY, maxY := o.Wire[0], o.Wire[0]
				for _, t := range o.Wire {
					minY, maxY = min(minY, t), max(maxY, t)
				}

				topY := config.WireStartY + minY*config.WireGap
				bottomY := config.WireStartY + maxY*config.WireGap
				centerY := (topY + bottomY) / 2
				height := (bottomY - topY) + config.OpHeight

				fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="%d" rx="%d" fill="%s" stroke="%s" stroke-width="2" />`,
					x,
					centerY-height/2,
					config.OpWidth,
					height,
					config.OpRX,
					config.Theme.SubFill,
					config.Theme.SubStroke,
				)

				fmt.Fprintf(&b, `<text x="%d" y="%d" text-anchor="middle" fill="%s" class="gate-label">%s</text>`,
					x+config.OpWidth/2,
					centerY+config.OpHeight/2-13,
					config.Theme.Text,
					o.Name,
				)
			case *Measurement:
				// wires
				for _, w := range o.Wire {
					for _, t := range o.Target {
						cy := config.WireStartY + w*config.WireGap
						ty := config.WireStartY + t*config.WireGap

						fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2" />`,
							x+config.OpWidth/2,
							cy+config.OpHeight/2,
							x+config.OpWidth/2,
							ty-4,
							config.Theme.MeasureStroke,
						)

						fmt.Fprintf(&b, `<polygon points="%d,%d %d,%d %d,%d" fill="%s" />`,
							x+config.OpWidth/2,
							ty,
							x+config.OpWidth/2-4,
							ty-6,
							x+config.OpWidth/2+4,
							ty-6,
							config.Theme.MeasureStroke,
						)
					}
				}

				// operation box
				for _, w := range o.Wire {
					y := config.WireStartY + w*config.WireGap
					fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="%d" rx="%d" fill="%s" stroke="%s" stroke-width="2" />`,
						x,
						y-config.OpHeight/2,
						config.OpWidth,
						config.OpHeight,
						config.OpRX,
						config.Theme.MeasureFill,
						config.Theme.MeasureStroke,
					)

					fmt.Fprintf(&b, `<path d="M %d %d A 10 10 0 0 1 %d %d" fill="none" stroke="%s" stroke-width="2" />`,
						x+10,
						y,
						x+30,
						y,
						config.Theme.MeasureStroke,
					)

					fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2" />`,
						x+20,
						y,
						x+26,
						y-10,
						config.Theme.MeasureStroke,
					)
				}
			case *Barrier:
				wires := make([]int, len(o.Wire))
				copy(wires, o.Wire)
				sort.Ints(wires)

				start, prev := wires[0], wires[0]
				for i := 1; i <= len(wires); i++ {
					if i == len(wires) || wires[i] != prev+1 {
						y1 := config.WireStartY + start*config.WireGap
						y2 := config.WireStartY + prev*config.WireGap

						fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2" stroke-dasharray="4 2" />`,
							x+config.OpWidth/2,
							y1-config.OpHeight/2,
							x+config.OpWidth/2,
							y2+config.OpHeight/2,
							config.Theme.BarrierStroke,
						)

						if i < len(wires) {
							start, prev = wires[i], wires[i]
						}

						continue
					}

					prev = wires[i]
				}
			}
		}

		// next layer
		x += config.WireGap
	}

	b.WriteString(`</svg>`)
	return b.String()
}

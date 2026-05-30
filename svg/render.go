package svg

import (
	"fmt"
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
}

var DefaultConfig = Config{
	WireGap:    56,
	WireStartX: 80,
	WireStartY: 42,
	OpWidth:    36,
	OpHeight:   36,
	OpRX:       8,
	FontSize:   13,
}

func Render(layout *Layout, config Config) string {
	// the size of the SVG canvas
	width := config.WireStartX + len(layout.Layers)*config.WireGap + config.WireGap/2
	height := config.WireStartY + len(layout.Wires)*config.WireGap

	// svg
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`,
		width, height,
		width, height,
	)

	// style
	b.WriteString(`<style>`)
	fmt.Fprintf(&b, `.gate-label { font-family: ui-monospace, monospace; font-size: %dpx; font-weight: 600; }`, config.FontSize)
	fmt.Fprintf(&b, `.wire-label { font-family: ui-monospace, monospace; font-size: %dpx; font-weight: 500; }`, config.FontSize)
	b.WriteString(`</style>`)

	// wires
	for i, w := range layout.Wires {
		y := config.WireStartY + i*config.WireGap
		fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#4b5563" stroke-width="2" />`,
			config.WireStartX, y,
			width, y,
		)

		fmt.Fprintf(&b, `<text x="%d" y="%d" text-anchor="end" fill="#4b5563" class="wire-label">%s</text>`,
			config.WireStartX-8, y+5,
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
				for _, c := range o.Controls {
					for _, t := range o.Targets {
						cy := config.WireStartY + c*config.WireGap
						ty := config.WireStartY + t*config.WireGap

						fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#0ea5e9" stroke-width="2" />`,
							x+config.OpWidth/2, cy,
							x+config.OpWidth/2, ty,
						)

						fmt.Fprintf(&b, `<circle cx="%d" cy="%d" r="6" fill="#0ea5e9" />`,
							x+config.OpWidth/2, cy,
						)
					}
				}

				// operation box
				minY, maxY := o.Targets[0], o.Targets[0]
				for _, t := range o.Targets {
					minY, maxY = min(minY, t), max(maxY, t)
				}

				topY := config.WireStartY + minY*config.WireGap
				bottomY := config.WireStartY + maxY*config.WireGap
				centerY := (topY + bottomY) / 2
				height := (bottomY - topY) + config.OpHeight

				fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="%d" rx="%d" fill="#1f2937" stroke="#0ea5e9" stroke-width="2" />`,
					x,
					centerY-height/2,
					config.OpWidth,
					height,
					config.OpRX,
				)

				fmt.Fprintf(&b, `<text x="%d" y="%d" text-anchor="middle" fill="#e5e7eb" class="gate-label">%s</text>`,
					x+config.OpWidth/2, centerY+config.OpHeight/2-13,
					o.Name,
				)
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

				fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="%d" rx="%d" fill="#1f2937" stroke="#8b5cf6" stroke-width="2" />`,
					x,
					centerY-height/2,
					config.OpWidth,
					height,
					config.OpRX,
				)

				fmt.Fprintf(&b, `<text x="%d" y="%d" text-anchor="middle" fill="#e5e7eb" class="gate-label">%s</text>`,
					x+config.OpWidth/2, centerY+config.OpHeight/2-13,
					o.Name,
				)
			case *Measurement:
				// wires
				for _, w := range o.Wire {
					for _, t := range o.Targets {
						cy := config.WireStartY + w*config.WireGap
						ty := config.WireStartY + t*config.WireGap

						fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#10b981" stroke-width="2" />`,
							x+config.OpWidth/2, cy+config.OpHeight/2,
							x+config.OpWidth/2, ty-4,
						)

						fmt.Fprintf(&b, `<polygon points="%d,%d %d,%d %d,%d" fill="#10b981" />`,
							x+config.OpWidth/2, ty,
							x+config.OpWidth/2-4, ty-6,
							x+config.OpWidth/2+4, ty-6,
						)
					}
				}

				// operation box
				for _, w := range o.Wire {
					y := config.WireStartY + w*config.WireGap
					fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="%d" rx="%d" fill="#1f2937" stroke="#10b981" stroke-width="2" />`,
						x, y-config.OpHeight/2,
						config.OpWidth, config.OpHeight, config.OpRX,
					)

					fmt.Fprintf(&b, `<path d="M %d %d A 10 10 0 0 1 %d %d" fill="none" stroke="#10b981" stroke-width="2" />`,
						x+10, y,
						x+30, y,
					)

					fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#10b981" stroke-width="2" />`,
						x+20, y,
						x+26, y-10,
					)
				}
			case *Barrier:
				for _, w := range o.Wire {
					y := config.WireStartY + w*config.WireGap
					fmt.Fprintf(&b, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#f59e0b" stroke-width="2" stroke-dasharray="4 2" />`,
						x+config.OpWidth/2, y-config.OpHeight/2,
						x+config.OpWidth/2, y+config.OpHeight/2,
					)
				}
			}
		}

		// next layer
		x += config.WireGap
	}

	b.WriteString(`</svg>`)
	return b.String()
}

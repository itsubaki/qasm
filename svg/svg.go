package svg

import xparser "github.com/itsubaki/qasm/parser"

func SVG(text string, config Config) (string, error) {
	program, err := xparser.Parse(text)
	if err != nil {
		return "", err
	}

	circuit, err := NewVisitor().Run(program)
	if err != nil {
		return "", err
	}

	layout := NewLayout(circuit)
	diagram := Render(layout, config)
	return diagram, nil
}

package visitor

import (
	"testing"

	"github.com/itsubaki/q"
	"github.com/itsubaki/qasm/environ"
	xparser "github.com/itsubaki/qasm/parser"
)

func TestEndSignal(t *testing.T) {
	cases := []struct {
		name string
		text string
	}{
		{name: "statement", text: `end;`},
		{name: "nested scope", text: `{ { end; } touch(q); }`},
		{name: "if", text: `if (true) { end; }`},
		{name: "else", text: `if (false) {} else { end; }`},
		{name: "for body", text: `for int i in [0:1] { end; }`},
		{name: "while body", text: `while (true) { end; }`},
		{name: "switch body", text: `switch (1) { case 1 { end; } }`},
		{name: "function", text: `stop();`},
		{name: "nested return", text: `def f() -> int { return stop(); } f();`},
		{name: "initializer", text: `int result=stop();`},
		{name: "constant", text: `const int result=stop();`},
		{name: "array literal", text: `array[int[64], 2] result={stop(),touch(q)};`},
		{name: "cast", text: `int(stop());`},
		{name: "additive left", text: `stop()+touch(q);`},
		{name: "additive right", text: `1+stop();`},
		{name: "multiplicative", text: `stop()*touch(q);`},
		{name: "unary", text: `-stop();`},
		{name: "power", text: `stop()**touch(q);`},
		{name: "equality", text: `stop()==touch(q);`},
		{name: "comparison", text: `stop()<touch(q);`},
		{name: "logical and", text: `true && stop();`},
		{name: "logical or", text: `false || stop();`},
		{name: "bitwise", text: `1 & stop();`},
		{name: "shift", text: `1 << stop();`},
		{name: "argument", text: `mod(stop(),touch(q));`},
		{name: "if condition", text: `if (stop()) { touch(q); }`},
		{name: "while condition", text: `while (stop()) { touch(q); }`},
		{name: "for range", text: `for int i in [0:stop()] { touch(q); }`},
		{name: "switch expression", text: `switch (stop()) { default { touch(q); } }`},
		{name: "switch case", text: `switch (1) { case stop(),touch(q) {} }`},
		{name: "gate argument", text: `U(stop(),touch(q),0) q;`},
		{name: "global phase", text: `gphase(stop());`},
		{name: "gate power", text: `pow(stop()) @ U(0,0,0) q;`},
		{name: "gate control", text: `ctrl(stop()) @ U(0,0,0) q;`},
		{name: "user gate argument", text: `gate g(theta) a { U(theta,0,0) a; } g(stop()) q;`},
		{name: "nested gate", text: `gate g a { U(stop(),0,0) a; } gate h a { g a; } h q;`},
		{name: "index", text: `q[stop()];`},
		{name: "index range", text: `q[0:stop()];`},
		{name: "reset operand", text: `reset q[stop()];`},
		{name: "measurement operand", text: `measure q[stop()];`},
		{name: "assignment index", text: `bit[1] result; result[stop()]=touch(q);`},
		{name: "measurement target", text: `bit[1] result; measure q -> result[stop()];`},
		{name: "qubit size", text: `qubit[stop()] result;`},
		{name: "array size", text: `array[int[64], stop()] result;`},
		{name: "array element size", text: `array[int[stop()], 1] result;`},
	}

	for _, c := range cases {
		text := `
				int touched=0; qubit q;
				def stop() -> int { end; }
				def touch(qubit arg) -> int { U(pi,0,pi) arg; return 1; }
			` + c.text + ` touched=1; U(pi,0,pi) q;`

		program, err := xparser.Parse(text)
		if err != nil {
			t.Fatal(err)
		}

		qsim, env := q.New(), environ.New()
		result := New(qsim, env).Visit(program)
		if _, ok := result.(endSignal); !ok {
			t.Fatalf("got=%v (%T), want endSignal", result, result)
		}

		if got := env.Variable["touched"]; got != int64(0) {
			t.Errorf("executed statement after end: touched=%v", got)
		}

		if !qsim.Qubit().IsZero() {
			t.Errorf("executed gate after end: %v", qsim.Qubit())
		}

		if result, err := Visit(text); err != nil || result != End {
			t.Errorf("Visit: got=(%v, %v), want=(%s, nil)", result, err, End)
		}
	}
}

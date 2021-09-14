gate X q { U(pi, 0.0, pi) q; }
gate H q { U(pi/2.0, 0.0, pi) q; }

gate CX q, p { ctrl @ X q, p; }

gate bell q, p { h q; cx q, p; }


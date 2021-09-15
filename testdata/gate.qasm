gate x q { U(pi, 0.0, pi) q; }
gate h q { U(pi/2.0, 0.0, pi) q; }

gate cx   q, p { ctrl @ x q, p; }
gate bell q, p { h q; cx q, p; }

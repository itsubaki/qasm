gate i q { U(0, 0, 0) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate y q { U(pi, pi/2.0, pi/2.0) q; }
gate z q { Z q; }

gate cx c, t { ctrl @ x c, t; }
gate cy c, t { ctrl @ y c, t; }
gate cz c, t { ctrl @ z c, t; }
gate ch c, t { ctrl @ h c, t; }

gate ccx c0, c1, t { ctrl @ ctrl @ x c0, c1, t; }
gate swap q, p { cx q, p; cx p, q; cx q, p; }

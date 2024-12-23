OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ x c, t; }

qubit[2] q;
reset q;

h q[0];
cx q[0], q[1];

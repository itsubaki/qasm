OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }
gate cx c, t { ctrl @ U(pi, 0, pi) c, t; }
gate cz c, t { ctrl @ U(0, pi, 0) c, t; }

qubit[3] q;

h q[1];
cx q[1], q[2];
cx q[0], q[1];
h q[0];

cx q[1], q[2];
cz q[0], q[2];

measure q[0];
measure q[1];

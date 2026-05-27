OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }

qubit[3] q;

h q[1];
ctrl @ x q[1], q[2];
ctrl @ x q[0], q[1];
h q[0];

measure q;

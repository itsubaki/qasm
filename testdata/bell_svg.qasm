OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate x q { U(pi, 0, pi) q; }

qubit[2] q;
bit[2] c;
reset q;

h q[0];
ctrl @ x q[0], q[1];

c = measure q;

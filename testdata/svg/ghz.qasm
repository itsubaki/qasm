OPENQASM 3.0;

qubit[3] q;
reset q;

h q[0];
ctrl @ x q[0], q[1];
ctrl @ x q[1], q[2];

measure q;

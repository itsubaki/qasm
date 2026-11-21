OPENQASM 3.0;

qubit[2] q;
U(pi/2, 0, pi) q[0];
ctrl @ U(pi, 0, pi) q[0], q[1];

// Measurement in a declaration.
bit[2] init = measure q;

// Measurement in an assignment.
bit assigned;
assigned = measure q[0];

// Equivalent arrow syntax.
bit arrow;
measure q[1] -> arrow;

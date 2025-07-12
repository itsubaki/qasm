OPENQASM 3.0;

qubit[2] q;
U(pi/2.0, 0, pi)    q[0];
ctrl @ U(pi, 0, pi) q[0], q[1];

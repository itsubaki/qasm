OPENQASM 3.0;
qubit[2] q;
bit[2] c;

gphase(pi/4);
h q[0];
ctrl @ gphase(pi/2) q[0];
ctrl(2) @ gphase(pi/3) q[0], q[1];
h q[1];

measure q -> c;

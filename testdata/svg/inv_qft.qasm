OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate r(theta) q { U(0, 0, theta) q; }

qubit[7] q;
reset q;

// inv_qft
h q[6];

ctrl @ r(-pi/2) q[6], q[5];
h q[5];

ctrl @ r(-pi/4) q[6], q[4];
ctrl @ r(-pi/2) q[5], q[4];
h q[4];

ctrl @ r(-pi/8) q[6], q[3];
ctrl @ r(-pi/4) q[5], q[3];
ctrl @ r(-pi/2) q[4], q[3];
h q[3];

ctrl @ r(-pi/16) q[6], q[2];
ctrl @ r(-pi/ 8) q[5], q[2];
ctrl @ r(-pi/ 4) q[4], q[2];
ctrl @ r(-pi/ 2) q[3], q[2];
h q[2];

ctrl @ r(-pi/32) q[6], q[1];
ctrl @ r(-pi/16) q[5], q[1];
ctrl @ r(-pi/ 8) q[4], q[1];
ctrl @ r(-pi/ 4) q[3], q[1];
ctrl @ r(-pi/ 2) q[2], q[1];
h q[1];

ctrl @ r(-pi/64) q[6], q[0];
ctrl @ r(-pi/32) q[5], q[0];
ctrl @ r(-pi/16) q[4], q[0];
ctrl @ r(-pi/ 8) q[3], q[0];
ctrl @ r(-pi/ 4) q[2], q[0];
ctrl @ r(-pi/ 2) q[1], q[0];
h q[0];

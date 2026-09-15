OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate r(theta) q { U(0, 0, theta) q;}

// N=15, a=7
qubit[3] q;
qubit[4] a;
reset q;
reset a;

h q;
x a[3];
barrier q, a;

// modexp
ctrl @ x q[0], a[1];
ctrl @ x q[0], a[2];
barrier q, a;

ctrl @ x          a[0], a[2];
ctrl(2) @ x q[1], a[2], a[0];
ctrl @ x          a[0], a[2];

ctrl @ x          a[3], a[1];
ctrl(2) @ x q[1], a[1], a[3];
ctrl @ x          a[3], a[1];
barrier q, a;

// inv_qft
h q[2];

ctrl @ r(-pi/2) q[2], q[1];
h q[1];

ctrl @ r(-pi/4) q[2], q[0];
ctrl @ r(-pi/2) q[1], q[0];
h q[0];
barrier q, a;

measure q;


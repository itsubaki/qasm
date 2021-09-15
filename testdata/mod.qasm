OPENQASM 3.0;
include "testdata/gate.qasm";

gate x1(a) q     { pow(3) @ U(pi, a, pi) q; }
gate x2(a) q     { inv    @ x1(a) q;}
gate x3(a) q     { inv    @ x2(a) q;}
gate x4(a) q     { inv    @ x3(a) q;}
gate cx1(a) q, p { ctrl   @ x4(a) q, p; }

qubit[2] q;
reset q;

h     q[0];
cx1(0) q[0], q[1];
print;

gate ncx  q0, q1     { negctrl @ x   q0, q1; }
gate cncx q0, q1, q2 { ctrl    @ ncx q0, q1, q2; }

qubit p;
reset q, p;
print;

x q[0];
cncx q[0], q[1], p;

OPENQASM 3.0;
include "testdata/gate.qasm";

gate X1(a) q    { pow(3) @ U(pi, a, pi) q; }
gate X2(a) q    { inv    @ X1(a) q;}
gate X3(a) q    { inv    @ X2(a) q;}
gate X4(a) q    { inv    @ X3(a) q;}
gate CX(a) q, p { ctrl   @ X4(a) q, p; }

qubit[2] q;
reset q;

H     q[0];
CX(0) q[0], q[1];
print;

gate nCX  q0, q1     { negctrl @ X   q0, q1; }
gate CnCX q0, q1, q2 { ctrl    @ nCX q0, q1, q2; }

qubit p;
reset q, p;
print;

X q[0];
CnCX q[0], q[1], p;

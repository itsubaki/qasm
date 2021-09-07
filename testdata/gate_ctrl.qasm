OPENQASM 3.0;
include "itsubaki/q.qasm";

gate X(a, b)    q0         { U(a, b, a)  q0; }
gate nCX(a, b)  q0, q1     { negctrl @ X(a, b)   q0, q1; }
gate CnCX(a, b) q0, q1, q2 { ctrl    @ nCX(a, b) q0, q1, q2; }

qubit[3] q;
reset q;

X(pi, 0) q[0];
CnCX(pi, 0) q[0], q[1], q[2];

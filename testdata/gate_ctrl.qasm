OPENQASM 3.0;
include "itsubaki/q.qasm";

gate X    q0         { U(pi, 0, pi)  q0; }
gate nCX  q0, q1     { negctrl @ X   q0, q1; }
gate CnCX q0, q1, q2 { ctrl    @ nCX q0, q1, q2; }

qubit[3] q;
reset q;

X q[0];
CnCX q[0], q[1], q[2];

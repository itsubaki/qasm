gate X   q0         { U(pi, 0, pi) q0; }
gate CX  q0, q1     { negctrl @ X  q0, q1; }
gate CnCX q0, q1, q2 { ctrl    @ CX q0, q1, q2; }

qubit[3] q;
reset q;

X q[0];
CnCX q[0], q[1], q[2];

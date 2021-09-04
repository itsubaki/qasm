gate X   q0         { U(pi, 0, pi) q0; }
gate CX  q0, q1     { ctrl @ X  q0, q1; }
gate CCX q0, q1, q2 { ctrl @ CX q0, q1, q2; }

qubit[3] q;
reset q;

X q[0], q[1];
CCX q[0], q[1], q[2];

gate X q { U(pi, 0.0, pi) q; }
gate X2 q { X q; }
gate CX q, p { ctrl @ x q, p; }

qubit[2] q;

h q[0];
CX q[0], q[1];

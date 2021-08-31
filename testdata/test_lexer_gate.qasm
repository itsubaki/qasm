gate X   q    { U(pi, 0.0, pi) q; }
gate X2  q    { X q; }
gate CX  q, p { ctrl @ X q, p; }
gate CX2 q, p { ctrl @ X2 q, p; }

qubit[2] q;

h   q[0];
CX2 q[0], q[1];

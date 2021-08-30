gate X q { U(pi, 0.0, pi) q; }
gate X2 q { X q; }
gate CX q, p { ctrl @ X2 q, p; }

qubit[2] q;

X2 q[0];

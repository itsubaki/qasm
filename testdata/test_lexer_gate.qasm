gate X q { U(pi, 0.0, pi) q; }
gate X2 q { X q; }
gate CX q, p { ctrl @ X q, p; }

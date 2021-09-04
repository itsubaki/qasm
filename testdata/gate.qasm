gate X q { U(pi, 0, pi) q; }

qubit[2] q;
reset q;

X q[0], q[1];

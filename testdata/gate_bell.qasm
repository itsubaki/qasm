gate H  q    { U(pi/2.0, 0, pi) q; }
gate X1 q    { pow(3) @ U(pi, 0, pi) q; }
gate X2 q    { inv    @ X1 q;}
gate X3 q    { inv    @ X2 q;}
gate X4 q    { inv    @ X3 q;}
gate CX q, p { ctrl   @ X4 q, p; }

qubit[2] q;
reset q;

H  q[0];
CX q[0], q[1];

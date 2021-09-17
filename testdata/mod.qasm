OPENQASM 3.0;
include "testdata/gate.qasm";

gate u1(a, b, c) q     { pow(3) @ U(a, b, c)  q; }
gate u2(a, b, c) q     { inv    @ u1(a, b, c) q; }
gate u3(a, b, c) q     { inv    @ u2(a, b, c) q; }
gate u4(a, b, c) q     { inv    @ u3(a, b, c) q; }
gate cu(a, b, c) q, p  { ctrl   @ u4(a, b, c) q, p; }

qubit[2] q;
reset q;

h q[0];
cu(pi, 0, pi) q[0], q[1];

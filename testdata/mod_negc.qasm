OPENQASM 3.0;
include "testdata/gate.qasm";

gate ncx  q0, q1     { negctrl @ x   q0, q1; }
gate cncx q0, q1, q2 { ctrl    @ ncx q0, q1, q2; }

qubit[3] q;
reset q;

x q[0];
cncx q[0], q[1], q[2];

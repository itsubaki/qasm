OPENQASM 3.0;
include "itsubaki/q.qasm";

qubit p;
qubit q[2];
reset p, q;

h p;

h  q[0];
cx q[0], q[1];

cx p,  q[0];
h  p;

cx q[0], q[1];
cz p, q[1];

measure p, q[0];

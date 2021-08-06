OPENQASM 3.0;
include "stdgates.qasm";

qubit    p;
qubit[2] q;
reset p, q;

x p;

h  q[0];
cx q[0], q[1];

cx p,  q[0];
h  p;

cx q[0], q[1];
cz p, q[1];

measure p, q[0];

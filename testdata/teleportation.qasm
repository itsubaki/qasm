OPENQASM 3.0;
include "stdgates.qasm";

qubit    phi;
qubit[2] q;
bit[2]   c;

reset phi, q;

x phi;

h  q[0]
cx q[0], q[1]
cx phi,  q[0]
h  phi

cx q[0], q[1];
cz phi,  q[1];

c[0] = measure phi;
c[1] = measure q[0];

OPENQASM 3.0;
include "itsubaki/q.qasm";

qubit q[2];
bit   c[2];

reset q;

h  q[0];
cx q[0], q[1];

measure q -> c;
c = measure q;
c[0] = measure q[0];
c[1] = measure q[1];

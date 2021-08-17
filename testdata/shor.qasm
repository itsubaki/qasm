OPENQASM 3.0;
include "itsubaki/q.qasm";

const N = 15;
const a = 7;

qubit r0[3];
qubit r1[4];

x r1[-1];
h r0;
cmodexp2(a, N) r0, r1;
iqft r0;

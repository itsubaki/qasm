OPENQASM 3.0;
include "itsubaki/q.qasm";

const N = 3 * 5;
const a = 7;

qubit[3] r0;
qubit[4] r1;
reset r0, r1;

x r1[-1];
h r0;
cmodexp2(a, N) r0, r1;
iqft r0;

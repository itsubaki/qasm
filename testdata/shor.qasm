OPENQASM 3.0;
include "testdata/gate.qasm";

const N = 15;
const a = 7;

qubit[3] r0;
qubit[4] r1;
reset r0, r1;

x r1[-1];
h r0;
CMODEXP2(a, N) r0, r1;
IQFT r0;

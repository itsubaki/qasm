OPENQASM 3.0;
include "testdata/stdgates.qasm";

const N = 3 * 5;
const a = 7;

qubit[3] r0;
qubit[4] r1;
bit[3] c;
reset r0, r1;

x r1[-1];
h r0;
CMODEXP2(a, N) r0, r1;
IQFT r0;

c = measure r0;

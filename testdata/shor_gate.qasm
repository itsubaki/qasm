OPENQASM 3.0;
include "itsubaki/q.qasm";

gate shor(a, N) r0, r1 {
    h r0;
    cmodexp2(a, N) r0, r1;
    iqft r0;
}

const N = 15;
const a = 7;

qubit[3] r0;
qubit[4] r1;
reset r0, r1;

x r1[-1];
shor(a, N) r0, r1;

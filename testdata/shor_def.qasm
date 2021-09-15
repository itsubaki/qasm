OPENQASM 3.0;
include "testdata/gate.qasm";

def shor(int[32] a, int[32] N) qubit[n] r0, qubit[m] r1 -> bit[n] {
    X r1[-1];
    H r0;
    CMODEXP2(a, N) r0, r1;
    IQFT r0;
    
    return measure r0;
}

const N = 15;
const a = 7;

qubit[3] r0;
qubit[4] r1;
bit[3]   c;
reset r0, r1;

c = shor(a, N) r0, r1;

def shor(int[32] a, int[32] N) qubit[n] r0, qubit[m] r1 -> bit[n] {
    h r0;
    cmodexp2(a, N) r0, r1;
    iqft r0;
    
    return measure r0;
}

const N = 15;
const a = 7;
bit[3]   c;

qubit[3] r0;
qubit[4] r1;
reset r0, r1;

x r1[-1];
c = shor(a, N) r0, r1;

def shor(a, N) qubit[n] r0, qubit[m] r1 -> bit[n] {
    x r1[-1];
    h r0;
    cmodexp2(a, N) r0, r1;
    iqft r0;
    return measure r0;
}

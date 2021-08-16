gate bell q0, q1 {
    h  q0;
    cx q0, q1;
}

gate shor(a, N) r0, r1 {
    x r1[-1];
    h r0;
    cmodexp2 a, N, r0, r1;
    iqft r0;
}

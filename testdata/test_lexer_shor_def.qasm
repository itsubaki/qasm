def shor(int[32] a, int[32] N) qubit[n] r0, qubit[m] r1 -> bit[n] { h r0; cmodexp2(a, N) r0, r1; iqft r0; return measure r0; }

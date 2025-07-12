OPENQASM 3.0;

def qft(qubit[3] q) {
    U(pi/2.0, 0, pi) q[0];
    ctrl @ U(0, 0, pi/2) q[0], q[1];
    ctrl @ U(0, 0, pi/4) q[0], q[2];

    U(pi/2.0, 0, pi) q[1];
    ctrl @ U(0, 0, pi/2) q[1], q[2];

    U(pi/2.0, 0, pi) q[2];

    // swap
    ctrl @ U(pi, 0, pi) q[0], q[2];
    ctrl @ U(pi, 0, pi) q[2], q[0];
    ctrl @ U(pi, 0, pi) q[0], q[2];
}

qubit[3] q;
U(pi, 0, pi) q[2];
qft(q);

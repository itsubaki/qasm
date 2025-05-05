OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate rz(theta) q { U(0, 0, theta) q; }
gate cx c, t { ctrl @ U(pi, 0, pi) c, t; }
gate crz(theta) c, t {
    ctrl @ rz(theta) c, t;
}

def qft(qubit[3] q) -> bool {
    h q[0];
    crz(pi/2) q[0], q[1];
    crz(pi/4) q[0], q[2];

    h q[1];
    crz(pi/2) q[1], q[2];

    h q[2];

    cx q[0], q[2];
    cx q[2], q[0];
    cx q[0], q[2];

    return true
}

qubit[3] q;
qft(q);

print;

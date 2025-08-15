OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate cccx c0, c1, c2, t { ctrl @ ctrl @ ctrl @ U(pi, 0, pi) c0, c1, c2, t; }

def oracle(qubit[4] q) {
    x q[2], q[3];
    h q[3];
    cccx q[0], q[1], q[2], q[3];
    h q[3];
    x q[2], q[3];
}

def diffuser(qubit[4] q) {
    h q;
    x q;
    h q[3];
    cccx q[0], q[1], q[2], q[3];
    h q[3];
    x q;
    h q;
}

qubit[4] q;
reset q;
h q;

int r = int(floor(pi/4 * sqrt(16.0)));
for int i in [0:r] {
    oracle(q);
    diffuser(q);
}

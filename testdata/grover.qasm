OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate cccz c0, c1, c2, t { ctrl(3) @ U(0, 0, pi) c0, c1, c2, t; }

// oracle for |110>|x>
def oracle(qubit[3] q, qubit a) {
    x q[2];

    x a;
    cccz q[0], q[1], q[2], a;
    x a;

    x q[2];
}

def diffuser(qubit[3] q, qubit a) {
    h q; h a;
    x q; x a;
    cccz q[0], q[1], q[2], a;
    x q; x a;
    h q; h a;
}

def G(qubit[3] q, qubit a) {
    oracle(q, a);
    diffuser(q, a);
}

const int n = 3;
qubit[n] q;
qubit a;
reset q;
reset a;

h q;
h a;

int N = 2**(n+1);
int R = int(pi/4 * sqrt(float(N)));

for int i in [1:R] {
    G(q, a);
}

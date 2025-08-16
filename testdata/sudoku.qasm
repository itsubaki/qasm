OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }

gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; }
gate cccz c0, c1, c2, t { ctrl(3) @ U(0, pi, 0) c0, c1, c2, t; }
gate ccccz c0, c1, c2, c3, t { ctrl(4) @ U(0, pi, 0) c0, c1, c2, c3, t; }
gate xor q0, q1, q2 { cx q0, q2; cx q1, q2; }

// The oracle constructs a Grover oracle that validates solutions to a 2x2 mini-sudoku puzzle.
// It enforces the following uniqueness constraints: a != b, c != d, a != c, b != d
// Valid solutions are [1, 0, 0, 1] and [0, 1, 1, 0].
def oracle(qubit[4] r, qubit[4] s, qubit a) {
    xor r[0], r[1], s[0];
    xor r[2], r[3], s[1];
    xor r[0], r[2], s[2];
    xor r[1], r[3], s[3];

    ccccz s[0], s[1], s[2], s[3], a;

    xor r[3], r[1], s[3];
    xor r[2], r[0], s[2];
    xor r[3], r[2], s[1];
    xor r[1], r[0], s[0];
}

def diffuser(qubit[4] r) {
    h r;
    x r;
    cccz r[0], r[1], r[2], r[3];
    x r;
    h r;
}

const int n = 4;
qubit[n] r;
qubit[4] s;
qubit a;

reset r;
reset s;
reset a;

h r;
h a;

int N = 2**n;
int M = 2;
int R = int(pi/4 * sqrt(float(N)/float(M)));

for int i in [0:R] {
    oracle(r, s, a);
    diffuser(r);
}

measure s;
measure a;

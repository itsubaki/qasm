OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate z q { U(0, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }

gate xor q0, q1, q2 {
    ctrl @ U(pi, 0, pi) q0, q2;
    ctrl @ U(pi, 0, pi) q1, q2;
}

// The oracle constructs a Grover oracle that checks solutions to a 2x2 sudoku puzzle.
// The oracle flips the phase when the following uniqueness constraints are satisfied: a != b, c != d, a != c, and b != d.
// The valid solutions are [1,0,0,1] and [0,1,1,0].
def oracle(qubit[4] r, qubit[4] s, qubit a) {
    xor r[0], r[1], s[0];
    xor r[2], r[3], s[1];
    xor r[0], r[2], s[2];
    xor r[1], r[3], s[3];

    ctrl(4) @ x s[0], s[1], s[2], s[3], a;

    xor r[1], r[3], s[3];
    xor r[0], r[2], s[2];
    xor r[2], r[3], s[1];
    xor r[0], r[1], s[0];
}

const int n = 4;
qubit[n] r;
qubit[4] s;
qubit a;

reset r;
reset s;
reset a;

h r;
x a;
h a;

// iteration 1
oracle(r, s, a);

h r;
x r;
ctrl(3) @ z r[0], r[1], r[2], r[3];
x r;
h r;

// iteration 2
oracle(r, s, a);

h r;
x r;
ctrl(3) @ z r[0], r[1], r[2], r[3];
x r;
h r;

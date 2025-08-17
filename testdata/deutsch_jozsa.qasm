OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; }

def constant(qubit q0, qubit q1) {
    x q1;
}

def balanced(qubit q0, qubit q1) {
    cx q0, q1;
}

qubit q0;
qubit q1;
reset q0;
reset q1;

x q1;
h q0;
h q1;

balanced(q0, q1);

h q0;
measure q0;


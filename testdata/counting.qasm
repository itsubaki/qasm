OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate crz(theta) c, t { ctrl @ U(0, 0, theta) c, t; }

gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; }
gate xor q0, q1, q2 { cx q0, q2; cx q1, q2; }
gate ccccz c0, c1, c2, c3, t { ctrl(3) @ U(0, 0, pi) c0, c1, c2, c3, t; }
gate cccccz c0, c1, c2, c3, c4, t { ctrl(5) @ U(0, 0, pi) c0, c1, c2, c3, c4, t; }

def oracle(qubit[4] r, qubit[4] s, qubit c, qubit a) {
    xor r[0], r[1], s[0];
    xor r[2], r[3], s[1];
    xor r[0], r[2], s[2];
    xor r[1], r[3], s[3];

    cccccz s[0], s[1], s[2], s[3], c, a;

    xor r[1], r[3], s[3];
    xor r[0], r[2], s[2];
    xor r[2], r[3], s[1];
    xor r[0], r[1], s[0];
}

def diffuser(qubit[4] r, qubit a) {
    h r;
    x r;
    ccccz r[0], r[1], r[2], r[3], a;
    x r;
    h r;
}

def swap(qubit[3] q) {
  cx q[0], q[2];
  cx q[2], q[0];
  cx q[0], q[2];
}

def inv_qft(qubit[3] q) {
  h q[2];

  crz(-pi/2) q[0], q[1];
  h q[1];

  crz(-pi/4) q[0], q[2];
  crz(-pi/2) q[1], q[2];
  h q[0];
}

qubit[3] c;
qubit[4] r;
qubit[4] s;
qubit a;

// initialize
reset c;
reset r;
reset s;
reset a;

h c;
h r;
x a;
h a;

// controlled-G**(2**0)
oracle(r, s, c[2], a);
diffuser(r, a);

// controlled-G**(2**1)
oracle(r, s, c[1], a);
diffuser(r, a);
oracle(r, s, c[1], a);
diffuser(r, a);

// controlled-G**(2**2)
oracle(r, s, c[0], a);
diffuser(r, a);
oracle(r, s, c[0], a);
diffuser(r, a);
oracle(r, s, c[0], a);
diffuser(r, a);
oracle(r, s, c[0], a);
diffuser(r, a);

// inverse qft
swap(c);
inv_qft(c);

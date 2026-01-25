OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate crz(theta) c, t { ctrl @ U(0, 0, theta) c, t; }

gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; }
gate xor q0, q1, q2 { cx q0, q2; cx q1, q2; }
gate c4z c0, c1, c2, c3, t { ctrl(4) @ U(0, 0, pi) c0, c1, c2, c3, t; }
gate c5z c0, c1, c2, c3, c4, t { ctrl(5) @ U(0, 0, pi) c0, c1, c2, c3, c4, t; }

def oracle(qubit[4] r, qubit[4] s, qubit c, qubit a) {
    xor r[0], r[1], s[0];
    xor r[2], r[3], s[1];
    xor r[0], r[2], s[2];
    xor r[1], r[3], s[3];

    x a;
    c5z s[0], s[1], s[2], s[3], c, a;
    x a;

    xor r[1], r[3], s[3];
    xor r[0], r[2], s[2];
    xor r[2], r[3], s[1];
    xor r[0], r[1], s[0];
}

def diffuser(qubit c, qubit[4] r) {
    h r;
    x r;
    c4z c, r[0], r[1], r[2], r[3];
    x r;
    h r;
}

def controlledG(qubit[4] r, qubit[4] s, qubit c, qubit a) {
  oracle(r, s, c, a);
  diffuser(c, r);
}

def swap(qubit[3] q) {
  cx q[0], q[2];
  cx q[2], q[0];
  cx q[0], q[2];
}

def inv_qft(qubit[3] q) {
  h q[2];
  crz(-pi/2) q[2], q[1];
  
  h q[1];
  crz(-pi/4) q[2], q[0];
  crz(-pi/2) q[1], q[0];
  
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

// controlled-G**(2**0)
controlledG(r, s, c[2], a);

// controlled-G**(2**1)
controlledG(r, s, c[1], a);
controlledG(r, s, c[1], a);

// controlled-G**(2**2)
controlledG(r, s, c[0], a);
controlledG(r, s, c[0], a);
controlledG(r, s, c[0], a);
controlledG(r, s, c[0], a);

// inverse qft
swap(c);
inv_qft(c);

// top 4
// [011 1001 0000 0][  3   9   0   0]( 0.0591 0.3592i): 0.1325
// [011 0110 0000 0][  3   6   0   0]( 0.0591 0.3592i): 0.1325
// [101 1001 0000 0][  5   9   0   0]( 0.0591-0.3592i): 0.1325
// [101 0110 0000 0][  5   6   0   0]( 0.0591-0.3592i): 0.1325

// measure c;
// M = N * sin(theta/2)**2, theta=2*pi*phi, phi=k/2**t
// 011(3) -> phi=0.3750, theta=0.7854, M=2.3431
// 101(5) -> phi=0.6250, theta=0.7854, M=2.3431

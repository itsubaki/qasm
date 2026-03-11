OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate cr(theta) c, t { ctrl @ U(0, 0, theta) c, t; }
gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; }
gate xor q0, q1, q2 { cx q0, q2; cx q1, q2; }
gate ccccz c0, c1, c2, c3, t { ctrl(4) @ U(0, 0, pi) c0, c1, c2, c3, t; }
gate cccccx c0, c1, c2, c3, c4, t { ctrl(5) @ U(pi, 0, pi) c0, c1, c2, c3, c4, t; }

def oracle(qubit[4] r, qubit[4] s, qubit c, qubit a) {
    xor r[0], r[1], s[0];
    xor r[2], r[3], s[1];
    xor r[0], r[2], s[2];
    xor r[1], r[3], s[3];

    cccccx s[0], s[1], s[2], s[3], c, a;

    xor r[1], r[3], s[3];
    xor r[0], r[2], s[2];
    xor r[2], r[3], s[1];
    xor r[0], r[1], s[0];
}

def diffuser(qubit c, qubit[4] r) {
    h r;
    x r;
    ccccz r[0], r[1], r[2], c, r[3];
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
  cr(-pi/2) q[2], q[1];
  
  h q[1];
  cr(-pi/4) q[2], q[0];
  cr(-pi/2) q[1], q[0];
  
  h q[0];
}

const int n = 3;
qubit[n] c;
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

// controlled-G
for int i in [0:n-1] {
  for int j in [0:(1<<i)-1] {
    controlledG(r, s, c[2-i], a);
  }
}

// inverse qft
swap(c);
inv_qft(c);

bit m = measure c;
// 011: phi=0.3750, theta=2.3562, M=2.3431
// 101: phi=0.6250, theta=3.9270, M=2.3431

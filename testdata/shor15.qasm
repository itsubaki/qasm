OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate h q { U(pi/2.0, 0, pi) q; }
gate cx c, t { ctrl @ U(pi, 0, pi) c, t; }
gate ccx c0, c1, t { ctrl(2) @ U(pi, 0, pi) c0, c1, t; }
gate crz(theta) c, t { ctrl @ U(0, 0, theta) c, t; }

def modexp(qubit[3] q, qubit[4] a) {
    // controlled-U^(2^0)
    cx q[2], a[1];
    cx q[2], a[2];

    // controlled-U^(2^1)
    cx        a[0], a[2];
    ccx q[1], a[2], a[0];
    cx        a[0], a[2];
    
    cx        a[3], a[1];
    ccx q[1], a[1], a[3];
    cx        a[3], a[1];
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

// N=15, a=7
qubit[3] q;
qubit[4] a;
reset q;
reset a;

x a[3];
h q;

modexp(q, a);
swap(q);
inv_qft(q);

measure a;
// measure q;
//
// 010 > 0.010 > 0.25 > 1/4; r=4.
// 110 > 0.110 > 0.75 > 3/4; r=4.
// gcd(pow(a, r/2)-1, N) = 3.
// gcd(pow(a, r/2)+1, N) = 5.

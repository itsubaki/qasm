OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; }

qubit psi;
U(1, 2, 3) psi;

// encode
qubit[2] enc;
cx psi, enc[0];
cx psi, enc[1];

// error (bit-flip)
x psi;

// add ancilla
qubit[2] a;

// error correction
cx psi,    a[0];
cx enc[0], a[0];
cx enc[0], a[1];
cx enc[1], a[1];

bit m0 = measure a[0];
bit m1 = measure a[1];

if(m0 && !m1) { x psi; }
if(m0 && m1)  { x enc[0]; }
if(!m0 && m1) { x enc[1]; }

// decode
cx psi, enc[1];
cx psi, enc[0];

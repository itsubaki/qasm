OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }
gate cx q0, q1 { ctrl @ U(pi, 0, pi) q0, q1; }

qubit phi;
U(1, 2, 3) phi;

// encoding
qubit[2] enc;
cx phi, enc[0];
cx phi, enc[1];

// error (bit-flip)
x phi;

// add ancilla
qubit[2] a;

// error correction
cx phi,    a[0];
cx enc[0], a[0];
cx enc[0], a[1];
cx enc[1], a[1];

bit m0 = measure a[0];
bit m1 = measure a[1];

if(m0 && !m1) { x phi; }
if(m0 && m1)  { x enc[0]; }
if(!m0 && m1) { x enc[1]; }

// decoding
cx phi, enc[1];
cx phi, enc[0];

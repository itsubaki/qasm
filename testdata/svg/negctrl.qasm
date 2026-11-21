OPENQASM 3.0;

qubit[3] q;

// Negative and positive controls.
negctrl @ x q[0], q[2];
ctrl @ x q[0], q[2];
negctrl(2) @ x q[0], q[1], q[2];
ctrl @ negctrl @ x q[0], q[1], q[2];
negctrl @ ctrl @ x q[0], q[1], q[2];

// A negative-controlled phase is shown as X P X.
negctrl @ gphase(pi/2) q[2];
negctrl @ ctrl @ gphase(pi/4) q[0], q[2];

measure q;

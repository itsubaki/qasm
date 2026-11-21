OPENQASM 3.0;

qubit[2] c;
qubit[2] t;
qubit s;

// Pair c[0] with t[0], and c[1] with t[1].
ctrl @ x c, t;

// Repeat the scalar target for both controls.
ctrl @ x c, t[0];

// Repeat the indexed scalar control for both targets.
negctrl @ x c[0], t;

// Broadcast arrays while repeating the scalar negative control.
ctrl @ negctrl @ x c, s, t;

// Broadcast controlled phases too.
ctrl @ gphase(pi/2) c;

measure c;
measure t;
measure s;

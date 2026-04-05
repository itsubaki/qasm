OPENQASM 3.0;

gate Rx(theta) q { U(theta, -pi/2, pi/2) q; }
gate W(theta) q { Rx(-2*theta) q; }
gate S(phi) q { U(0, 0, -2*phi) q; }

qubit q;
reset q;

const float theta = pi/6;
W(theta) q;
S(pi/4) q;
W(theta) q;
S(-pi/4) q;
W(theta) q;

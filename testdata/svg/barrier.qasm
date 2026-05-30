OPENQASM 3.0;

gate x q { U(pi, 0, pi) q; }

qubit[3] q;
qubit[3] a;

x q;
barrier q;
barrier a;
barrier q, a;
barrier q[0];
barrier q[0], q[1];
barrier q[0], a;

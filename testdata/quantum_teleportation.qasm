OPENQASM 3.0;

gate h q { U(pi/2.0, 0, pi) q; }
gate cx c, t { ctrl @ U(pi, 0, pi) c, t; }
gate cz c, t { ctrl @ U(0, pi, 0) c, t; }

qubit phi;
qubit a;
qubit t;

reset phi;
reset a;
reset t;

U(1, 2, 3) phi;

h a;
cx a, t;
cx phi, a;
h phi;

cx a, t;
cz phi, t;

measure phi;
measure a;

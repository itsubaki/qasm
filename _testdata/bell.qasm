OPENQASM 3.0;
include "_testdata/stdgates.qasm";

qubit q0;
qubit q1;

U(pi/2.0, 0, pi) q0;
ctrl @ U(pi, 0, pi) q0, q1;

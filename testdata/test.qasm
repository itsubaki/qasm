OPENQASM 3.0;
include "testdata/stdgates.qasm";

gate bell q, p {
    h q;
    cx q, p;
}

qubit[3] q;
reset    q;

x q[0];
ctrl @ bell q[0], q[1], q[2];

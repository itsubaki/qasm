OPENQASM 3.0;
include "itsubaki/q.qasm";

gate bell q0, q1 {
    h  q0;
    cx q0, q1;
}

qubit q;
qubit p;
reset q, p;

bell q, p;

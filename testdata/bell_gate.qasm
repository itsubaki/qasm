OPENQASM 3.0;
include "itsubaki/q.qasm";

gate bell q0, q1 {
    h  q0;
    cx q0, q1;
}

qubit[2] q;
bit[2]   c;

reset q;

bell q[0], q[1];
measure q -> c;

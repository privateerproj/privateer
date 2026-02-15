#!/bin/sh

set -x

STATUS=0

./privateer completion || STATUS=1
./privateer env || STATUS=1
./privateer generate-plugin -p ./test/data/CCC.VPC_2025.01.yaml -n example || STATUS=1
./privateer help || STATUS=1
./privateer list || STATUS=1
./privateer run -b ./test/data/ || STATUS=1
./privateer version || STATUS=1

exit $STATUS

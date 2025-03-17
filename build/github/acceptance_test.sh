#!/bin/sh

set -x

STATUS=0

./privateer completion
if [ $? -ne 0 ]; then
    STATUS=1
fi

./privateer generate-plugin -p ./test/data/CCC.VPC_2025.01.yaml -n example
if [ $? -ne 0 ]; then
    STATUS=1
fi

./privateer help
if [ $? -ne 0 ]; then
    STATUS=1
fi

./privateer list
if [ $? -ne 0 ]; then
    STATUS=1
fi

./privateer run -b ./test/data/
if [ $? -ne 0 ]; then
    STATUS=1
fi

./privateer version
if [ $? -ne 0 ]; then
    STATUS=1
fi

exit $STATUS



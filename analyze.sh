#!/bin/bash

out=all.csv #sans_unspecified.csv
flags=out_file=$out #,include_unspecified=false

rm -rf $out

go install ./protoc-gen-states

if [ -z $GOOGLEAPIS ]; then
  echo "Cloning googleapis to current directory"
  git clone --depth=1 https://github.com/googleapis/googleapis.git
  echo "export GOOGLEAPIS=$(pwd)/googleapis - DO THIS IN YOUR SHELL TO RETAIN"
  export GOOGLEAPIS=$(pwd)/googleapis
fi

pushd $GOOGLEAPIS

find google/cloud -name '*.proto' -type f -exec dirname {} \; | uniq | while read line; do
    protoc -I. --states_out=$flags:. $line/*.proto
done

popd

mv $GOOGLEAPIS/$out .

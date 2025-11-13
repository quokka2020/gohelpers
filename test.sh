#!/bin/bash


echo "Run vet and test"
for d in $(find . -type d ! \( -path \*.git -prune \) ! \( -path \*.vscode -prune \))
do
    (echo "vet and test '$d'" && cd $d && go vet && go test)
done

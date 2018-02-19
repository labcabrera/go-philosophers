#!/bin/bash

rm *.dat

./main -option=gen-rules -filename=rules-01.dat -count=1000
./main -option=gen-rules -filename=rules-02.dat -count=10000
./main -option=run -ruleFile1=rules-01.dat -ruleFile2=rules-02.dat -count=1
#./main option=process-rules filename = rules-01.dat
#!/bin/bash

RFOLDER=./random-rules

if [ ! -d $RFOLDER ];then
  mkdir random-rules
fi

./main -option=gen-rules -filename=$RFOLDER/rules-01.dat -count=50000
./main -option=gen-rules -filename=$RFOLDER/rules-02.dat -count=100000
./main -option=run -ruleFile1=$RFOLDER/rules-01.dat -ruleFile2=$RFOLDER/rules-02.dat -count=1

#./main option=process-rules filename = rules-01.dat

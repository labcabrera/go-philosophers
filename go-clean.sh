#!/bin/bash

if [ -f main ]; then
  rm main
fi

if [ -d random-rules ]; then
  rm -r random-rules
fi

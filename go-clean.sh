#!/bin/bash

if [ -f main ]; then
  rm main
fi

find . -name 'rules-*' -exec rm {} +


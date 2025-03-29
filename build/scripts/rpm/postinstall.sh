#!/bin/sh
BIN=/usr/local/bin/survey-repository
chmod u+x $BIN
chcon -t bin_t $BIN

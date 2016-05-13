#!/usr/local/bin/zsh
go run main.go | coffee -s -p > tmp.js; node tmp.js < tmp.dat

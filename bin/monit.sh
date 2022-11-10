#!/bin/bash

cd "/home/ubuntu/server"

str=$2
args=(${str//,/ })
length=${#args[@]}
program=${args[0]}
process=${args[0]}
if [ ${length} -eq 2 ]; then
    process=${program}"_"${args[1]}
fi

case $1 in
    "start")
        nohup ./${program} $3 > log/${process}_$(date +%Y_%m_%d_%H_%M_%s).log 2>&1 &
        echo $! > pid/${process}.pid
    ;;
    "stop")
        kill -15 `cat pid/${process}.pid`
        rm -f pid/${process}.pid
    ;;
    *)
        echo "Usage:"
        echo $0 "start name"
        echo $0 "stop name"
    ;;
esac

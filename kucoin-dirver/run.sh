#!/bin/bash

if [ $2 = "start" ]
then
    go build -o app
    nohup ./app $1 &> run_$1.log 2>&1&
    echo $! > run_$1.pid
    echo "run has been started succesfuly!"
fi

if [ $2 = "stop" ]
then
    kill -9 `cat run_$1.pid`
    rm run_$1.pid
    rm run_$1.log
    echo "run has been stopped"
fi
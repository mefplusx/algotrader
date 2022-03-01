#!/bin/bash

if [ $1 = "start" ]
then
    go build -o app
    nohup ./app -f &> find.log 2>&1&
    echo $! > find.pid
    echo "find has been started succesfuly!"
fi

if [ $1 = "stop" ]
then
    kill -9 `cat find.pid`
    rm find.pid
    rm find.log
    echo "find has been stopped"
fi
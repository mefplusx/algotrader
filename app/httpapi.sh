#!/bin/bash

if [ $1 = "start" ]
then
    go build -o app
    nohup ./app -ha  &> httpapi.log 2>&1&
    echo $! > httpapi.pid
    echo "httpapi has been started succesfuly!"
fi

if [ $1 = "stop" ]
then
    kill -9 `cat httpapi.pid`
    rm httpapi.pid
    rm httpapi.log
    echo "httpapi has been stopped"
fi
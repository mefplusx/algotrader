#!/bin/bash

if [ $1 = "start" ]
then
    go build -o app
    nohup ./app -b -ts &> bot.log 2>&1&
    echo $! > bot.pid
    echo "bot has been started succesfuly!"
fi

if [ $1 = "stop" ]
then
    kill -9 `cat bot.pid`
    rm bot.pid
    rm bot.log
    echo "bot has been stopped"
fi
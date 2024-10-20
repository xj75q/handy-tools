#!/bin/bash

while true
do
    now=$(date '+%H:%M:%S')
    #echo $now
    if [ $now = "13:30:00" ];then
       zenity --notification --text "your time cost"
    elif [ $now = "15:00:00" ];then
       zenity --notification --text "your time cost"
    elif [ $now = "17:00:00" ];then
       zenity --notification --text "your time cost"
    elif [ $now = "19:00:00" ];then
       zenity --notification --text "your time cost"
    fi
    sleep 1s
done


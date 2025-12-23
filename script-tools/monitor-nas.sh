#!/bin/bash

pwd=$(pwd)
monitor_min=15    #检查充电、电池模式间隔分钟数
battery_min=5     #掉电模式下检查电池间隔数
battery_low=20    #低电量百分比

notify_status="【切换到电池】NAS 切换到电池充电，检查设备供电..."
notify_battery="【低电量通知】NAS 电池电量低于20%"
notify_file="email"

checkfile(){
    cd $pwd
    if [ ! -f "$pwd/$notify_file" ]; then
        echo "请将两个文件放在同一目录下..."
        exit 1
    fi
}

battery(){
    while true
    do	
      sleep ${battery_min}m
      battery_level=$(acpi -b | grep -oP '[0-9]+(?=%)' | head -1)
      if [ "$battery_level" -lt $battery_low ]; then
         msg=$notify_battery"，当前电量为"${battery_level}"%"
         echo $notify_battery
         ./email send -s "$notify_battery" -d "$msg"
         exit
      fi
    done
}

monitor(){
    checkfile
    while true
    do	
      sleep ${monitor_min}m
      ac_state=$(acpi -a | awk '{print $3}')
      if [ "$ac_state" = "off-line" ]; then
            echo $notify_status 
            ./email send -s "$notify_status" -d "$notify_status"
            battery
            break
      fi
    done

}

#==================exeuce====================================

monitor


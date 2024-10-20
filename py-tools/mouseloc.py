#-*-coding=utf-8-*-

import pyautogui as pag
import time
import os


if __name__ == "__main__":
    try:
        while True:
            x,y=pag.position()
            print("当前鼠标坐标为：",x,y)
            time.sleep(0.2)
            os.system('clear')

    except KeyboardInterrupt:
        print(">> end...\n")

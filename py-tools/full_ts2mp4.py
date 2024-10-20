# -*- coding:utf-8 -*-
import requests
import re
import os
import shutil
import time
import argparse


class MergeFiles():
    def __init__(self, input_filename, output_path):
        self.filename = input_filename  # 这是ts文件的文件目录名或者文件名
        try:
            self.filePath = output_path + self.filename
        except Exception as e:
            print(e)
        self.ffmpeg = "F:\\ts2mp4\\ffmpeg-4.2.1\\bin\\ffmpeg.exe"

    def sort_key(self, str):
        re_digits = re.compile(r'(\d+)')
        pieces = re_digits.split(str)  # 切成数字与非数字
        pieces[1::2] = map(int, pieces[1::2])  # 将数字部分转成整数
        return pieces

    def sort_file(self):
        for root, dirs, files in os.walk(self.filePath):  # 遍历统计
            new_file_list = sorted(files, key=self.sort_key)
            return new_file_list

    def merge_file_to_ts(self):
        shell_str = "+".join(self.sort_file())
        command_str = 'copy /b ' + shell_str + " " + self.filePath + '.ts'
        os.chdir(self.filePath)
        print(command_str)
        result = os.system(command_str)
        if result == 0:
            print("合并ts文件成功")
            return self.filePath + ".ts"
        else:
            print("+++++", result)

    def switchToMp4(self, tsFile_path):
        output = str.replace(tsFile_path, ".ts", ".mp4")
        cmd_str = self.ffmpeg + " -i " + tsFile_path + " -c copy " + output
        result = os.system(cmd_str)
        if result == 0:
            print("1、已经将ts文件转换成对应的mp4文件")
            return 0

    def del_ts_files(self, full_ts_file):
        for root, dirs, files in os.walk(self.filePath):  # 遍历统计
            try:
                os.chdir("..")
                shutil.rmtree(root)
                print("2、该文件夹下所有ts目录删除")
            except Exception as e:
                print("2、删除分散ts文件时出错:", e)
        if os.path.exists(full_ts_file):
            os.remove(full_ts_file)
            print("3、已经将外部拼接成的ts文件删除")
        else:
            print("要删除的full_ts不存在")

    def run(self):
        # ts_file =self.merge_file_to_ts()
        # print(ts_file)
        ts_file = self.filePath + ".ts"
        result = self.switchToMp4(ts_file)
        if result == 0:
            time.sleep(3)
            self.del_ts_files(ts_file)
            print("4、<SUCCESS> mp4文件已经成功生成,同时已经将多余文件删除")


def parse_parm():
    parser = argparse.ArgumentParser()
    parser.add_argument('-i', '--input', help='input filename', required=True,
                        type=str)  # , default=False ,,action='store_true'
    parser.add_argument('-o', '--output', help='output to save path', required=True,
                        type=str)  # default=False ,,action='store_true'
    args = parser.parse_args()
    input_str = args.input
    output_str = args.output
    return input_str, output_str


if __name__ == "__main__":
    input_filename, output_path = parse_parm()
    mf = MergeFiles(input_filename, output_path)
    mf.run()

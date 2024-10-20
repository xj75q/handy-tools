
# -*- coding: UTF-8 -*-

# 功能 : 将输入路径下的所有webp格式转为png or jpg

import argparse
import os
from PIL import Image

class ConvertWebp(object):
    def __init__(self,input,output,img_type):
        self.CURRENT_PATH= input
        self.output = output if output is not None else self.CURRENT_PATH
        self.IMG_EXP = img_type if img_type is not None else ".jpg"

    def imgList(self,filepath):
        image_list =[]
        for path, file_dir, files in os.walk(filepath):
            for file_name in files:
                name, e = os.path.splitext(file_name)
                if e.lower() == ".webp":
                    img_name=os.path.join(path, name)
                    image_list.append(img_name)
        return image_list

    def convertImage(self):
        for webpPath in self.imgList(self.CURRENT_PATH):
            img = Image.open(webpPath + ".webp")
            img.load()
            img_name = webpPath.split(os.path.sep)[-1]
            img.save("{}{}{}{}".format(self.output,os.path.sep,img_name,self.IMG_EXP))

def parse_parm():
    parser = argparse.ArgumentParser()
    parser.add_argument('-i', '--input', help='input img dir', required=True,
                        type=str)  # , default=False ,,action='store_true'
    parser.add_argument('-o', '--output', help='output to save dir', required=False,
                        type=str)  # default=False ,,action='store_true'
    parser.add_argument('-t', '--type', help='save the file type', required=False,
                        type=str)  # default=False ,,action='store_true'
    args = parser.parse_args()
    input_str = args.input
    output_str = args.output
    type_str = args.type
    return input_str, output_str, type_str


if __name__ == "__main__":
    input_path, output_path,type_str = parse_parm()
    a = ConvertWebp(input_path,output_path,type_str)
    a.convertImage()

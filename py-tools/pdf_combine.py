import os
import argparse
from PyPDF2 import PdfFileMerger


class PdfCombine():
    def __init__(self, input_path, output_path):
        self.input = input_path
        self.output = output_path
        self.new_name = "combine.pdf"

    def pdfHandler(self):
        # target_path = 'G:\\小说pdf_转换后'
        pdf_lst = [f for f in os.listdir(self.input) if f.endswith('.pdf')]
        pdf_lst = [os.path.join(self.input, filename) for filename in pdf_lst]

        file_merger = PdfFileMerger()
        for pdf in pdf_lst:
            file_merger.append(pdf)  # 合并pdf文件
        try:
            if self.output is None:
                writefile = self.input + self.new_name
                file_merger.write(writefile)
            else:
                writefile = self.output + self.new_name
                file_merger.write(writefile)

        except Exception as e:
            print("there is a err: ", e)


def parse_parm():
    parser = argparse.ArgumentParser()
    parser.add_argument('-i', '--input', help='input pdf dir', required=True,
                        type=str)  # , default=False ,,action='store_true'
    parser.add_argument('-o', '--output', help='output to save dir', required=False,
                        type=str)  # default=False ,,action='store_true'
    args = parser.parse_args()
    input_str = args.input
    output_str = args.output
    return input_str, output_str


if __name__ == "__main__":
    input_path, output_path = parse_parm()
    pdfCtr = PdfCombine(input_path, output_path)
    pdfCtr.pdfHandler()

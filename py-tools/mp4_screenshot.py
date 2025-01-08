import os
import shutil
import time
import re
import argparse
from multiprocessing.pool import Pool
class Mp4ToScreenshot():
    def __init__(self,input_dir,output_dir,per_time):
        #ffmpeg程序所在目录
        self.ffmpeg_path = r"F:\bin\ffmpeg.exe"
        #MP4所在目录
        self.mp4_dir = input_dir
        #文件名
        self.dir_name = self.mp4_dir.split("\\")[-1]
        #输出文件所在路径
        self.save_dir = output_dir+"\\"+self.dir_name
        if per_time ==None:
            self.time_interval = 5
        else:
            self.time_interval = per_time
        self.single_dir = True

    def get_sfname(self,file):
        if ".mp4" in file:
            file_name = file.replace(".mp4", "")
            screenshot_file_name = file_name + "_截图"
            return screenshot_file_name
        elif ".avi" in file:
            file_name = file.replace(".avi", "")
            screenshot_file_name = file_name + "_截图"
            return screenshot_file_name
        elif ".wmv" in file:
            file_name = file.replace(".wmv", "")
            screenshot_file_name = file_name + "_截图"
            return screenshot_file_name
        elif ".flv" in file:
            file_name = file.replace(".flv", "")
            screenshot_file_name = file_name + "_截图"
            return screenshot_file_name
        elif ".mkv" in file:
            file_name = file.replace(".mkv", "")
            screenshot_file_name = file_name + "_截图"
            return screenshot_file_name
        elif ".f4v" in file:
            file_name = file.replace(".f4v", "")
            screenshot_file_name = file_name + "_截图"
            return screenshot_file_name
        elif ".rmvb" in file:
            file_name = file.replace(".rmvb", "")
            screenshot_file_name = file_name + "_截图"
            return screenshot_file_name

    def video_type(self,file):
        type_list = [".mp4",".avi",".wmv",".flv",".mkv",".f4v",".rmvb"]
        for f_type in type_list:
            if f_type in file:
                return True

    def exec_ffmpeg(self,save_dir,work_dir,screenshot_file_name,count):
        start_time = time.time()
        # 参考："ffmpeg -i 少年派-读书.mp4 -vf fps=fps=1/5 -f image2 D:\000\out%d.jpg"
        save_path = '"{0}\\{1}.jpg"'.format(save_dir, "%d".zfill(2))  # +"\\"+"%d.jpg".zfill(3)
        # 这里的fps当中的5代表每隔5秒
        cmd_param = "{0} -i {1} -vf fps=fps=1/{2} -f image2 {3} -loglevel quiet"
        cmd_str = cmd_param.format(self.ffmpeg_path, work_dir, self.time_interval, save_path)
        # cmd_str = cmd_param.format(ffmpeg,work_dir,self.time_interval,save_path)
        output_log_1 = ">> === ({0}) 的截图正在进行中，请稍后...".format(screenshot_file_name)
        print(output_log_1)
        result = os.system(cmd_str)
        if result == 0:
            count += 1
            output_log = "## {0}> {1}的截图已经保存成功在({2})路径下 ^_^".format(count, screenshot_file_name, save_dir)
            print(output_log)

        else:
            print("??截屏出错，请打开log查看。。。")
        end_time = time.time()
        time_cost = (end_time-start_time)/60
        print("## 此次视频转换操作耗时%.2f min \n"%time_cost)
       


    def copy_non_video_file(self,next_path,root):
        start_time = time.time()
        if len(next_path) != 1:
            save_dir = self.save_dir + "\\" + next_path[-1]
            time.sleep(1)
            shutil.copytree(root, save_dir,
                            ignore=shutil.ignore_patterns('*.mp4', '*.avi', '*.wmv', '*.flv', '*.swf', '*.SWF',
                                                          '*.rmvb', '*.baiduyun.downloading'))
        else:
            time.sleep(1)
            shutil.copytree(root, self.save_dir,
                            ignore=shutil.ignore_patterns('*.mp4', '*.avi', '*.wmv', '*.flv', '*.swf', '*.SWF',
                                                          '*.rmvb', '*.baiduyun.downloading'))
        end_time = time.time()
        cost_time = (end_time-start_time)/60
        print("===复制其他文件耗时%.2f min,已经对应文件复制到相应的目标路径下==="%cost_time)

    def run(self):
        start_time = time.time()
        count=0
        index = 0
        video_pool = Pool(5)
        non_video_pool = Pool(5)
        for root,dirs,files in os.walk(self.mp4_dir):
            index+=1
            for file in files:
                work_dir = '"{0}\\{1}"'.format(root, file)
                v_type =self.video_type(file)
                if v_type:
                    screenshot_file_name = self.get_sfname(file)
                    next_path =root.split("\\")[-1]
                    if screenshot_file_name !=None:
                        if next_path==self.dir_name:
                            save_dir = self.save_dir+"\\"+screenshot_file_name
                        else:
                            next_full_path = root.split(self.dir_name+"\\")[-1]
                            save_dir = self.save_dir+"\\"+next_full_path+"\\"+screenshot_file_name

                        if not os.path.exists(save_dir):  # 路径不存在时创建一个
                            os.makedirs(save_dir)
                        time.sleep(5)
                        video_pool.apply_async(self.exec_ffmpeg(save_dir,work_dir,screenshot_file_name,count))
                else:
                    if self.dir_name in root:
                        index+=1
                        next_path = root.split(self.dir_name + "\\")
                        non_video_pool.apply_async(self.copy_non_video_file(next_path,root))

        video_pool.close()
        video_pool.join()
        non_video_pool.close()
        non_video_pool.join()
        end_time=time.time()
        time_cost = (end_time-start_time)/60
        print("########所有任务执行完成,总耗时%.2f min#############\n"%time_cost)



def parse_parm():
    parser = argparse.ArgumentParser()
    parser.add_argument('-i','--input',help='input mp4 dir',required=True,type=str)#, default=False ,,action='store_true'
    parser.add_argument('-o','--output',help='output to the save path',required=True,type=str) #default=False ,,action='store_true'
    parser.add_argument('-p','--per',help='per time',required=False,type=str) #default=False ,,action='store_true'
    args = parser.parse_args()
    input_str = args.input
    output_str = args.output
    per_time = args.per
    return input_str,output_str,per_time

if __name__ =="__main__":
    input_dir,output_dir,per_time =parse_parm()
    mmt = Mp4ToScreenshot(input_dir,output_dir,per_time)
    mmt.run()


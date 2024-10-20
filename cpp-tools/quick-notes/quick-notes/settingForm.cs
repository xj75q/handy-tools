using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace speed_notes
{
    public partial class settingForm : Form
    {
        public settingForm()
        {
            InitializeComponent();
        }

        private void settingForm_Load(object sender, EventArgs e)
        {
            this.ShowIcon = false;
            this.ShowInTaskbar = false;
            this.TopMost = true;

            // 添加选项列表
            fileType.Items.Add("true");
            fileType.Items.Add("false");
            fileType.DropDownStyle = ComboBoxStyle.DropDownList;
            // 初始选择第一个选
            fileType.SelectedIndex = 0;
        
        }

    
        private void textPath_TextChanged(object sender, EventArgs e)
        {
            this.textPath.WordWrap = false;
            this.textPath.Multiline = false;
            this.textPath.ScrollBars = ScrollBars.Horizontal;
        }


        private void textPath_KeyPress(object sender, KeyPressEventArgs e)
        {
            if (e.KeyChar != (char)Keys.Back)
            {
                e.Handled = true;
            }
        }

        private void selectPath_Click(object sender, EventArgs e)
        {
            System.Windows.Forms.FolderBrowserDialog dialog = new System.Windows.Forms.FolderBrowserDialog();
            dialog.Description = "请选择文件夹";
            dialog.RootFolder = Environment.SpecialFolder.MyComputer;
            dialog.ShowNewFolderButton = true;
            if (textPath.Text.Length > 0)
                dialog.SelectedPath = textPath.Text;
            if (dialog.ShowDialog() == DialogResult.OK)
            {
                this.textPath.Text = dialog.SelectedPath;
            }


        }

        private  void settingOK_Click(object sender, EventArgs e)
        {
            MessageBoxButtons okCancel= MessageBoxButtons.OKCancel;
            string typeStr;
            StringBuilder str = new StringBuilder();
            string name = this.fileName.Text;
            string path = this.textPath.Text;
            string ftype = this.fileType.Text;
            if (path == "")
            {
                MessageBox.Show("文件路径不能为空");
                return;
            }
            if (name == "")
            {
                MessageBox.Show("文件名不能为空");
                return;
            }

            if (ftype == "true")
            {
                typeStr = "设置为txt";

            }
            else
            {
                typeStr = "不设任何后缀";
            }

            string filepath = string.Format("{0}{1}{2}", "文件路径为：", path, "\n"); 
            string filename = string.Format("{0}{1}{2}", "文件名为   ：", name, "\n");
            string filetype = string.Format("{0}{1}{2}", "文件后缀为：", typeStr, "\n");
            str.Append(filepath);
            str.Append(filename);
            str.Append(filetype);
            
            DialogResult dr = MessageBox.Show(str.ToString(), "提示", okCancel);

            if (dr == DialogResult.OK)

            {
                 this.cmdSetting(path,name, this.fileType.Text);
                this.Close();
            }
            

        }

        private void settingCancel_Click(object sender, EventArgs e)
        {
            this.Close();
        }

        //todo 以后增加messagebox返回输出
        private  void cmdSetting(string path, string fname, string ftype) {
            MainForm item = new MainForm();
            string cmdstr = String.Format("{0}{1}{2}{3}{4}{5}","record config -p ",path," -n ",fname," -t=",ftype);
             item.WinCmd(cmdstr);
        
        }

    }
}

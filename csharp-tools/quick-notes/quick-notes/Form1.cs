using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;
using System.Diagnostics;



namespace speed_notes
{
    public partial class MainForm : Form
    {
        public MainForm()
        {
            InitializeComponent();
            int ScreenWidth = System.Windows.Forms.Screen.PrimaryScreen.WorkingArea.Width;
           // int ScreenHeight = System.Windows.Forms.Screen.PrimaryScreen.WorkingArea.Height;
            this.Location = new System.Drawing.Point(ScreenWidth - 500, 20);
            this.TopMost = true;
        }

        private void MainForm_Load(object sender, EventArgs e)
        {
            this.ShowIcon = false;
            this.ShowInTaskbar = false;
            textBox.AutoSize = false;
            this.TitleClick += new EventHandler(MainForm_TitleClick);
            this.textBox.Focus();
            

        }

        private void MainForm_TitleClick(object sender, EventArgs e)
        {
            settingForm setting = new settingForm();
            setting.Show();

        }

        protected override void WndProc(ref Message m)
        {
            const int WM_NCLBUTTONDBLCLK = 0xA3;

            if (m.Msg == WM_NCLBUTTONDBLCLK)
            {

                TitleClick(this, null);
            }
            base.WndProc(ref m);

        }

        public event EventHandler TitleClick;


        private void textBox_TextChanged(object sender, EventArgs e)
        {
            this.textBox.WordWrap = false;
            this.textBox.Multiline = true;
            this.textBox.ScrollBars = ScrollBars.Vertical;

        }

        private  void textBox_KeyPress(object sender, KeyEventArgs e)
        {

            if ((e.Modifiers == Keys.Shift && e.KeyCode == Keys.Enter) || (e.Modifiers == Keys.Control & e.KeyCode == Keys.Enter))
            {

                return;
            }


            else if (e.KeyCode == Keys.Enter)
            {
                var str = textBox.Text;
                if (str.EndsWith(Environment.NewLine))
                {
                     textBox.Text = str.Substring(0, str.Length - Environment.NewLine.Length);
                     this.buttonOK_Click(sender, e);
                }
                
            }

        }

        private  void buttonOK_Click(object sender, EventArgs e)

        {
            string inputText = textBox.Text;
            string sendText;
            if (inputText.Contains("\""))
            {
                sendText = inputText.Replace("\"", "\"\"\"");
            }
            else
            {
                sendText = inputText;
            }

            if (inputText.Contains(System.Environment.NewLine))
            {
                sendText = sendText.Replace(System.Environment.NewLine, "·^");
            }
            else
            {
                sendText = inputText;
            }
            //todo 读取配置文件，判断有无当前日期，如果没有则写入
            //string now = DateTime.Now.ToString("yyyy-MM-dd");
           

            string cmdstr = String.Format("{0}{1}{2}", "record note -c ^\"", sendText, "\"");    
            this.WinCmd(cmdstr);
            //this.textBox.Focus();  
           // this.textBox.Select(0, 0);
          //  this.textBox.ScrollToCaret();
            this.textBox.Clear();
            this.hideWindow();
        }

       

        private void hideWindow()
        {
            int sideThickness = 6;
            //隐藏到屏幕左边缘
            if (this.Left == 0)
            {
                this.Left = sideThickness - this.Width;
            }
            //隐藏到屏幕右边缘
            else if (this.Left == Screen.PrimaryScreen.WorkingArea.Width - this.Width)
            {
                this.Left = Screen.PrimaryScreen.WorkingArea.Width - sideThickness;
            }
            //隐藏到屏幕上边缘
            else if (this.Top == 0 && this.Left > 0 && this.Left < Screen.PrimaryScreen.WorkingArea.Width - this.Width)
            {
                this.Top = sideThickness - this.Height;
            }
            //隐藏到屏幕上边缘
            else
            {
                this.Top = sideThickness - this.Height;
            }

        }


        private void buttonCancel_Click(object sender, EventArgs e)
        {
            // 将焦点设置到文本框
            
            this.textBox.Focus();
            this.hideWindow();
        }



        private void timer1_Tick(object sender, EventArgs e)
        {
            timer1.Interval = 200;
            this.AutoSideHideOrShow();
        }

        private void AutoSideHideOrShow()
        {
            int sideThickness = 4;//边缘的厚度，窗体停靠在边缘隐藏后留出来的可见部分的厚度

            //如果窗体最小化或最大化了则什么也不
            if (this.WindowState == FormWindowState.Minimized || this.WindowState == FormWindowState.Maximized)
            {
                return;
            }




            //如果鼠标在窗体内
            if (Cursor.Position.X >= this.Left && Cursor.Position.X < this.Right && Cursor.Position.Y >= this.Top && Cursor.Position.Y < this.Bottom)
            {
                //如果窗体离屏幕边缘很近，则自动停靠在该边缘
                if (this.Top <= sideThickness)
                {
                    this.Top = 0;
                }
                if (this.Left <= sideThickness)
                {
                    this.Left = 0;
                }
                if (this.Left >= Screen.PrimaryScreen.WorkingArea.Width - this.Width - sideThickness)
                {
                    this.Left = Screen.PrimaryScreen.WorkingArea.Width - this.Width;
                }
            }
            //当鼠标离开窗体以后
            else
            {
                //隐藏到屏幕左边缘
                if (this.Left == 0)
                {
                    this.Left = sideThickness - this.Width;
                }
                //隐藏到屏幕右边缘
                else if (this.Left == Screen.PrimaryScreen.WorkingArea.Width - this.Width)
                {
                    this.Left = Screen.PrimaryScreen.WorkingArea.Width - sideThickness;
                }
                //隐藏到屏幕上边缘
                else if (this.Top == 0 && this.Left > 0 && this.Left < Screen.PrimaryScreen.WorkingArea.Width - this.Width)
                {
                    this.Top = sideThickness - this.Height;
                }
            }
        }

        
        
        public void WinCmd(string cmdline)
        {
           
            using (var process = new Process())
            {
                process.StartInfo.FileName = "cmd.exe";
                process.StartInfo.UseShellExecute = false;
                process.StartInfo.RedirectStandardInput = true;
                process.StartInfo.RedirectStandardOutput = true;
                process.StartInfo.RedirectStandardError = true;
                process.StartInfo.CreateNoWindow = true;

                process.Start();
                process.StandardInput.AutoFlush = true;
                Console.WriteLine(cmdline);
                process.StandardInput.WriteLine(cmdline);
                process.StandardInput.WriteLine("exit");
                process.Close();
            }
            //string output = process.StandardOutput.ReadToEnd();
            // process.WaitForExit();

        }
       


        private void noteNotify_MouseDoubleClick(object sender, MouseEventArgs e)
        {
            if (this.Visible)
            {
                this.WindowState = FormWindowState.Normal;
                this.Visible = true;
                this.Top = 20;
                
            }
           
        }

        private void hideItem_Click(object sender, EventArgs e)
        {
            this.noteNotify.Visible = true;
            this.hideWindow();
        }

        
        private void originItem_Click(object sender, EventArgs e)
        {
            
            this.WindowState = FormWindowState.Normal;
            this.noteNotify.Visible = true;
            this.Top = 20;
        }
        

        private void closeItem_Click(object sender, EventArgs e)
        {
            if (MessageBox.Show("你确定要退出吗？", "系统提示", MessageBoxButtons.YesNo, MessageBoxIcon.Information, MessageBoxDefaultButton.Button1) == DialogResult.Yes)
            {

                this.noteNotify.Visible = false;
                this.Close();
                this.Dispose();
                System.Environment.Exit(System.Environment.ExitCode);

            }
        }

           /*
            public async Task<string> WinCmd(string cmdline)
            {
                return await Task.Run(() =>
                {
                    using (var process = new Process())
                    {
                        process.StartInfo.FileName = "cmd.exe";
                        process.StartInfo.UseShellExecute = false;
                        process.StartInfo.RedirectStandardInput = true;
                        process.StartInfo.RedirectStandardOutput = true;
                        process.StartInfo.RedirectStandardError = true;
                        process.StartInfo.CreateNoWindow = true;

                        process.Start();
                        process.StandardInput.AutoFlush = true;
                        process.StandardInput.WriteLine(cmdline + " &exit");

                        string output = process.StandardOutput.ReadToEnd();

                        process.WaitForExit();
                        process.Close();

                        return output;
                    }
                });

            }
            */
            
            
        }
    
}

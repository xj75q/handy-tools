namespace speed_notes
{
    partial class MainForm
    {
        /// <summary>
        /// 必需的设计器变量。
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// 清理所有正在使用的资源。
        /// </summary>
        /// <param name="disposing">如果应释放托管资源，为 true；否则为 false。</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows 窗体设计器生成的代码

        /// <summary>
        /// 设计器支持所需的方法 - 不要修改
        /// 使用代码编辑器修改此方法的内容。
        /// </summary>
        private void InitializeComponent()
        {
            this.components = new System.ComponentModel.Container();
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(MainForm));
            this.label1 = new System.Windows.Forms.Label();
            this.buttonCancel = new System.Windows.Forms.Button();
            this.buttonOK = new System.Windows.Forms.Button();
            this.timer1 = new System.Windows.Forms.Timer(this.components);
            this.textBox = new System.Windows.Forms.TextBox();
            this.fileSystemWatcher1 = new System.IO.FileSystemWatcher();
            this.noteNotify = new System.Windows.Forms.NotifyIcon(this.components);
            this.NoteMenuStrip = new System.Windows.Forms.ContextMenuStrip(this.components);
            this.hideItem = new System.Windows.Forms.ToolStripMenuItem();
            this.originItem = new System.Windows.Forms.ToolStripMenuItem();
            this.closeItem = new System.Windows.Forms.ToolStripMenuItem();
            ((System.ComponentModel.ISupportInitialize)(this.fileSystemWatcher1)).BeginInit();
            this.NoteMenuStrip.SuspendLayout();
            this.SuspendLayout();
            // 
            // label1
            // 
            this.label1.AutoSize = true;
            this.label1.Location = new System.Drawing.Point(34, 52);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(77, 12);
            this.label1.TabIndex = 0;
            this.label1.Text = "请输入文字：";
            // 
            // buttonCancel
            // 
            this.buttonCancel.Location = new System.Drawing.Point(157, 279);
            this.buttonCancel.Name = "buttonCancel";
            this.buttonCancel.Size = new System.Drawing.Size(75, 23);
            this.buttonCancel.TabIndex = 2;
            this.buttonCancel.Text = "取消";
            this.buttonCancel.UseVisualStyleBackColor = true;
            this.buttonCancel.Click += new System.EventHandler(this.buttonCancel_Click);
            // 
            // buttonOK
            // 
            this.buttonOK.Location = new System.Drawing.Point(254, 279);
            this.buttonOK.Name = "buttonOK";
            this.buttonOK.Size = new System.Drawing.Size(75, 23);
            this.buttonOK.TabIndex = 3;
            this.buttonOK.Text = "确认";
            this.buttonOK.UseVisualStyleBackColor = true;
            this.buttonOK.Click += new System.EventHandler(this.buttonOK_Click);
            // 
            // timer1
            // 
            this.timer1.Enabled = true;
            this.timer1.Interval = 200;
            this.timer1.Tick += new System.EventHandler(this.timer1_Tick);
            // 
            // textBox
            // 
            this.textBox.Location = new System.Drawing.Point(36, 77);
            this.textBox.Multiline = true;
            this.textBox.Name = "textBox";
            this.textBox.ScrollBars = System.Windows.Forms.ScrollBars.Both;
            this.textBox.Size = new System.Drawing.Size(293, 180);
            this.textBox.TabIndex = 4;
            this.textBox.WordWrap = false;
            this.textBox.TextChanged += new System.EventHandler(this.textBox_TextChanged);
            this.textBox.KeyDown += new System.Windows.Forms.KeyEventHandler(this.textBox_KeyPress);
            // 
            // fileSystemWatcher1
            // 
            this.fileSystemWatcher1.EnableRaisingEvents = true;
            this.fileSystemWatcher1.SynchronizingObject = this;
            // 
            // noteNotify
            // 
            this.noteNotify.ContextMenuStrip = this.NoteMenuStrip;
            this.noteNotify.Icon = ((System.Drawing.Icon)(resources.GetObject("noteNotify.Icon")));
            this.noteNotify.Text = "speed-note";
            this.noteNotify.Visible = true;
            this.noteNotify.MouseDoubleClick += new System.Windows.Forms.MouseEventHandler(this.noteNotify_MouseDoubleClick);
            // 
            // NoteMenuStrip
            // 
            this.NoteMenuStrip.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.hideItem,
            this.originItem,
            this.closeItem});
            this.NoteMenuStrip.Name = "NoteMenuStrip";
            this.NoteMenuStrip.Size = new System.Drawing.Size(101, 70);
            // 
            // hideItem
            // 
            this.hideItem.Name = "hideItem";
            this.hideItem.Size = new System.Drawing.Size(100, 22);
            this.hideItem.Tag = "";
            this.hideItem.Text = "隐藏";
            this.hideItem.Click += new System.EventHandler(this.hideItem_Click);
            // 
            // originItem
            // 
            this.originItem.Name = "originItem";
            this.originItem.Size = new System.Drawing.Size(100, 22);
            this.originItem.Text = "还原";
            this.originItem.Click += new System.EventHandler(this.originItem_Click);
            // 
            // closeItem
            // 
            this.closeItem.Name = "closeItem";
            this.closeItem.Size = new System.Drawing.Size(100, 22);
            this.closeItem.Text = "退出";
            this.closeItem.Click += new System.EventHandler(this.closeItem_Click);
            // 
            // MainForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 12F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(375, 324);
            this.Controls.Add(this.textBox);
            this.Controls.Add(this.buttonOK);
            this.Controls.Add(this.buttonCancel);
            this.Controls.Add(this.label1);
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.FixedDialog;
            this.Icon = ((System.Drawing.Icon)(resources.GetObject("$this.Icon")));
            this.MaximizeBox = false;
            this.MinimizeBox = false;
            this.Name = "MainForm";
            this.StartPosition = System.Windows.Forms.FormStartPosition.Manual;
            this.Text = "speed-note";
            this.Load += new System.EventHandler(this.MainForm_Load);
            ((System.ComponentModel.ISupportInitialize)(this.fileSystemWatcher1)).EndInit();
            this.NoteMenuStrip.ResumeLayout(false);
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.Label label1;
        private System.Windows.Forms.Button buttonCancel;
        private System.Windows.Forms.Button buttonOK;
        private System.Windows.Forms.Timer timer1;
        private System.Windows.Forms.TextBox textBox;
        private System.IO.FileSystemWatcher fileSystemWatcher1;
        private System.Windows.Forms.NotifyIcon noteNotify;
        private System.Windows.Forms.ContextMenuStrip NoteMenuStrip;
        private System.Windows.Forms.ToolStripMenuItem hideItem;
        private System.Windows.Forms.ToolStripMenuItem originItem;
        private System.Windows.Forms.ToolStripMenuItem closeItem;
    }
}


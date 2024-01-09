namespace speed_notes
{
    partial class settingForm
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.fileType = new System.Windows.Forms.ComboBox();
            this.label1 = new System.Windows.Forms.Label();
            this.label2 = new System.Windows.Forms.Label();
            this.label3 = new System.Windows.Forms.Label();
            this.label4 = new System.Windows.Forms.Label();
            this.fileName = new System.Windows.Forms.TextBox();
            this.settingCancel = new System.Windows.Forms.Button();
            this.settingOK = new System.Windows.Forms.Button();
            this.folderBrowserDialog = new System.Windows.Forms.FolderBrowserDialog();
            this.textPath = new System.Windows.Forms.TextBox();
            this.selectPath = new System.Windows.Forms.Button();
            this.SuspendLayout();
            // 
            // fileType
            // 
            this.fileType.Location = new System.Drawing.Point(129, 145);
            this.fileType.Name = "fileType";
            this.fileType.Size = new System.Drawing.Size(159, 20);
            this.fileType.TabIndex = 15;
            // 
            // label1
            // 
            this.label1.AutoSize = true;
            this.label1.Location = new System.Drawing.Point(52, 69);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(53, 12);
            this.label1.TabIndex = 5;
            this.label1.Text = "文件路径";
            // 
            // label2
            // 
            this.label2.AutoSize = true;
            this.label2.Location = new System.Drawing.Point(52, 111);
            this.label2.Name = "label2";
            this.label2.Size = new System.Drawing.Size(41, 12);
            this.label2.TabIndex = 6;
            this.label2.Text = "文件名";
            // 
            // label3
            // 
            this.label3.AutoSize = true;
            this.label3.Location = new System.Drawing.Point(52, 153);
            this.label3.Name = "label3";
            this.label3.Size = new System.Drawing.Size(53, 12);
            this.label3.TabIndex = 7;
            this.label3.Text = "文件后缀";
            // 
            // label4
            // 
            this.label4.AutoSize = true;
            this.label4.Font = new System.Drawing.Font("宋体", 8F);
            this.label4.ForeColor = System.Drawing.SystemColors.GrayText;
            this.label4.Location = new System.Drawing.Point(54, 189);
            this.label4.Name = "label4";
            this.label4.Size = new System.Drawing.Size(240, 11);
            this.label4.TabIndex = 8;
            this.label4.Text = "（默认为txt格式，如不需要可选择fasle关闭）";
            // 
            // fileName
            // 
            this.fileName.Location = new System.Drawing.Point(129, 108);
            this.fileName.Name = "fileName";
            this.fileName.Size = new System.Drawing.Size(159, 21);
            this.fileName.TabIndex = 14;
            // 
            // settingCancel
            // 
            this.settingCancel.Location = new System.Drawing.Point(62, 242);
            this.settingCancel.Name = "settingCancel";
            this.settingCancel.Size = new System.Drawing.Size(75, 23);
            this.settingCancel.TabIndex = 10;
            this.settingCancel.Text = "取消";
            this.settingCancel.UseVisualStyleBackColor = true;
            this.settingCancel.Click += new System.EventHandler(this.settingCancel_Click);
            // 
            // settingOK
            // 
            this.settingOK.Location = new System.Drawing.Point(213, 242);
            this.settingOK.Name = "settingOK";
            this.settingOK.Size = new System.Drawing.Size(75, 23);
            this.settingOK.TabIndex = 11;
            this.settingOK.Text = "确认";
            this.settingOK.UseVisualStyleBackColor = true;
            this.settingOK.Click += new System.EventHandler(this.settingOK_Click);
            // 
            // textPath
            // 
            this.textPath.Location = new System.Drawing.Point(129, 65);
            this.textPath.Name = "textPath";
            this.textPath.ScrollBars = System.Windows.Forms.ScrollBars.Horizontal;
            this.textPath.Size = new System.Drawing.Size(91, 21);
            this.textPath.TabIndex = 12;
            this.textPath.WordWrap = false;
            this.textPath.TextChanged += new System.EventHandler(this.textPath_TextChanged);
            this.textPath.KeyPress += new System.Windows.Forms.KeyPressEventHandler(this.textPath_KeyPress);
            // 
            // selectPath
            // 
            this.selectPath.Location = new System.Drawing.Point(226, 65);
            this.selectPath.Name = "selectPath";
            this.selectPath.Size = new System.Drawing.Size(66, 21);
            this.selectPath.TabIndex = 13;
            this.selectPath.Text = "选择路径";
            this.selectPath.UseVisualStyleBackColor = true;
            this.selectPath.Click += new System.EventHandler(this.selectPath_Click);
            // 
            // settingForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 12F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(356, 309);
            this.Controls.Add(this.selectPath);
            this.Controls.Add(this.textPath);
            this.Controls.Add(this.settingOK);
            this.Controls.Add(this.settingCancel);
            this.Controls.Add(this.fileName);
            this.Controls.Add(this.label4);
            this.Controls.Add(this.label3);
            this.Controls.Add(this.label2);
            this.Controls.Add(this.label1);
            this.Controls.Add(this.fileType);
            this.Name = "settingForm";
            this.Text = "设置";
            this.Load += new System.EventHandler(this.settingForm_Load);
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion
        private System.Windows.Forms.ComboBox fileType;
        private System.Windows.Forms.Label label1;
        private System.Windows.Forms.Label label2;
        private System.Windows.Forms.Label label3;
        private System.Windows.Forms.Label label4;
        private System.Windows.Forms.TextBox fileName;
        private System.Windows.Forms.Button settingCancel;
        private System.Windows.Forms.Button settingOK;
        private System.Windows.Forms.FolderBrowserDialog folderBrowserDialog;
        private System.Windows.Forms.TextBox textPath;
        private System.Windows.Forms.Button selectPath;
    }
}
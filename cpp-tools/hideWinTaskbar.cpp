#include <stdio.h>
#include <windows.h>
#include <iostream>

using namespace std;

void HideWindow() {
	HWND hwnd = GetForegroundWindow();
	if (hwnd) {
		ShowWindow(hwnd, SW_HIDE);
	}
}

int main(int argc, char* argv[])
{
	Sleep(1000);
	HideWindow();
	HWND task;
	task = FindWindow(L"Shell_TrayWnd", NULL);
	ShowWindow(task, SW_HIDE);//隐藏任务栏

	//system("pause");
	if (0 == RegisterHotKey(NULL, 1, MOD_SHIFT, VK_F1)) {
		cout << "HotKey error : " << GetLastError() << endl;
	}
	if (0 == RegisterHotKey(NULL, 2, MOD_SHIFT, VK_F2)) {
		cout << "HotKey error : " << GetLastError() << endl;
	}
	if (0 == RegisterHotKey(NULL, 3, MOD_SHIFT, VK_F3)) {
		cout << "exit error : " << GetLastError() << endl;
	}

	// 消息循环
	MSG msg = { 0 };
	while (GetMessage(&msg, NULL, 0, 0)) {
		switch (msg.message) {
		case WM_HOTKEY:
		{
			if (1 == msg.wParam) {
				ShowWindow(task, SW_SHOW);//显示任务栏
			}
			else if (2 == msg.wParam) {
				ShowWindow(task, SW_HIDE);//隐藏任务栏
			}
			else if (3 == msg.wParam) {
				ShowWindow(task, SW_SHOW);
				return 0;  // 显示任务栏并退出hide程序
			}
			break;
		}
		default:
			break;
		}

	}
	cout << "finished." << endl;
	return 0;
}


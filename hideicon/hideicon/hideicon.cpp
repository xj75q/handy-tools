#include <windows.h>
#include <winuser.h>
#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <shobjidl_core.h>
#include <stdlib.h>
using namespace std;
typedef vector <int> HOTKEY;

map <HOTKEY, int> AreDown;
bool HotKeyDown(HOTKEY hotkey)
{
	for (auto& i : hotkey)
		if (!(GetAsyncKeyState(i) & 0x8000))
		{
			AreDown[hotkey] = false;
			return false;
		}
	if (AreDown[hotkey])
		return false;
	else
	{
		AreDown[hotkey] = true;
		return true;
	}
}

HOTKEY add = { VK_MENU, '1' };
HOTKEY hide = { VK_MENU, '2' };
HOTKEY display = { VK_MENU, '3' };

vector <string> hidewindows;
bool ishide;
int main()
{
	
	
	system("title HideIcon");
	HWND hWnd = ::GetForegroundWindow();
	::SetWindowPos(hWnd, HWND_TOPMOST, 0, 0, 600, 400, SWP_NOMOVE);


	printf("\n 将鼠标放在需要隐藏的窗口界面，按alt + 1 选中窗口\n");
	for (;;)
	{
		if (HotKeyDown(add))
		{
			system("cls");
			HWND hWnd = GetForegroundWindow();
			char text[0xFF];
			GetWindowTextA(
				hWnd,
				text,
				0xFF
			);
			string str = text;
			hidewindows.push_back(str);
			for (auto& i : hidewindows)
				//cout << "\n你当前选择的窗口为：" << i << endl;
			    printf("\n 你当前选择的窗口为:%s\n", i.data());
			cout << endl;
		}

		if (HotKeyDown(hide))
		{
			system("cls");
			for (auto& i : hidewindows)
			{
				printf("\n [ %s ]的窗口将在任务栏隐藏...\n", i.data());

				HWND targetWindow = FindWindowA(NULL, i.data());

				if (targetWindow != NULL)
				{

					SetWindowLongPtrW(targetWindow, GWL_EXSTYLE,WS_EX_TOOLWINDOW);
					SetWindowPos(targetWindow, NULL, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_FRAMECHANGED);
				
					//cout << "该任务栏图标隐藏成功，按回车键退出..." << endl;
					printf("\n 该任务栏图标隐藏成功，按回车键退出... \n");
					getchar();
					exit(0);

				}
				else
				{
					std::cout << "未找到目标窗口" << std::endl;
				}
				
			}
			ishide = ishide ? false : true;
		}
		
		if (HotKeyDown(display))
		{
			system("cls");
			for (auto& i : hidewindows)
			{
				
				printf("\n [ %s ]的窗口将在任务栏显示...\n", i.data());

				HWND targetWindow = FindWindowA(NULL, i.data());

				if (targetWindow != NULL)
				{

					SetWindowLongPtrW(targetWindow, GWL_EXSTYLE, WS_EX_APPWINDOW);
					SetWindowPos(targetWindow, NULL, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_FRAMECHANGED);
				}
				else
				{
					std::cout << "未找到目标窗口" << std::endl;
				}

			}
			ishide = ishide ? false : true;
		}
	}
}

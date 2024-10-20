!e::
IfWinExist,ahk_exe Evernote.exe
{
    WinActivate
}
else
{
    Run \Yinxiang Biji\Evernote.exe
    WinWait,ahk_exe Evernote.exe
    WinActivate
}
return


!r::
IfWinExist,ahk_exe msedge.exe
{
    WinActivate
}
else
{
    Run C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe
    WinWait,ahk_exe msedge.exe
    WinActivate
}
return


!x::
IfWinExist,ahk_exe 喜马拉雅.exe
{
    WinActivate
}
else
{
    Run \喜马拉雅.exe
    WinWait,ahk_exe 喜马拉雅.exe
    WinActivate
}
return


!p::
IfWinExist,ahk_exe cloudmusic.exe
{
    WinActivate
}
else
{
    Run \网易云音乐PC版\cloudmusic.exe
    WinWait,ahk_exe cloudmusic.exe
    WinActivate
}
return


!o::
IfWinExist,ahk_exe QQMusic.exe
{
    WinActivate
}
else
{
    Run \qqmusic\QQMusic.exe
    WinWait,ahk_exe QQMusic.exe
    WinActivate
}
return






!m::Send {Volume_Mute}
!q::WinMinimize,A
!l::WinMaximize,A 
!k::WinRestore,A

Up::
MouseClick,WheelUp,,,1,0,D,R
return


Down::
MouseClick,WheelDown,,,1,0,D,R
return

CapsLock & Up:: Up
CapsLock & Down:: Down

!,::MouseMove, 228 , 416
!.::MouseMove, 1171 , 416

!Enter::Click


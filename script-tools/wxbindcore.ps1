$wxid=(Get-Process -Name "Weixin" | Select-Object).Id
foreach ($element in $wxid) {
  $process = Get-Process -Id $element

  $process.ProcessorAffinity = [System.IntPtr]128 #第7个核心
  $process.MaxWorkingSet = 90mb
}
$wxocrid=(Get-Process -Name "WeChatAppEx" | Select-Object).Id
foreach ($element in $wxocrid) {
  $process = Get-Process -Id $element
  $process.ProcessorAffinity = [System.IntPtr]64 #第6个核心
  $process.MaxWorkingSet = 40mb
}
$edgeid=(Get-Process -Name "msedge*" | Select-Object).Id
foreach ($element in $edgeid) {
  $process = Get-Process -Id $element
  $process.ProcessorAffinity = 0x0009 #绑定0和3个。
}


Param(
    [string]$command,
    [string]$yamlFile
  )
$escCommand =  $command -replace '"','\"'
$evalArgs = @("eval", "$escCommand", $yamlFile)

Write-Output ".\yq.exe"  $evalArgs 
& ".\yq.exe" $evalArgs 
# Start-Process ".\yq.exe" -ArgumentList $evalArgs -Wait -NoNewWindow -PassThru


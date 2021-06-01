
# https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_functions_advanced_parameters?view=powershell-7.1
Param(
    [Parameter(Position=0, Mandatory=$true, ValueFromPipeline=$false,
    HelpMessage="enter valid command, {eval }")]
    [ValidateSet("e","eval")]
    [String[]]
    $Command,

    [Parameter(Position=1, Mandatory=$true, ValueFromPipeline=$false)]
    [String[]]
    $Filter,
    
    [Parameter(Mandatory=$false, ValueFromPipeline=$false)]
    [Switch]
    $flags,
    
    [Parameter(Position=1, Mandatory=$true, ValueFromPipeline=$true)]
    [SupportsWildcards()]
    [String[]]
    [string]$YamlFilePath
  )
$escFilter =  $Filter -replace '"','\"'
$evalArgs = @("$Command", "$escFilter", $YamlFilePath)

Write-Output ".\yq.exe"  $evalArgs 
& "${PSScriptRoot}\yq.exe" $evalArgs

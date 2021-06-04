
# https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_functions_advanced_parameters?view=powershell-7.1
Param(
    [Parameter(Position=1, Mandatory=$true, ValueFromPipeline=$false,
    HelpMessage="enter valid command, {eval }")]
    [ValidateSet("e","eval","all")]
    [String]
    $Command,

    [Parameter(Position=2, Mandatory=$true, ValueFromPipeline=$false)]
    [String]
    $Filter,

    [Parameter(Mandatory=$false, ValueFromPipeline=$true)]
    [SupportsWildcards()]
    [String[]]
    $FilePath,

    [Parameter(Position=2, Mandatory=$true, ValueFromPipeline=$false,
    HelpMessage="sets indent level for output")]
    [Int16]
    $Indent,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
            HelpMessage="force print with colors")]
    [Switch]
    $ColorsSwitch,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
            HelpMessage="force print with no colors")]
    [Switch]
    $NoColorsSwitch,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
            HelpMessage="set exit status if there are no matches or null or false is returned")]
    [Switch]
    $ExitStatusSwitch,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
            HelpMessage="update the yaml file inplace of first yaml file given")]
    [Switch]
    $InPlaceSwitch,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
            HelpMessage="don't print document separators (---)")]
    [Switch]
    $NoDocSepSwitch,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
            HelpMessage="Don't read input, simply evaluate the expression given. Useful for creating yaml docs from scratch.")]
    [Switch]
    $NullInputSwitch,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
            HelpMessage="output as json. Set indent to 0 to print json in one line.")]
    [Switch]
    $ToJsonSwitch,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
    HelpMessage="pretty print, shorthand for '... style = `"`"'")]
    [Switch]
    $PrettyPrintSwitch,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
    HelpMessage="unwrap scalar, print the value with no quotes, colors or comments (default true)")]
    [Switch]
    $UnwrapScalarSwitch,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
            HelpMessage="verbose mode")]
    [Switch]
    $VerboseSwitch,

    [Parameter(Mandatory=$false, ValueFromPipeline=$false,
    HelpMessage="print version and quit")]
    [Switch]
    $VersionSwitch
)
$escFilter =  $Filter -replace '"','\"'
$evalArgs = @("$Command", "$escFilter", $FilePath)

Write-Output ".\yq.exe"  $evalArgs 
& "${PSScriptRoot}\yq.exe" $evalArgs

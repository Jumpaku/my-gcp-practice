# spreadsheets-drive

## Prerequisites

- Add a GCP project
    - Enable Google Drive API in API and Services
    - Enable Google Forms API in API and Services
    - Create a service account with an access key
- Add a shared drive in Google Workspace
    - Add the service account to the shared drive as a content manager

## Examples

Export environment variable for authentication:

```sh
export GOOGLE_APPLICATION_CREDENTIALS=credentials/my-gcp-practice-477709-t4-cc9c47f8716c.json
```

### clone-spreadsheet-in-shared-drive-and-update

```sh
xtracego run clone-spreadsheet-in-shared-drive/main.go -- \
  -source="${ORIGINAL_SPREADSHEET_ID}" \
  -destination="${FOLDER_ID_WHERE_SPREADSHEET_IS_CLONED}" \
  -filename="${CLONED_SPREADSHEET_NAME}"
```

### get-formatted-sheet-data

```shell
xtracego run get-formatted-sheet-data/main.go -- -spreadsheet="${SPREADDSHEET_ID}"
```

### save-user-entered-sheet-data

```shell
xtracego run save-user-entered-sheet-data/main.go -- -spreadsheet="${SPREADDSHEET_ID}" -sheet="${SHEET_ID}" \
  < save-user-entered-sheet-data/example.csv
```

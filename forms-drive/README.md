# forms-drive

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

### clone-form-in-shared-drive-and-update

```sh
xtracego run clone-form-in-shared-drive/main.go -- \
  -source="${ORIGINAL_FORM_ID}" \
  -destination="${FOLDER_ID_WHERE_FORM_IS_CLONED}" \
  -filename="${CLONED_FORM_NAME}"
```

### find-all-forms-in-shared-drive

```shell
xtracego run find-all-forms-in-shared-drive/main.go -- -folder="${PARENT_FOLDER_ID}"
xtracego run find-all-forms-in-shared-drive/main.go -- -folder="${PARENT_FOLDER_ID}" -filename="${KEYWORD_INCLUDED_IN_FILENAME}"
xtracego run find-all-forms-in-shared-drive/main.go -- -form="${FORM_ID}"
```

### list-form-responses

```shell
xtracego run list-form-responses/main.go -- -form="${FORM_ID}"
```

### change-access-and-publish-state

```shell
xtracego run change-access-and-publish-state/main.go -- -form="${FORM_ID}" -publish=accepting -access=anyone 
xtracego run change-access-and-publish-state/main.go -- -form="${FORM_ID}" -publish=accepting -access=domain -domain="${DOMAIN}"
xtracego run change-access-and-publish-state/main.go -- -form="${FORM_ID}" -publish=accepting -access=limited 
xtracego run change-access-and-publish-state/main.go -- -form="${FORM_ID}" -publish=not_accepting -access=anyone 
xtracego run change-access-and-publish-state/main.go -- -form="${FORM_ID}" -publish=not_accepting -access=domain -domain="${DOMAIN}"
xtracego run change-access-and-publish-state/main.go -- -form="${FORM_ID}" -publish=not_accepting -access=limited 
xtracego run change-access-and-publish-state/main.go -- -form="${FORM_ID}" -publish=not_published -access=anyone 
xtracego run change-access-and-publish-state/main.go -- -form="${FORM_ID}" -publish=not_published -access=domain -domain="${DOMAIN}" 
xtracego run change-access-and-publish-state/main.go -- -form="${FORM_ID}" -publish=not_published -access=limited 
```

### update-form

```shell
xtracego run update-form/main.go -- -form="${FORM_ID}" -title="${UPDATED_TITLE} -description="${UPDATED_DESCRIPTION}"
```

# Cloud Build trigger for GitHub (gcloud commands)

1. Install Cloud Build GitHub App into a selected GitHub repository
   1. https://github.com/apps/google-cloud-build
2. Create GitHub PAT and Save in Secret Manager
   1. scopes:
      - repo
      - read:user
      - (read:org) 
3. Allow Cloud Build service account to access Secret Manager
   1. member: service-${PROJECT_NUMBER}@gcp-sa-cloudbuild.iam.gserviceaccount.com
   2. role: roles/secretmanager.secretAccessor
4. Create CloudBuild connection with the PAT and the App ID
    ```shell
    gcloud builds connections create github my-gcp-practice-github-connection \
      --project=${PROJECT_ID} \
      --region=asia-northeast1 \
      --app-installation-id="${APP_INSTALLATION_ID}" \
      --authorizer-token-secret-version="projects/${PROJECT_ID}/secrets/GITHUB_PAT/versions/latest"
    ```
5. Create CloudBuild repository with the connection
    ```shell
    gcloud builds repositories create my-gcp-practice-github-repository \
        --project=${PROJECT_ID} \
        --connection=my-gcp-practice-github-connection \
        --region=asia-northeast1 \
        --remote-uri=https://github.com/Jumpaku/my-gcp-practice.git
    ```
6. Add a webhook secret and save in Secret Manager 
7. Create CloudBuild trigger with the repository
    ```shell
    gcloud builds triggers create github \
        --name="my-github-trigger" \
        --repo-name="my-gcp-practice-github-repository" \
        --repo-owner="Jumpaku" \
        --branch-pattern="^main$"
    ```
8. Set up a GitHub repository webhook


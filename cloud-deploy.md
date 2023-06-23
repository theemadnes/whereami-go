# instructions for setting up Cloud Deploy for `whereami-go`

make sure Cloud Deploy API is enabled
```
export PROJECT_ID=cicd-system-demo-01

gcloud services enable clouddeploy.googleapis.com --project ${PROJECT_ID}
```

set up perms per https://cloud.google.com/deploy/docs/deploy-app-gke

```
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member=serviceAccount:$(gcloud projects describe ${PROJECT_ID} \
    --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
    --role="roles/clouddeploy.jobRunner"

gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member=serviceAccount:$(gcloud projects describe ${PROJECT_ID} \
    --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
    --role="roles/container.developer"

gcloud iam service-accounts add-iam-policy-binding $(gcloud projects describe ${PROJECT_ID} \
    --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
    --member=serviceAccount:$(gcloud projects describe ${PROJECT_ID} \
    --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
    --role="roles/iam.serviceAccountUser" \
    --project=${PROJECT_ID}
```

set up targets & pipeline

```
gcloud deploy apply --file=cloud-deploy/clouddeploy.yaml --region=us-central1 --project=${PROJECT_ID}
```

```
gcloud deploy releases create test-release-001 \
  --project=${PROJECT_ID} \
  --region=us-central1 \
  --delivery-pipeline=whereami-go-01 \
  --images=whereami-go-image=us-central1-docker.pkg.dev/cicd-system-demo-01/cicd-demo-01/whereami@sha256:ea6cc4e509b41949fe6a3b0cfb312adea6940ef2c2766f244116822ddc7cedfb
```
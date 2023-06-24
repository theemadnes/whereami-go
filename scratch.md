# testing out Binary Auth for `whereami-go`

```
# using https://cloud.google.com/binary-authorization/docs/cloud-build

PROJECT_ID=cicd-system-demo-01

#gcloud config set project ${PROJECT_ID}

PROJECT_NUMBER=$(gcloud projects list --filter="${PROJECT_ID}" --format="value(PROJECT_NUMBER)")

gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
  --role roles/binaryauthorization.attestorsViewer

gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
  --role roles/cloudkms.signerVerifier

gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member serviceAccount:${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com \
  --role roles/containeranalysis.notes.attacher
```

 - deployment/whereami: creating container whereami
    - pod/whereami-7ff67b8d85-cvqfq: FailedKillPod: error killing pod: failed to "KillPodSandbox" for "e0f20c64-4a75-48fc-b83f-624abf837acd" with KillPodSandboxError: "rpc error: code = Unknown desc = failed to destroy network for sandbox \"4a4fec98ca945c70549bcc8ecf158a7cac5e552d5f7f48f2992c57293e82ce23\": plugin type=\"cilium-cni\" failed (delete): failed to find plugin \"cilium-cni\" in path [/home/kubernetes/bin]"
    - pod/whereami-7ff67b8d85-r62sc: creating container whereami
 - deployment/whereami: container whereami terminated with exit code 2

kubectl port-forward svc/whereami 1234:80
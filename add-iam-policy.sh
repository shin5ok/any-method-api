
gcloud projects add-iam-policy-binding --member=serviceAccount:$(gcloud projects describe --format=json $PROJECT | jq .projectNumber -r)-compute@developer.gserviceaccount.com --role=roles/editor $PROJECT

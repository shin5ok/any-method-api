IMAGE=$1

gcloud builds submit -t $IMAGE --project=$GOOGLE_CLOUD_PROJECT

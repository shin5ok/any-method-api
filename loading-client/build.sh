IMAGE=$1

gcloud builds submit -t $IMAGE --project=$PROJECT

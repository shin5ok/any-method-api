date=$(date '+%Y%m%d%H%M')
TAG=$date
echo "Building with TAG=$TAG"
gcloud builds submit -t gcr.io/$PROJECT/loading-client:$TAG --project=$PROJECT

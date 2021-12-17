TAG=${TAG:-0.01}
echo "Building with TAG=$TAG"
gcloud builds submit -t gcr.io/$PROJECT/loading-client:$TAG --project=$PROJECT

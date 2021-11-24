TAG=${TAG:-0.01}
echo "Building with TAG=$TAG"
gcloud builds submit --pack=image=gcr.io/$PROJECT/any-method-api:$TAG --project=$PROJECT

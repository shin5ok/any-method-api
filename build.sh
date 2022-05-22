date=$(date '+%Y%m%d%H%M')
TAG=$date
APPNAME=any-method-api
echo "Building with TAG=$TAG"
gcloud builds submit --pack=image=gcr.io/$PROJECT/${APPNAME}:$TAG --project=$PROJECT

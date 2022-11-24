date=$(date '+%Y%m%d%H%M')
TAG=${TAG:-${date}}
echo "Building with TAG=$TAG"
IMAGE=asia-northeast1-docker.pkg.dev/$PROJECT/my-app/$APPNAME

gcloud builds submit -t $IMAGE:$TAG --project=$PROJECT

date=$(date '+%Y%m%d%H%M')
TAG=${TAG:-${date}}
APPNAME=${1:-any-method-api}
echo "Building with TAG=$TAG"
gcloud builds submit -t gcr.io/$PROJECT/${APPNAME}:$TAG --project=$PROJECT

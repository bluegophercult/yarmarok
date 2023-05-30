ARG GCLOUD_SDK_VERSION=335.0.0-alpine

FROM google/cloud-sdk:$GCLOUD_SDK_VERSION

RUN apk add --update --no-cache openjdk8-jre

RUN gcloud components install cloud-firestore-emulator beta --quiet

ENV FIRESTORE_PROJECT_ID "local-test-firestore-project-id"
RUN gcloud config set project "${FIRESTORE_PROJECT_ID}"

ENV PORT 8080
EXPOSE "$PORT"

ENTRYPOINT ["gcloud", "beta", "emulators", "firestore", "start", "--host-port=0.0.0.0:8080"]

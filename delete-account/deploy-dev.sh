gcloud functions deploy DeleteUserData \
--trigger-event providers/firebase.auth/eventTypes/user.delete \
--trigger-resource linum-dev \
--source src \
--region europe-west1 \
--runtime go116
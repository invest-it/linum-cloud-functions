gcloud functions deploy DeleteUserData \
--trigger-event providers/firebase.auth/eventTypes/user.delete \
--trigger-resource linum-5d9f6 \
--source src \
--runtime go116
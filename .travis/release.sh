#!/bin/bash
set -e # Exit with nonzero exit code if anything fails

# Pull requests and commits to other branches shouldn't try to deploy
if [ "$TRAVIS_PULL_REQUEST" != "false" -o "$TRAVIS_BRANCH" != "master" ]; then
    exit 0
fi

# Save some useful information
REPO=`git config remote.origin.url`
SSH_REPO=${REPO/https:\/\/github.com\//git@github.com:}
git remote set-url origin "$SSH_REPO"

git config user.name "Travis CI"
git config user.email "travis@travis-ci.org"

# Get the deploy key by using Travis's stored variables to decrypt deploy_key.enc
ENCRYPTED_KEY_VAR="encrypted_${ENCRYPTION_LABEL}_key"
ENCRYPTED_IV_VAR="encrypted_${ENCRYPTION_LABEL}_iv"
ENCRYPTED_KEY=${!ENCRYPTED_KEY_VAR}
ENCRYPTED_IV=${!ENCRYPTED_IV_VAR}
openssl aes-256-cbc -K $ENCRYPTED_KEY -iv $ENCRYPTED_IV -in github_deploy_key.enc -out github_deploy_key -d
chmod 600 github_deploy_key
eval `ssh-agent -s`
ssh-add github_deploy_key

# Now that we're all set up, we can release.
make release
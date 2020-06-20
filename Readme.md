# Long Trail

Long Trail is a web application built specifically for [GMCMF](https://gmcmf.org) students to easily keep track of when and how long they practice.

Long Trail is hosted on AWS and uses [dispatch](https://github.com/flick-web/dispatch).

## Deploying

Deploying your own version of Long Trail is fairly straightforward. You need an AWS account and a few command line tools, but everything else is fairly automatic.

Prerequisites:

- An AWS account
- Correctly configured `awscli` (or `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` set if running from GitLab CI)
- `jq` installed
- (Optional) Somewhere to host a static website (GitHub/GitLab pages, Netlify, etc)

Once everything is set up, just clone this repository, `cd` into it, and run:

```sh
./scripts/deploy.sh # Deploy the backend on AWS
./scripts/generate-config.sh # Generate website configuration from the deployed backend
cd web # Change to the frontend directory
npm install # Install dependencies
npm install --global parcel-bundler # Install the parcel CLI tool
parcel index.html # Run the site locally for testing
```

You can then view the website locally (on http://localhost:1234 by default). The website will automatically be configured to use the your backend stack on AWS, so you can register and start using it! To host the site somewhere, you'll just need to build the site with parcel and upload the resulting files to any kind of static site host.

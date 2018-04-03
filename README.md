# Tumbler to Twitter

*Stefan Arentz, April 2018*

This is an Amazon Lambda function that periodically takes a look at a Tumblr photo blog and mirrors all posts there to Twitter. It does this on a 5 minute schedule, so posts will be at worst delayed by a few minutes.

I use this as a workaround for the extremely limited way Instagram can post to Twitter; instead of telling Instagram to post to Twitter, I tell it to post to Tumblr and this function then *properly* posts the photo to Twitter with correctly inlined media. This means it will show up nicely in your feed instead of just with a link to instagram.

## Deploying to AWS

```
env GOOS=linux go build && zip tumblr2twitter.zip tumblr2twitter
aws s3 mb s3://tumblr2twitter-artifacts
sam package --template-file template.yaml --s3-bucket tumblr2twitter-artifacts \
  --output-template-file package.yaml
sam deploy --template-file ./package.yaml --stack-name tumblr2twitter \
  --capabilities CAPABILITY_IAM
```

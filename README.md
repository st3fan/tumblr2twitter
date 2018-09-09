# Tumbler to Twitter

*Stefan Arentz, April 2018*

This is a small program that you can run from cron to periodically look at a Tumblr photo blog and mirror all posts there to Twitter.

I use this as a workaround for the extremely limited way Instagram can post to Twitter; instead of telling Instagram to post to Twitter, I tell it to post to Tumblr and this function then *properly* posts the photo to Twitter with correctly inlined media. This means it will show up nicely in your feed instead of just with a link to instagram.

You need to set the following environment variables to make this work


- TUMBLR_BLOG_NAME
- TUMBLR_CONSUMER_KEY
- TUMBLR_CONSUMER_SECRET
- TWITTER_ACCESS_TOKEN
- TWITTER_ACCESS_TOKEN_SECRET
- TWITTER_CONSUMER_KEY
- TWITTER_CONSUMER_SECRET

This README is mostly a reminder to myself how this all works. But if you find this useful, drop me a line. You can also see my posted photos at [twitter.com/satefan/media](https://twitter.com/satefan/media).


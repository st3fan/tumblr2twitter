// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/boltdb/bolt"
	"github.com/st3fan/html2text"
	"github.com/st3fan/tumblrclient"
	tumblr "github.com/tumblr/tumblr.go"
)

func download(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Download failed with status %d %s", res.StatusCode, res.Status)
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func postTweetWithMedia(twitter *anaconda.TwitterApi, status string, media anaconda.Media) (anaconda.Tweet, error) {
	values := url.Values{}
	values.Set("media_ids", strconv.FormatInt(media.MediaID, 10))
	return twitter.PostTweet(status, values)
}

func haveSeenPostBefore(db *bolt.DB, ID int) (bool, error) {
	seen := false
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Posts"))
		if bucket == nil {
			return nil
		}
		seen = bucket.Get([]byte(strconv.Itoa(ID))) != nil
		return nil
	})
	return seen, err
}

func rememberPost(db *bolt.DB, ID int, tweetURL string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("Posts"))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(strconv.Itoa(ID)), []byte(tweetURL))
	})
}

func cleanupCaption(caption string) string {
	re := regexp.MustCompile(`<a[^>]*>(.*?)</a>`)
	caption = strings.TrimSpace(html2text.HTML2Text(re.ReplaceAllString(caption, "")))
	caption = strings.ReplaceAll(caption, "&hellip;", "\U00002026")
	return caption
}

func main() {
	db, err := bolt.Open("/home/stefan/etc/tumblr2twitter.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	client := tumblrclient.NewClient(os.Getenv("TUMBLR_CONSUMER_KEY"), os.Getenv("TUMBLR_CONSUMER_SECRET"))
	// TODO Verify that this works

	twitter := anaconda.NewTwitterApiWithCredentials(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"), os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	// TODO Verify that this works

	blog := client.GetBlog(os.Getenv("TUMBLR_BLOG_NAME"))
	posts, err := blog.GetPosts(nil)
	if err != nil {
		log.Fatal("Failed to get Tumblr blog posts: ", err)
	}

	for i := 0; i < int(posts.TotalPosts); i++ {
		post := posts.Get(uint(i))
		if photoPost, ok := post.(*tumblr.PhotoPost); ok && photoPost != nil {
			log.Printf("Looking at blog post <%d> at <%s>", photoPost.Id, photoPost.PostUrl)

			// We can currently only handle posts with one photo attached
			if len(photoPost.Photos) == 1 {
				// Check if we already posted this to Twitter
				seenPostBefore, err := haveSeenPostBefore(db, int(photoPost.Id))
				if err != nil {
					log.Printf("Failed to check if we have seen post before: %s\n", err)
					continue
				}

				if seenPostBefore {
					log.Printf("Seen post before, skipping")
					continue
				}

				log.Printf("Not seen post before, going to post to Twitter")

				// Even if the upload fails, we remember this as done
				if err := rememberPost(db, int(photoPost.Id), ""); err != nil {
					log.Printf("Failed to remember post: %s", err)
				}

				// Fetch the photo from Tumblr
				image, err := download(photoPost.Photos[0].OriginalSize.Url)
				if err != nil {
					log.Printf("Failed to download image from post <%d>: %s", photoPost.Id, err)
					continue
				}

				// Upload the photo to Twitter
				media, err := twitter.UploadMedia(base64.StdEncoding.EncodeToString(image))
				if err != nil {
					log.Printf("Failed to upload image to Twitter: %s", err)
					continue
				}

				// Post a tweet with the media we just uploaded
				tweet, err := postTweetWithMedia(twitter, cleanupCaption(photoPost.Caption), media)
				if err != nil {
					log.Printf("Failed to post tweet to Twitter: %s", err)
					continue
				}
				
				tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.User.ScreenName, tweet.IdStr)
				log.Printf("Posted Tumblr post <%s> as tweet <%s>\n", photoPost.PostUrl, tweetURL)

				if err := rememberPost(db, int(photoPost.Id), tweetURL); err != nil {
					log.Printf("Failed to remember post: %s", err)
				}
			}
		}
	}
}

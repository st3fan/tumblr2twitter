package main

import (
	"testing"
)

func Test_CleanupCaption(t *testing.T) {
	caption := cleanupCaption("<p>Beautiful day today. <br/><a href=\"https://www.instagram.com/p/BnphlH-FaVB6_FWjwqVAgIZEQAjswSD90hJ_bw0/?utm_source=ig_tumblr_share&amp;igshid=9qod8zb07dfw\">https://www.instagram.com/p/BnphlH-FaVB6_FWjwqVAgIZEQAjswSD90hJ_bw0/?utm_source=ig_tumblr_share&amp;igshid=9qod8zb07dfw</a></p>")
	if caption != "Beautiful day today." {
		t.Fatal("Not expected output: ", caption)
	}
}

func Test_CleanupCaption_Hellip(t *testing.T) {
	caption := cleanupCaption(`<p>Testing &hellip; ignore please &hellip; <br/> <a href="https://www.instagram.com/p/B7MK1qGJkZCelx898PUUoIwU3Aq18RX5d1TKaI0/?igshid=p1w1z9s9wp1">https://www.instagram.com/p/B7MK1qGJkZCelx898PUUoIwU3Aq18RX5d1TKaI0/?igshid=p1w1z9s9wp1</a></p>`)
	if caption != "Testing \U00002026 ignore please \U00002026" {
		t.Fatal("Not excpted output: ", caption)
	}
}


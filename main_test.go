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

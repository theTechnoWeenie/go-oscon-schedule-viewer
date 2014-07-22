OSCON 2014 DIY Schedule Viewer
========================

This is a super quick and dirty oscon schedule view that was written in the go language.

I wanted to learn the go language (golang.org) a little better, and I liked the challenge of creating a schedule viewer.

This will query the URL (http://www.oreilly.com/pub/sc/osconfeed), parse the JSON for the event, and print the relavent (`mySessionSerails`) sessions to the console displaying the start time, end time, location and whether or not the session is currently... well... in session.

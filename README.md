Snap Release Email Generator
============================

This small utility helps generating the release emails for snaps. It
uses the trello board to retrieve the curently released snaps and
construct email message.

Requirements
============

trello_secrets.json file that contains the Application Key and Token.

```
konrad at annapurna in snap-release-email-generator (master) % cat trello_secrets.json 
{"app_id":"APP ID","token":"TOKEN"}

```

Usage
====

Usage of ./snap-release-email-generator:
  -publisher string
        Snap publisher's name (default "Simon")
  -secrets string
        Trello Secrets configuration (default "trello_secrets.json")

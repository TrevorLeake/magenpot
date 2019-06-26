# Magenpot (Magento Honeypot)

*Built as a detached fork from [Drupot](www.github.com/d1str0/Drupot) (Drupal honeypot).*

## Installation
Magenpot supports go modules.

`go get github.com/trevorleake/magenpot`

`go build`

## Running Magenpot
`./magenpot -c config.toml`

## Identifying the Magento Version
A simple way to identify the version of Magento a site (example.com) is running, is to try the the magento_version path (example.com/magento_version). This will not always work, as the file doesn't have to be served.

## Configuration
`config.toml.example` contains an example of *all* currently available configuration options.

### Magento
    [magento]
    port = 80
    magento_version_text = Magento/2.3 (Enterprise)

`port` allows you to set the http port to listen on. Currently, this is only ever served over http. Future versions will support https.

`magento_version_text` allows you to set what is returned in the magento_version file and thereby mimic different versions of Magento.

### hpfeeds
    [hpfeeds]
    enabled = true
    host = "hpfeeds.threatstream.com"
    port = 10000
    ident = "magenpot"
    auth = "somesecret"
    channel = "magenpot.events"
    meta = "Magento scan event detected"

hpfeeds can be enabled for logging if wanted. Supply host, port, ident, auth,
and channel information relevant to an hpfeeds broker you want to report to.

`meta` provides a static string to send in every hpfeeds request. Could be use
to differentiate Magneto versions hosted by honeypot or used to differentiate
Magenpot data in busy hpfeeds channels.

### Fetch Public IP
    [fetch_public_ip]
    enabled = true
    urls = ["http://icanhazip.com/", "http://ifconfig.me/ip"]


If enabled, Magenpot will attempt to fetch the public IP of itself from the listed URLs. If enabled and no public IP can be fetched, Magenpot will quit.

### Example Magento Sites
Try appending the "magento_version" path to check these sites' versions.
* https://www.bulkpowders.co.uk/

Please be nice to these sites. :)

package app

import (
	"net/http"
	"net/url"
	"time"

	_ "github.com/mailru/easyjson/gen"
)

// A NelReport describes a single network error report.
//
// (The name is a bit of a misnomer; it can also represent other report payloads
// uploaded via the Reporting API, in which case ReportType will tell you what
// kind of report it is, and RawBody will contain the unparsed JSON `body`
// field.  It can also represent information about successful HTTP requests,
// collected and delivered via a user agent's NEL stack, in which case the Type
// field will be "ok".)
type NelReport struct {
	// The number of milliseconds between when the report was generated by
	// the user agent and when it was uploaded.
	Age int
	// The type of report.  For NEL, this will be "network-error".
	ReportType string
	// The URL of the request that this report describes.
	URL string
	// UserAgent represents the value of the User-Agent header in the request that the report is about
	UserAgent string
	// The referrer information of the request, as determined by the
	// referrer policy associated with its client.
	Referrer string
	// The active sampling rate for this request, expressed as a fraction
	// between 0.0 and 1.0 (inclusive).
	SamplingFraction float32
	// The IP address of the host to which the user agent set the request.
	ServerIP string
	// The ALPN ID  of the network protocol used to fetch the resource.
	Protocol string
	// The method of the HTTP request (e.g. GET, POST)
	Method string
	// The status code of the HTTP response, if available.
	StatusCode int
	// The elapsed number of milliseconds between the start of the resource
	// fetch and when it was aborted by the user agent.
	ElapsedTime int
	// The phase of the request in which the failure occurred. One of
	// {dns, connection, application}; a successful request always has a value of application
	Phase string
	// The description of the error type.  For reports about successful
	// requests, this will be "ok".  See the NEL spec for the authoritative
	// list of possible values for failed requests.
	Type string

	// For non-NEL reports, this will contain the unparsed JSON content of
	// the report's `body` field.
	RawBody []byte

	// An arbitrary set of extra data that you can attach to your reports.
	// Annotations
}

// ReportBatch is a collection of reports that should all be processed together.
// We will create a new batch for each upload that the collector receives.
// Certain processors might join batches together or split them up.
type ReportBatch struct {
	Reports []NelReport

	// When this batch was received by the collector
	Time time.Time

	// The URL that was used to upload the report.
	CollectorURL url.URL

	// The IP address of the client that uploaded the batch of reports.  You can
	// typically assume that's the same IP address that was used for the original
	// requests.  The IP address will be encoded as a string; for example,
	// "192.0.2.1" or "2001:db8::2".
	ClientIP string

	// The user agent of the client that uploaded the batch of reports.
	ClientUserAgent string

	// The key-value pairs of the HTTP header that is received by the collector.
	// This can be used to get additional information. One example is to get the
	// remote address of the client when the collector runs behind a proxy.
	Header http.Header

	// An arbitrary set of extra data that you can attach to this batch of
	// reports.
	// Annotations
}

// Annotations lets you attach an arbitrary collection of extra data to each
// individual report, and to each report batch.  Each annotation has a name and
// an arbitrary type; it's up to you to make sure that your processors don't
// make conflicting assumptions about the type of an annotation with a
// particular name.
// type Annotations struct {
// 	Annotations map[string]interface{}
// }

/*
[
   {
      "age":9348,
      "body":{
         "elapsed_time":573,
         "method":"GET",
         "phase":"application",
         "protocol":"h2",
         "referrer":"https://tv3.darklibria.it/",
         "sampling_fraction":1.0,
         "server_ip":"142.132.244.104",
         "status_code":302,
         "type":"ok"
      },
      "type":"network-error",
      "url":"https://cache.libria.fun/videos/media/ts/8339/10/720/fff82.ts",
      "user_agent":"Mozilla/5.0 (Linux; arm_64; Android 10; HRY-LX1T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 YaApp_Android/23.36.1 YaSearchBrowser/23.36.1 BroPP/1.0 SA/3 Mobile Safari/537.36"
   }
]
*/
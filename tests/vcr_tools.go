package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

var HeadersToIgnore = []string{
	"Authorization",
	"X-Amz-Date",
	"X-Amz-Content-Sha256",
	"User-Agent",
	"Amz-Sdk-Request",
	"Amz-Sdk-Invocation-Id",
	"X-Amzn-Requestid",
	"Date",
	"X-Amz-Id-2",
	"X-Amz-Request-Id",
	"Content-Length",
}

func deepEqualContents(x, y any) bool {
	if reflect.ValueOf(x).IsNil() {
		if reflect.ValueOf(y).IsNil() {
			return true
		} else {
			return reflect.ValueOf(y).Len() == 0
		}
	} else {
		if reflect.ValueOf(y).IsNil() {
			return reflect.ValueOf(x).Len() == 0
		} else {
			return reflect.DeepEqual(x, y)
		}
	}
}

// TODO fix rest of matchers

func bodyMatches(r *http.Request, i cassette.Request) bool {
	if r.Body != nil {
		var buffer bytes.Buffer
		if _, err := buffer.ReadFrom(r.Body); err != nil {
			return false
		}

		r.Body = io.NopCloser(bytes.NewBuffer(buffer.Bytes()))
		if buffer.String() != i.Body {
			return false
		}
	} else {
		if len(i.Body) != 0 {
			return false
		}
	}

	return true
}

// modified version from default matcher
var CustomMatcher = func(r *http.Request, i cassette.Request) bool {
	if r.Method != i.Method {
		return false
	}

	if r.URL.String() != i.URL {
		return false
	}

	if r.Proto != i.Proto {
		return false
	}

	if r.ProtoMajor != i.ProtoMajor {
		return false
	}

	if r.ProtoMinor != i.ProtoMinor {
		return false
	}

	requestHeader := r.Header.Clone()
	cassetteRequestHeaders := i.Headers.Clone()

	for _, header := range HeadersToIgnore {
		delete(requestHeader, header)
		delete(cassetteRequestHeaders, header)
	}

	if !deepEqualContents(requestHeader, cassetteRequestHeaders) {
		return false
	}

	if !bodyMatches(r, i) {
		return false
	}

	if !deepEqualContents(r.TransferEncoding, i.TransferEncoding) {
		return false
	}

	if r.Host != i.Host {
		return false
	}

	// Only ParseForm for non-GET requests since that would use query params
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		err := r.ParseForm()
		if err != nil {
			return false
		}
	}
	if !deepEqualContents(r.Form, i.Form) {
		return false
	}

	if !deepEqualContents(r.Trailer, i.Trailer) {
		return false
	}

	if r.RemoteAddr != i.RemoteAddr {
		return false
	}

	if r.RequestURI != i.RequestURI {
		return false
	}

	return true
}

func prettyPrint(raw string) string {
	rawBytes := []byte(raw)

	if json.Valid(rawBytes) {
		var v interface{}
		if err := json.Unmarshal(rawBytes, &v); err == nil {
			if pretty, err := json.MarshalIndent(v, "", "  "); err == nil {
				return string(pretty)
			}
		}
	}

	return ""
}

func prettyFormatInteraction(i *cassette.Interaction) error {
	// comment out because this messages with matcher
	// if i.Request.Body != "" {
	// 	// Read and reformat
	// 	reqRaw := []byte(i.Request.Body)
	// 	i.Request.Body = prettyPrint(reqRaw)
	// }

	if i.Response.Body != "" {
		prettyResponseBody := prettyPrint(i.Response.Body)
		if prettyResponseBody != "" {
			i.Response.Body = prettyResponseBody
		}
	}

	return nil
}

var environmentIdRegex *regexp.Regexp = regexp.MustCompile(`e-[a-zA-Z0-9]{8,}`)
var instanceIdRegex *regexp.Regexp = regexp.MustCompile(`i-[a-zA-Z0-9]{8,}`)
var dbIdRegex *regexp.Regexp = regexp.MustCompile(`db-[a-zA-Z0-9]{8,}`)
var amiIdRegex *regexp.Regexp = regexp.MustCompile(`ami-[a-zA-Z0-9]{8,}`)
var vpcIdRegex *regexp.Regexp = regexp.MustCompile(`vpc-[a-zA-Z0-9]{8,}`)
var subnetIdRegex *regexp.Regexp = regexp.MustCompile(`subnet-[a-zA-Z0-9]{8,}`)
var namespaceRegex *regexp.Regexp = regexp.MustCompile(`ns-[a-zA-Z0-9]{8,}`)
var securityGroupIdRegex *regexp.Regexp = regexp.MustCompile(`sg-[a-zA-Z0-9]{8,}`)
var expandedSecurityGroupIdRegex *regexp.Regexp = regexp.MustCompile(`securitygroup-[a-zA-Z0-9]{8,}`)
var volumeIdRegex *regexp.Regexp = regexp.MustCompile(`vol-[a-zA-Z0-9]{8,}`)
var attachmentIdRegex *regexp.Regexp = regexp.MustCompile(`eni-attach-[a-zA-Z0-9]{8,}`)
var reservationIdRegex *regexp.Regexp = regexp.MustCompile(`r-[a-zA-Z0-9]{8,}`)

var accountIdInArn *regexp.Regexp = regexp.MustCompile(`:[0-9]{10,}:`)
var longNumberInXmlTag *regexp.Regexp = regexp.MustCompile(`>[0-9]{8,}<`) // we're going to assume that any numeric xml values are identifications of some sort, so just sanitize it
var uuidv2Regex *regexp.Regexp = regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
var iso8601Regex *regexp.Regexp = regexp.MustCompile(`\\d{4}-\\d\\d-\\d\\dT\\d\\d:\\d\\d:\\d\\d(\\.\\d+)?(([+-]\\d\\d:\\d\\d)|Z)?`)

// Updated idTagRegex for better content capture, case-insensitive, dot-all
var idTagRegex *regexp.Regexp = regexp.MustCompile(`(?is)<(id|DbiResourceId|HostedZoneId)>(.*?)</(?:id|DbiResourceId|HostedZoneId)>`)
var keyNameTagRegex *regexp.Regexp = regexp.MustCompile(`(?i)<(keyName)>.+</(keyName)>`)
var masterUsernameTagRegex *regexp.Regexp = regexp.MustCompile(`(?i)<(MasterUsername)>.+</(MasterUsername)>`)
var nextTokenTagRegex *regexp.Regexp = regexp.MustCompile(`(?i)<(NextToken)>.+</(NextToken)>`)

// Regex to detect path-like IDs and capture the prefix
var idPathPrefixRegex = regexp.MustCompile(`^(.*\/)[^/]+$`)

const sanitizedIDInPath = "000000"
const sanitizedIDDefault = "000000000000"

var ipv4Regex *regexp.Regexp = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
var macAddressRegex *regexp.Regexp = regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`)

var beanstalkSecurityGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBSecurityGroup-[0-9A-Z]{10,}`)
var beanstalkLoadBalancerSecurityGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBLoadBalancerSecurityGroup-[0-9A-Z]{10,}`)
var beanstalkAutoScalingGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBAutoScalingGroup-[0-9A-Z]{10,}`)

var amazonawsUrlRegex *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.[a-zA-Z0-9]+\.amazonaws\.com`)
var beanstalkUrlSubdomainRegex *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.elasticbeanstalk\.com`)
var internalUrlRegex *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.[a-zA-Z0-9]+\.internal`)

// Regex to find <Versions> blocks (case-insensitive, dot-all)
var versionsBlockRegex *regexp.Regexp = regexp.MustCompile(`(?is)<(Versions)>(.*?)</Versions>`)

// Regex to find <member> tags within a block (case-insensitive, dot-all)
var memberTagRegex *regexp.Regexp = regexp.MustCompile(`(?is)<(member)>(.*?)</member>`)

func sanitizeBody(body string) string {
	body = uuidv2Regex.ReplaceAllString(body, "00000000-0000-0000-0000-000000000000")
	body = environmentIdRegex.ReplaceAllString(body, "e-aaaaaaaaaa")
	body = instanceIdRegex.ReplaceAllString(body, "i-aaaaaaaaaa")
	body = dbIdRegex.ReplaceAllString(body, "db-AAAAAAAAAA")
	body = amiIdRegex.ReplaceAllString(body, "ami-aaaaaaaaaa")
	body = vpcIdRegex.ReplaceAllString(body, "vpc-aaaaaaaaaa")
	body = subnetIdRegex.ReplaceAllString(body, "subnet-aaaaaaaaaa")
	body = namespaceRegex.ReplaceAllString(body, "ns-aaaaaaaaaa")
	body = securityGroupIdRegex.ReplaceAllString(body, "sg-aaaaaaaaaa")
	body = expandedSecurityGroupIdRegex.ReplaceAllString(body, "securitygroup-aaaaaaaaaa")
	body = volumeIdRegex.ReplaceAllString(body, "vol-aaaaaaaaaa")
	body = attachmentIdRegex.ReplaceAllString(body, "eni-attach-aaaaaaaaaa")
	body = reservationIdRegex.ReplaceAllString(body, "r-aaaaaaaaaa")

	body = accountIdInArn.ReplaceAllString(body, ":0000000000:")
	body = longNumberInXmlTag.ReplaceAllString(body, ">00000000<")
	body = iso8601Regex.ReplaceAllString(body, "2020-01-01T00:00:00.000Z")

	// Conditional sanitization for ID-like tags
	body = idTagRegex.ReplaceAllStringFunc(body, func(match string) string {
		submatches := idTagRegex.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match // Should not happen, but safeguard
		}
		openingTagName := submatches[1] // e.g., "Id", "HostedZoneId"
		originalContent := submatches[2]

		var sanitizedValue string
		pathSubmatches := idPathPrefixRegex.FindStringSubmatch(originalContent)
		if pathSubmatches != nil {
			// Content is path-like, preserve prefix
			pathPrefix := pathSubmatches[1]
			sanitizedValue = pathPrefix + sanitizedIDInPath
		} else {
			// Content is not path-like, use default sanitization
			sanitizedValue = sanitizedIDDefault
		}
		return fmt.Sprintf("<%s>%s</%s>", openingTagName, sanitizedValue, openingTagName)
	})

	body = masterUsernameTagRegex.ReplaceAllString(body, "<$1>aaaaaaaaaaaa</$2>")
	body = keyNameTagRegex.ReplaceAllString(body, "<$1>aaaaaaaaaa</$2>")
	body = nextTokenTagRegex.ReplaceAllString(body, "<$1>BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB</$2>")

	body = ipv4Regex.ReplaceAllString(body, "0.0.0.0")
	body = macAddressRegex.ReplaceAllString(body, "00:00:00:00:00:00")

	body = beanstalkSecurityGroupNameRegex.ReplaceAllString(body, "AWSEBSecurityGroup-AAAAAAAAAAAA")
	body = beanstalkLoadBalancerSecurityGroupNameRegex.ReplaceAllString(body, "AWSEBLoadBalancerSecurityGroup-AAAAAAAAAAAA")
	body = beanstalkAutoScalingGroupNameRegex.ReplaceAllString(body, "AWSEBAutoScalingGroup-AAAAAAAAAAAA")

	body = amazonawsUrlRegex.ReplaceAllString(body, "subdomain.us-west-2.service.amazonaws.com")
	body = beanstalkUrlSubdomainRegex.ReplaceAllString(body, "subdomain.us-west-2.elasticbeanstalk.com")
	body = internalUrlRegex.ReplaceAllString(body, "subdomain.us-west-2.service.internal")

	// --- New Versions/Member Sanitization ---
	body = versionsBlockRegex.ReplaceAllStringFunc(body, func(versionsBlock string) string {
		// Extract the content between <Versions> and </Versions>
		matches := versionsBlockRegex.FindStringSubmatch(versionsBlock)
		if len(matches) < 3 {
			return versionsBlock // Should not happen with the regex, but safe guard
		}
		openingTag := matches[1] // e.g., "Versions" or "versions"
		content := matches[2]

		memberCounter := 0
		sanitizedContent := memberTagRegex.ReplaceAllStringFunc(content, func(memberBlock string) string {
			memberMatches := memberTagRegex.FindStringSubmatch(memberBlock)
			if len(memberMatches) < 3 {
				return memberBlock // Safeguard
			}
			memberOpeningTag := memberMatches[1] // e.g., "member" or "Member"

			// Generate the incremental value (repeat digit 8 times)
			digit := memberCounter % 10
			sanitizedValue := strings.Repeat(fmt.Sprintf("%d", digit), 8)
			memberCounter++

			// Reconstruct the member tag preserving case
			return fmt.Sprintf("<%s>%s</%s>", memberOpeningTag, sanitizedValue, memberOpeningTag)
		})

		// Reconstruct the Versions block preserving case
		return fmt.Sprintf("<%s>%s</%s>", openingTag, sanitizedContent, openingTag)
	})

	return body
}

func sanitizeAndFormatBodyHook(i *cassette.Interaction) error {
	if i.WasReplayed() {
		// don't reformat again if we are playing back
		return nil
	}

	for _, header := range HeadersToIgnore {
		delete(i.Request.Headers, header)
		delete(i.Response.Headers, header)
	}
	i.Request.ContentLength = 0

	i.Request.Body = sanitizeBody(i.Request.Body)
	i.Response.Body = sanitizeBody(i.Response.Body)

	if err := prettyFormatInteraction(i); err != nil {
		return err
	}

	if i.Request.Body != "" {
		parsedQuery, parseErr := url.ParseQuery(i.Request.Body)
		if parseErr == nil {
			changed := false
			for _, qpSanitizer := range requestQueryParamSanitizers {
				if value, ok := parsedQuery[qpSanitizer.ParamName]; ok {
					if len(value) > 0 && value[0] != "" {
						parsedQuery.Set(qpSanitizer.ParamName, qpSanitizer.SanitizedValue)
						changed = true
					}
				}
			}

			if changed {
				// Re-encode the modified query string and update Request.Body
				i.Request.Body = parsedQuery.Encode()

				// Update Request.Form as well for consistency
				if i.Request.Form == nil {
					i.Request.Form = make(url.Values)
				}
				for key, values := range parsedQuery {
					i.Request.Form[key] = values
				}
			}
		} else {
			fmt.Printf("Warning: Could not parse request body as query params: %v\n", parseErr)
		}
	}

	return nil
}

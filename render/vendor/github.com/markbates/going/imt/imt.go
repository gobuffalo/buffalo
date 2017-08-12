package imt

var Application = struct {
	Atom       string
	VND        string
	ECMAScript string
	JSON       string
	JavaScript string
	Binary     string
	Dash       string
	PDF        string
	PostScript string
	RDF        string
	RSS        string
	SOAP       string
	WOFF       string
	XHTML      string
	XML        string
	DTD        string
	XOP        string
	ZIP        string
	GZIP       string
}{
	Atom:       "application/atom+xml",
	VND:        "application/vnd.dart",
	ECMAScript: "application/ecmascript",
	JSON:       "application/json",
	JavaScript: "application/javascript",
	Binary:     "application/octet-stream",
	Dash:       "application/dash+xml",
	PDF:        "application/pdf",
	PostScript: "application/postscript",
	RDF:        "application/rdf+xml",
	RSS:        "application/rss+xml",
	SOAP:       "application/soap+xml",
	WOFF:       "application/font-woff",
	XHTML:      "application/xhtml+xml",
	XML:        "application/xml",
	DTD:        "application/xml-dtd",
	XOP:        "application/xop+xml",
	ZIP:        "application/zip",
	GZIP:       "application/gzip",
}

var Audio = struct {
	L24    string
	MP4    string
	MPEG   string
	Ogg    string
	FLAC   string
	Opus   string
	Vorbis string
	Wave   string
}{
	L24:    "audio/L24",
	MP4:    "audio/mp4",
	MPEG:   "audio/mpeg",
	Ogg:    "audio/ogg",
	FLAC:   "audio/flac",
	Opus:   "audio/opus",
	Vorbis: "audio/vorbis",
	Wave:   "audio/vnd.wave",
}

var Image = struct {
	GIF   string
	JPEG  string
	PJPEG string
	PNG   string
	BMP   string
	SVG   string
	TIFF  string
	DJVU  string
}{
	GIF:   "image/gif",
	JPEG:  "image/jpeg",
	PJPEG: "image/pjpeg",
	PNG:   "image/png",
	BMP:   "image/bmp",
	SVG:   "image/svg+xml",
	TIFF:  "image/tiff",
	DJVU:  "image/vnd.djvu",
}

var Multipart = struct {
	Mixed       string
	Alternative string
	Related     string
	FormData    string
	Signed      string
	Encrypted   string
}{
	Mixed:       "multipart/mixed",
	Alternative: "multipart/alternative",
	Related:     "multipart/related",
	FormData:    "multipart/form-data",
	Signed:      "multipart/signed",
	Encrypted:   "multipart/encrypted",
}

var Text = struct {
	CMD      string
	CSS      string
	CSV      string
	HTML     string
	Markdown string
	Plain    string
	RTF      string
	VCARD    string
}{
	CMD:      "text/cmd",
	CSS:      "text/css",
	CSV:      "text/csv",
	HTML:     "text/html",
	Markdown: "text/markdown",
	Plain:    "text/plain",
	RTF:      "text/rtf",
	VCARD:    "text/vcard",
}

var Video = struct {
	AVI       string
	MPEG      string
	MP4       string
	OGG       string
	Quicktime string
	WEBM      string
	Matroska  string
	WMV       string
	Flash     string
}{
	AVI:       "video/avi",
	MPEG:      "video/mpeg",
	MP4:       "video/mp4",
	OGG:       "video/ogg",
	Quicktime: "video/quicktime",
	WEBM:      "video/webm",
	Matroska:  "video/x-matroska",
	WMV:       "video/x-ms-wmv",
	Flash:     "video/x-flv",
}

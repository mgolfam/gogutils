# gogutils

[![Go Reference](https://pkg.go.dev/badge/github.com/mgolfam/gogutils.svg)](https://pkg.go.dev/github.com/mgolfam/gogutils)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-00ADD8?logo=go)](go.mod)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue.svg)](LICENSE)

A lightweight, dependency-light collection of Go utilities for everyday backend work — HTTP + cURL, file management, cryptography, compression, networking, JSON, struct mapping, time handling, and structured logging.

Every package is self-contained: import only what you need.

## Installation

```bash
go get github.com/mgolfam/gogutils
```

```go
import "github.com/mgolfam/gogutils/utils"
```

Requires **Go 1.21+**.

## Packages at a glance

| Package | Import path | Purpose |
| --- | --- | --- |
| `utils` | `.../utils` | Crypto helpers, networking, UUIDs, JSON, struct mapping, time, strings |
| `utils/compression` | `.../utils/compression` | Gzip & Deflate compress/decompress |
| `httpclient` | `.../httpclient` | HTTP client with caching, proxy, multipart, SOAP, and a cURL parser |
| `filemanager` | `.../filemanager` | Read/write files, checksums, Base64, MIME detection, multipart uploads |
| `crypt` | `.../crypt` | SHA-256 and Base36 encoding |
| `config` | `.../config` | Load a JSON config file into a struct |
| `glog` | `.../glog` | Minimal leveled logger |
| `services` | `.../services` | IP geolocation lookup |
| `dto` | `.../dto` | Shared data-transfer types (`IpInfo`, `LangConf`) |
| `enums`, `enums/doctype` | `.../enums`, `.../enums/doctype` | Password-type and document-type constants |

---

## Usage

### `utils` — crypto, hashing & random

```go
import "github.com/mgolfam/gogutils/utils"

utils.HashSha256("hello")           // hex SHA-256
utils.HashMd5("hello")              // hex MD5
utils.HashSha1("hello")             // hex SHA-1

utils.RandomString(16, true, true, true) // upper + lower + digits
utils.RandomInt(100)                      // 0..99

utils.ConvertTo36Base(123456)             // int64 -> base36 string
utils.ConvertToBase10From36("2n9c")       // base36 string -> int64
```

### `utils` — networking

```go
ips, _  := utils.GetLocalIPAddresses()
pub, _  := utils.GetPublicIPAddress()
mx, _   := utils.ResolveMX("example.com")
addrs, _:= utils.ResolveDNS("example.com")
open    := utils.IsPortOpen("127.0.0.1", 6379)
mac, _  := utils.GetMACAddress("eth0")
```

### `utils` — UUIDs & JSON

```go
utils.UUID()          // standard UUID v4
utils.CleanUUID()     // UUID without dashes

utils.PrintJSON(v)               // pretty-print any value
s := utils.ToJsonString(v)       // marshal to string
flat := utils.FlattenJSON(nested) // flatten nested map with dotted keys
```

### `utils` — struct mapping (reflection)

Copy fields between structs by matching field names/types, or via a `mapby` tag.

```go
type Src struct{ FullName string `mapby:"name"` }
type Dst struct{ Name string }

var dst Dst
utils.FillFromMapByTags(&dst, Src{FullName: "Ada"})   // maps by `mapby` tag
utils.FillByFieldNameAndType(&dst, src)               // maps by matching name+type
utils.FillByFieldNameAndTypeSlice(&dstSlice, srcSlice)
```

### `utils` — time

```go
utils.NowUnixSeconds()                      // current Unix time (s)
utils.Today()                               // "2006-01-02"
utils.Now(utils.TIME_FORMAT_TS)             // formatted now
utils.CalculateAge("1990-05-01", "2006-01-02")
utils.Time2Unix("2024-01-01 00:00:00", utils.TIME_FORMAT_TS)
```

### `utils/compression`

The package name is `utils`; alias it on import to avoid clashing with the main `utils` package.

```go
import compression "github.com/mgolfam/gogutils/utils/compression"

gz, _   := compression.Gzip(data)
raw, _  := compression.Gunzip(gz)

df, _   := compression.Deflate(data)
raw2, _ := compression.Inflate(df)
```

### `httpclient` — requests with caching & proxy

```go
import "github.com/mgolfam/gogutils/httpclient"

resp, err := httpclient.SendRequest(httpclient.HttpConfig{
    Method:  "GET",
    URL:     "https://api.example.com/data",
    Headers: map[string]string{"Accept": "application/json"},
    Timeout: 10 * time.Second,
    Cache:   true,   // write response to cache
    CacheTtl: 300,   // seconds
    UseProxy: false,
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(resp.StatusCode, string(resp.Body), resp.FromCache)
```

Also available: `SendMultipartFormData(FormDataConfig)`, `SoapCall(SoapConfig)`, and `Download(url, headers, filePath, userAgent)`.

### `httpclient` — cURL parser

Turn a `curl` command into a ready-to-send `HttpConfig`.

```go
cfg, err := httpclient.ParseCurlCommand(`curl -X POST https://api.example.com/login \
    -H "Content-Type: application/json" \
    -d '{"user":"ada"}'`)
if err != nil {
    log.Fatal(err)
}
resp, _ := httpclient.SendRequest(*cfg)
```

### `filemanager`

```go
import "github.com/mgolfam/gogutils/filemanager"

data, _ := filemanager.ReadFileBytes("in.dat")
_ = filemanager.WriteFileBytes("out.dat", data, 0o644)

sum, _  := filemanager.CalculateFileChecksum("out.dat") // SHA-256
b64, _  := filemanager.File2Base64("photo.png")
_ = filemanager.Base642File(b64, "copy.png")

ext, _  := filemanager.FileExtensionFromBytes(data)     // sniff MIME/extension

// Streaming writer implementing io.WriteCloser
w, _ := filemanager.NewFileWriter("stream.log")
defer w.Close()
w.Write([]byte("line\n"))
```

### `crypt`

```go
import "github.com/mgolfam/gogutils/crypt"

crypt.Sha256("hello")        // hex SHA-256
crypt.EncodeBase36(123456)   // int64 -> base36
crypt.DecondeBase36("2n9c")  // base36 -> int64
```

### `config`

Load a JSON file into any struct and set the log level in one call.

```go
import "github.com/mgolfam/gogutils/config"

type AppConfig struct {
    Port int    `json:"port"`
    DB   string `json:"db"`
}

var cfg AppConfig
if err := config.LoadConfig("config.json", "INFO", &cfg); err != nil {
    log.Fatal(err)
}
```

### `glog`

```go
import "github.com/mgolfam/gogutils/glog"

glog.LogLevel.Label = glog.INFO   // DEBUG | INFO | WARN | ERROR
glog.LogLevel.Load()

glog.Log("server started")
glog.LogL(glog.ERROR, "failed:", err) // only prints when >= configured level
```

### `services` — IP geolocation

```go
import "github.com/mgolfam/gogutils/services"

info := services.GetIpInfo("8.8.8.8") // *dto.IpInfo
fmt.Println(info.Country, info.City, info.Isp)
```

### `enums`

```go
import (
    "github.com/mgolfam/gogutils/enums"
    "github.com/mgolfam/gogutils/enums/doctype"
    "github.com/mgolfam/gogutils/utils"
)

// Password types: enums.PASS_TYPE_PLAIN | PASS_TYPE_MD5 | PASS_TYPE_SHA256
ok := utils.IsPasswordEqual(dbHash, input, enums.PASS_TYPE_SHA256)

// Document types: doctype.PersonalPhoto, doctype.Password, doctype.ID_Card, doctype.Shenasname
```

---

## Project layout

```text
gogutils/
├── config/       JSON config loader
├── crypt/        SHA-256, Base36
├── dto/          shared data-transfer types
├── enums/        password- and document-type constants
├── filemanager/  file I/O, checksums, Base64, MIME, multipart
├── glog/         leveled logger
├── httpclient/   HTTP client, cache, proxy, multipart, SOAP, cURL parser
├── services/     IP geolocation
└── utils/        crypto, networking, UUID, JSON, mapping, time, strings
    └── compression/  gzip & deflate
```

## Testing

Run the full suite:

```bash
go test ./...            # all packages
go test -race ./...      # with the race detector
go test -cover ./...     # with coverage
```

Or via the `Makefile` (run `make help` to list all targets):

```bash
make test     # run all unit tests
make race     # tests with the race detector
make cover    # tests with a coverage summary
make check    # fmt + vet + test
```

The deterministic packages (`crypt`, `dto`, `utils`, `utils/compression`, `filemanager`, `httpclient`'s cURL parser) ship with unit tests. Network- and I/O-dependent helpers (`utils/network`, `services`, live HTTP requests) are intentionally left to integration testing.

## Contributing

Contributions are welcome. Please keep packages focused and dependency-light, run `go vet ./...` and `go test ./...`, and format with `gofmt` before opening a pull request.

## License

Released under the [BSD 3-Clause License](LICENSE). © 2024 Milad Golfam.

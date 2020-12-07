Golden is a small library to help with testing using golden files. Its main
purpose (but not only one)  is to provide an easy way to describe 
HTTP request / response as YAML files. 

[![Go Report Card](https://goreportcard.com/badge/github.com/rzajac/golden)](https://goreportcard.com/report/github.com/rzajac/golden)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/rzajac/golden)

# Installation 

```
go get github.com/rzajac/golden
```

# Usage

## Asserting

Lets say we have a golden file looking like this:

```yaml
bodyType: json
body: |
    { "key1": "val1" }
```

How to use it in test.

```go
func Test_Assert(t *testing.T) {
    // --- Given --- 
    gld := golden.File(golden.Open(t, "testdata/file.yaml", nil))
    
    // --- When ---
    data := []byte(`{
        "key1": "val1"
    }`)
    
    // --- Then --- 
    gld.Assert(data)
}
```

Because the `bodyType` was set to `json` the `data` in the test doesn't 
have to be formatted exactly the same way as it's in the golden file. The
library is smart enough to compare data represented as JSON not the strings.  

If you need exact match set `bodyType` to `text`.  

## Unmarshalling

```go
type Data struct {
    Key1 string `json:"key1"`
}

func Test_Unmarshal(t *testing.T) {
    // --- Given ---
    gld := golden.File(golden.Open(t, "../testdata/file.yaml", nil))

    // --- When ---
    data := &Data{}
    gld.Unmarshall(data)

    // --- Then ---
    if data.Key1 != "val1" {
        t.Errorf("expected `%s` got `%s`", "val1", data.Key1)
    }
}
```

In this case golden file body will be unmarshalled (using `json.Unmarshal`)
to structure `Data`. Any errors during unmarshalling will be handled by 
`Unmarshall` method.

## Testing HTTP request / response

Golden file describing the HTTP request and response:

```yaml
request:
    method: POST
    path: /some/path
    query: key0=val0&key1=val1
    headers:
        - 'Authorization: Bearer token'
        - 'Content-Type: application/json'
    bodyType: json
    body: |
        {
          "key2": "val2"
        }

response:
    statusCode: 200
    headers:
        - 'Content-Type: application/json'
    bodyType: json
    body: |
        { "success": true }
```

Example test using golden file:

```go
func Test_Endpoint(t *testing.T) {
    // --- Given ---
    pth := "testdata/request.yaml"
    gld := golden.Exchange(golden.Open(t, pth, nil))

    // Setup mocks.
    srvH, mckS := SrvMock()
    mckS.On("CheckUserAccess", "token").Return(true, nil)

    // Prepare request recorder.
    rec := httptest.NewRecorder()

    // --- When ---
    srvH.ServeHTTP(rec, gld.Request.Request())

    // --- Then ---
    mckS.AssertExpectations(t)
    gld.Response.Assert(rec.Result())
}
```

## Golden files as templates

Golden files can also be used as Go templates when more dynamic approach 
is needed.

```yaml
request:
    method: POST
    path: /some/path
    query: key0=val0&key1=val1
    headers:
        - 'Authorization: Bearer {{ .token }}'
        - 'Content-Type: application/json'
    bodyType: json
    body: |
        {
          "key2": "val2"
        }

response:
    statusCode: 200
    headers:
        - 'Content-Type: application/json'
    bodyType: json
    body: |
        { "success": true }
```

```go
func Test_Endpoint(t *testing.T) {
    // --- Given ---
    token := GetTestToken()
    tplD := make(golden.Map).Add("token", token)
    tpl := "testdata/request.yaml"
    gld := golden.Exchange(golden.Open(t, tpl, tplD))

    // Setup mocks.
    srvH, mckS := SrvMock()
    mckS.On("CheckUserAccess", token).Return(true, nil)

    // Prepare request recorder.
    rec := httptest.NewRecorder()

    // --- When ---
    srvH.ServeHTTP(rec, gld.Request.Request())

    // --- Then ---
    mckS.AssertExpectations(t)
    gld.Response.Assert(rec.Result())
}
```

Check out the documentation to see full API.

## License

Apache License, Version 2.0

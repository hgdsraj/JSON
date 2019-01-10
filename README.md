# JSON

JSON allows Java like syntax for fast JSON parsing with low memory usage.

### Usage

Get the package using your favorite dependency manager or `go get github.com/josuemnb/JSON` to install to your $GOPATH.

Example

    import (
        "github.com/josuemnb/JSON"
    )

    json := JSON.Load(s)
    if o,ok := json.GetObject("j"); ok {
      println("J",j.GetString("ok"))
    }
    if rows, ok := json.GetArray("rows"); ok {
        for rows.Next() {
          row := rows.Current()
          str := row.GetString("mystring")
          f := row.GetFloat("float")
        }
    }

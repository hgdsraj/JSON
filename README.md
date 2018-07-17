# JSON

Fast Json Parsing

Low memory usage

Syntax like Java Json Object

Example

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

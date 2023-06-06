# `enc` WIP

Replacement for `encoding/json`, providing an intermediate layer of abstraction between the encoded data and the typed data.

## Interstitial

```go
nodes := []enc.Node{
  enc.Nil{},
  enc.String("foo bar"),
  enc.Number(3.14),
  enc.Map{"null":enc.Nil{}},
  enc.List{enc.String("first")},
}
```

Rationale: in the built in `encoding/json` library, everything gets converted between `[]byte` and specific types.

This means that when customization is needed, you end up implementing `json.Unmarshaler` or `json.RawMessage`.

This likely requires multiple scans to the data, and writing convoluted code or simply more complex.

Using an intermediate layer, which maps the primitive types in the JSON format, allows a single scan of the `[]byte` data, 
and only high level computation on it.

Those primitive types can be then used when custom marshaling or unmarshaling.

## enc.Pairs

This is a special type that can Marshal itself as a JSON object, but is implemented as a list of pairs, which then guarantee the order.

To keep the usability of this library high, we opted to avoid OrderedMap which are clumsy to use, and instead allow you to choose between the 
fast `enc.Map`, or the ordered `enc.Pairs`.

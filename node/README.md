# Node

Nodes are the basic building blocks for the Halation runtime.

Nodes are struct that meet the `node.Node` interface, which can usually be accomplished by including `node.BaseNode` and implementing `TypeName()`.

Additionally, Nodes may have fields that are set by the user either at runtime or configuration time.
Take the following timer node as an example:

```go
type Timer struct {
  Description string
  Promises []node.Promise

  Delay time.Duration
}
```

When the timer is be configured with a set delay during configuration time (e.g. `{ delay: "50s" }`), the `Delay` field is set to a value.
However, when the delay field is set at runtime, a struct in the `Promises` slice contains the instructions to do so:
```go
Timer.Promises = []node.Promise{
  {
    FieldName: "Delay",
    SupplyNodeName: node.NodeName{"nyiyui.ca/halation/...", "calculate-delay-time"},
  },
}
```

To set an order of operations, set FieldName to `"activate"`:
```go
Timer.Promises = []node.Promise{
  {
    FieldName: "activate",
    SupplyNodeName: node.NodeName{"nyiyui.ca/halation/...", "0-pre-show"},
  },
}
```

Field names may introspect:
- `State.Channels[0].ChannelID`
- `State.ExtraProperties["paused"]`

## Controlling which field can be set at runtime

Field that can be set at runtime must have a blank `halation:""` tag (`halation:"abc"` is illegal).

## TODOs

- test nodes with complex dependency graphs (only DAGs for now)
```
  *--> B --*
 /          \
A ---> C --> E
 \
  *--> D
```
(NOTE: E only activates when both B and C's promises are fulfilled.)

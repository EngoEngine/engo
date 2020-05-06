# Parenting Demo

Each button is connected by a line. The button above is a "parent" to the buttons
below. If it is off, the children will NOT turn on. When a parent is turned on, all
children will also turn on. When a parent is turned off, all children will also
turn off. Children can be toggled without affecting the parent.

## What does it do?
* It shows how to implement a System that utilizes the parent feature of the ecs.BasicEntity

## What are important aspects of the code?

### Use AppendChild to create parent-child relationships
rectBasics is a `[]ecs.BasicEntity`
```go
rectBasics[0].AppendChild(&rectBasics[1])
rectBasics[0].AppendChild(&rectBasics[2])
rectBasics[0].AppendChild(&rectBasics[3])

rectBasics[1].AppendChild(&rectBasics[4])
rectBasics[1].AppendChild(&rectBasics[5])

rectBasics[2].AppendChild(&rectBasics[6])

rectBasics[3].AppendChild(&rectBasics[7])
```

### OnClickSystem

The OnClickSystem in this demo has a map of the entities and their associated
IDs. This allows the update to easily check if the entity's parents and Children
are actually part of the system, so the system can update it.

The Update is able to get the children, as well as all the children-of-Children
by using `ent.Descendents()` and just the children wwith `ent.Children()`. It
is able to access the parent with `ent.Parent()`.

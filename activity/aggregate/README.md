---
title: Aggregate
weight: 4603
---

# Aggregate
This activity allows you to aggregate data and calculate an average or sliding average.


## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/aggregate
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "function",
      "type": "string",
      "required": true,
      "allowed" : ["avg","sum","min","max","count"]
    },
    {
      "name": "windowSize",
      "type": "integer",
      "required": true,
      "allowed" : ["avg","sum","min","max","count"]
    },
    {
      "name": "value",
      "type": "number"
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "number"
    },
    {
      "name": "report",
      "type": "boolean"
    }
  ]
}
```

## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| function    | True     | The aggregate fuction. Supported: avg (Default),sum,min,max,count |
| windowType  | True     | The window type of the values to aggregate Supported: tumbling (Default), sliding, timeTumbling, timeSliding |
| windowSize  | True     | The window size of the values to aggregate |
| value       | False    | The value to aggregate |


## Example
The below example aggregates a 'temperature' attribute with a moving window of size 5:

```json
"id": "aggregate_4",
"name": "Aggregate",
"description": "Simple Aggregator Activity",
"activity": {
  "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/aggregate",
  "input": {
    "function": "avg",
    "windowType": "tumbling",
    "windowSize": "5"
  },
  "mappings": {
    "input": [
      {
        "type": "assign",
        "value": "temperature",
        "mapTo": "value"
      }
    ]
  }
```

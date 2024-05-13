package server

// run flags
const AFlag = "a"
const ADefault = "localhost:8080"
const AUsage = "specify the url"

const IFlag = "i"
const IDefault = 3
const IUsage = "interval to save storage"

const FFlag = "f"
const FDefault = "tmp/metrics-db.json"
const FUsage = "path to save snapshot"

const RFlag = "r"
const RDefault = true
const RUsage = "restore start snapshot?"

// routing paths
const MainPath = "/"
const GetMetricByURLPath = "/value/:metricType/:metricName"
const GetMetricByJSONPath = "/value"
const UpdateByURLPath = "/update/:metricType/:metricName/:metricValue"
const UpdateByJSONPath = "/update"

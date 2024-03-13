package server

// run flags
const AFlag = "a"
const ADefault = "localhost:8080"
const AUsage = "specify the url"

// routing paths
const MainPath = "/"
const GetMetricByURLPath = "/value/:metricType/:metricName"
const GetMetricByJSONPath = "/value"
const UpdateByURLPath = "/update/:metricType/:metricName/:metricValue"
const UpdateByJSONPath = "/update"

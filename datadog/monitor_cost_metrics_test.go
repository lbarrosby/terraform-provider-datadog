package datadog

import (
	"encoding/json"
	"testing"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestMonitorCostMetricsRoundTrip(t *testing.T) {
	variables := getMonitorFormulaQuerySchema().Elem.(*schema.Resource)
	cost := variables.Schema["cloud_cost_query"].Elem.(*schema.Resource)
	require.Empty(t, cost.Schema["data_source"].ValidateDiagFunc("metrics", nil))
	require.NotEmpty(t, cost.Schema["data_source"].ValidateDiagFunc("invalid", nil))
	for _, aggregator := range []string{"avg", ""} {
		input := map[string]interface{}{"data_source": "metrics", "name": "query1", "query": "avg:system.cpu.user{*}", "aggregator": aggregator}
		query := buildMonitorFormulaAndFunctionCloudCostQuery(input)
		require.NotNil(t, query.MonitorFormulaAndFunctionMetricQueryDefinition)
		encoded, err := json.Marshal(query)
		require.NoError(t, err)
		var decoded datadogV1.MonitorFormulaAndFunctionQueryDefinition
		require.NoError(t, json.Unmarshal(encoded, &decoded))
		output := buildTerraformCostMonitorVariables([]datadogV1.MonitorFormulaAndFunctionQueryDefinition{decoded})[0]["cloud_cost_query"].([]map[string]interface{})[0]
		require.Equal(t, input["data_source"], output["data_source"])
		require.Equal(t, input["name"], output["name"])
		require.Equal(t, input["query"], output["query"])
		if aggregator != "" {
			require.Equal(t, datadogV1.MonitorFormulaAndFunctionCostAggregator(aggregator), *output["aggregator"].(*datadogV1.MonitorFormulaAndFunctionCostAggregator))
		} else {
			require.NotContains(t, output, "aggregator")
		}
	}
}

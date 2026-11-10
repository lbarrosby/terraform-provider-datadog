package fwprovider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestMonitorCostMetricsRoundTrip(t *testing.T) {
	require.Contains(t, (&monitorResource{}).getAllowCloudCostDataSource(), "metrics")
	r := &monitorResource{}
	for _, aggregator := range []types.String{types.StringValue("avg"), types.StringNull()} {
		input := CloudCostQuery{DataSource: types.StringValue("metrics"), Name: types.StringValue("query1"), Query: types.StringValue("avg:system.cpu.user{*}"), Aggregator: aggregator}
		queries := r.buildCloudCostQueryStruct([]CloudCostQuery{input})
		require.NotNil(t, queries[0].MonitorFormulaAndFunctionMetricQueryDefinition)
		encoded, err := json.Marshal(queries)
		require.NoError(t, err)
		var decoded []datadogV1.MonitorFormulaAndFunctionQueryDefinition
		require.NoError(t, json.Unmarshal(encoded, &decoded))
		options := datadogV1.MonitorOptions{}
		options.SetVariables(decoded)
		state := &monitorResourceModel{}
		r.updateVariablesState(context.Background(), state, &options)
		require.Equal(t, []CloudCostQuery{input}, state.Variables[0].CloudCostQuery)
	}
}

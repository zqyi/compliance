package output

import (
	"encoding/json"
	"fmt"

	"github.com/prometheus/compliance/promql/comparer"
	"github.com/prometheus/compliance/promql/config"
)

// JSON produces JSON-based output for a number of query results.
func JSON(results []*comparer.Result, includePassing bool, cfg *config.Config) {
	tweaks := cfg.QueryTweaks
	buf, err := json.Marshal(map[string]interface{}{
		"totalResults":   len(results), // Needed because we may exclude passing results.
		"results":        results,
		"includePassing": includePassing,
		"queryTweaks":    tweaks,
	})
	if err != nil {
		panic(err)
	}
	fmt.Print(string(buf))
}

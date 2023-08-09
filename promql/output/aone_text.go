package output

import (
	"fmt"
	"strings"

	"github.com/prometheus/compliance/promql/comparer"
	"github.com/prometheus/compliance/promql/config"
)

// Text produces text-based output for a number of query results.
func AoneText(results []*comparer.Result, includePassing bool, cfg *config.Config) {
	fmt.Println(strings.Repeat("=", 80))
	var test_type string
	parallel_mode := cfg.TestTargetConfig.Headers["x-sls-parallel-enable"]
	pushdown_mode := cfg.TestTargetConfig.Headers["x-sls-pushdown-enable"]
	// fmt.Println("parallel_mode ", parallel_mode)
	// fmt.Println("pushdown_mode ", pushdown_mode)

	if parallel_mode != "true" && pushdown_mode != "true" {
		test_type = "Original"
	} else if parallel_mode == "true" && pushdown_mode != "true" {
		test_type = "Parallel"
	} else if parallel_mode != "true" && pushdown_mode == "true" {
		test_type = "Pushdown"
	} else if parallel_mode == "true" && pushdown_mode == "true" {
		test_type = "Parallel_Pushdown"
	}

	fmt.Printf("Query Mode: %s\n", test_type)

	tweaks := cfg.QueryTweaks
	successes := 0
	unsupported := 0
	for _, res := range results {
		if res.Success() {
			successes++
			if !includePassing {
				continue
			}
		}
		if res.Unsupported {
			unsupported++
		}

		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("QUERY: %v\n", res.TestCase.Query)
		fmt.Printf("START: %v, STOP: %v, STEP: %v\n", res.TestCase.Start, res.TestCase.End, res.TestCase.Resolution)
		fmt.Printf("RESULT: ")
		if res.Success() {
			fmt.Println("PASSED")
		} else if res.Unsupported {
			fmt.Println("UNSUPPORTED: ")
			fmt.Printf("Query is unsupported: %v\n", res.UnexpectedFailure)
		} else {
			fmt.Printf("FAILED: ")
			if res.UnexpectedFailure != "" {
				fmt.Printf("Query failed unexpectedly: %v\n", res.UnexpectedFailure)
			}
			if res.UnexpectedSuccess {
				fmt.Println("Query succeeded, but should have failed.")
			}
			if res.Diff != "" {
				fmt.Println("Query returned different results:")
				fmt.Println(res.Diff)
			}
		}
	}

	fmt.Println("General query tweaks:")
	if len(tweaks) == 0 {
		fmt.Println("None.")
	}
	for _, t := range tweaks {
		fmt.Println("* ", t.Note)
	}

	total := len(results)
	fails := total - successes

	fmt.Printf("CUSTOM_NAME_%s: %s\n", test_type, test_type)
	fmt.Printf("CUSTOM_TITLE_%s: %d/%d\n", test_type, successes, total)

	if fails > 0 {
		fmt.Printf("CUSTOM_VAL_%s: 失败 %d\n", test_type, fails)
		fmt.Printf("CUSTOM_CLASS_%s: text-danger\n", test_type)
		fmt.Print("CUSTOM_STATUS: 301\n")
	} else {
		fmt.Printf("CUSTOM_VAL_%s: 通过\n", test_type)
		fmt.Printf("CUSTOM_CLASS_%s: text-success\n", test_type)
	}

}

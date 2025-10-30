package main

import (
	"armur-codescanner/internal/tools"
	"fmt"
)

func main() {
	dirPath := "/home/wharpheus/tools/domain-ninja/dnc-foundry/src"

	fmt.Println("🚀 DNC Foundry Solidity Optimization Analysis")
	fmt.Println("==========================================")

	fmt.Println("\n🔍 Running Gas Optimizer...")
	gasResults := tools.RunGasOptimizer(dirPath)
	fmt.Println("Gas Optimization Results:")
	printResults(gasResults)

	fmt.Println("\n🔗 Running LP Pairing Checks...")
	lpResults := tools.RunLPPairingChecks(dirPath)
	fmt.Println("LP Pairing Check Results:")
	printResults(lpResults)

	fmt.Println("\n🔒 Running Security Scans...")
	securityResults := tools.RunSlither(dirPath)
	fmt.Println("Security Results (Reentrancy & Access Control):")
	printResults(securityResults)

	fmt.Println("\n💰 Running DeFi Optimizations...")
	defiResults := tools.RunDeFiOptimizations(dirPath)
	fmt.Println("DeFi Optimization Results:")
	printResults(defiResults)

	fmt.Println("\n🛡️ HARDENED SECURITY & OPTIMIZATION ANALYSIS COMPLETE!")
	fmt.Println("==========================================")
}

func printResults(results map[string]interface{}) {
	if results == nil {
		fmt.Println("  - No results")
		return
	}

	for category, data := range results {
		fmt.Printf("  %s:\n", category)
		switch v := data.(type) {
		case []interface{}:
			for _, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					fmt.Printf("    - %v\n", itemMap)
				} else {
					fmt.Printf("    - %v\n", item)
				}
			}
		default:
			fmt.Printf("    - %v\n", v)
		}
	}
}

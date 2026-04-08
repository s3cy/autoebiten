// Package autoui provides EbitenUI automation helpers for autoebiten.
// It enables CLI-based widget tree inspection, search, method invocation,
// and visual debugging for LLM-assisted E2E testing.
//
// Usage:
//
//	ui := ebitenui.UI{Container: root}
//	autoui.Register(&ui)
//
// Commands registered:
//   - autoui.tree: Export full widget tree as XML
//   - autoui.at: Get widget at (x,y) coordinates
//   - autoui.find: Simple key=value attribute search
//   - autoui.xpath: XPath 1.0 queries
//   - autoui.call: Reflection-based method invocation
//   - autoui.highlight: Visual debugging highlights
package autoui
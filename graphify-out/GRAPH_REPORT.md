# Graph Report - .  (2026-04-17)

## Corpus Check
- 174 files · ~174,627 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1470 nodes · 3602 edges · 38 communities detected
- Extraction: 49% EXTRACTED · 50% INFERRED · 0% AMBIGUOUS · INFERRED: 1786 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Input Capture|Input Capture]]
- [[_COMMUNITY_Custom Commands|Custom Commands]]
- [[_COMMUNITY_Command Execution|Command Execution]]
- [[_COMMUNITY_autoui Features|autoui Features]]
- [[_COMMUNITY_autoui Finder|autoui Finder]]
- [[_COMMUNITY_LaunchMode|Launch/Mode]]
- [[_COMMUNITY_autoui Core|autoui Core]]
- [[_COMMUNITY_RegistryCaller|Registry/Caller]]
- [[_COMMUNITY_Custom Handlers|Custom Handlers]]
- [[_COMMUNITY_Testkit Config|Testkit Config]]
- [[_COMMUNITY_RadioGroup|RadioGroup]]
- [[_COMMUNITY_ReflectionTabBook|Reflection/TabBook]]
- [[_COMMUNITY_Doc Templates|Doc Templates]]
- [[_COMMUNITY_Script Executor|Script Executor]]
- [[_COMMUNITY_Script AST|Script AST]]
- [[_COMMUNITY_State Export Tests|State Export Tests]]
- [[_COMMUNITY_Widget State|Widget State]]
- [[_COMMUNITY_Custom Data|Custom Data]]
- [[_COMMUNITY_Highlight|Highlight]]
- [[_COMMUNITY_RPC Messages|RPC Messages]]
- [[_COMMUNITY_Launch Socket Tests|Launch Socket Tests]]
- [[_COMMUNITY_AfterDraw|AfterDraw]]
- [[_COMMUNITY_Doc Rewrite|Doc Rewrite]]
- [[_COMMUNITY_Script Parser|Script Parser]]
- [[_COMMUNITY_OutputCR|Output/CR]]
- [[_COMMUNITY_Config Tests|Config Tests]]
- [[_COMMUNITY_Widget State Extract|Widget State Extract]]
- [[_COMMUNITY_Messages Tests|Messages Tests]]
- [[_COMMUNITY_Custom|Custom]]
- [[_COMMUNITY_Doc|Doc]]
- [[_COMMUNITY_Errors|Errors]]
- [[_COMMUNITY_Caller Export|Caller Export]]
- [[_COMMUNITY_autoui Doc|autoui Doc]]
- [[_COMMUNITY_testkit Doc|testkit Doc]]
- [[_COMMUNITY_Doc Verify|Doc Verify]]
- [[_COMMUNITY_Socket Env|Socket Env]]
- [[_COMMUNITY_Go Code Extract|Go Code Extract]]
- [[_COMMUNITY_Troubleshooting|Troubleshooting]]

## God Nodes (most connected - your core abstractions)
1. `contains()` - 82 edges
2. `WalkTree()` - 49 edges
3. `createTestNineSlice()` - 45 edges
4. `ExtractWidgetInfo()` - 43 edges
5. `NewContext()` - 35 edges
6. `InvokeMethod()` - 35 edges
7. `Documentation Template System Rewrite` - 35 edges
8. `NewCommandExecutor()` - 31 edges
9. `main()` - 28 edges
10. `ProcessRequest()` - 25 edges

## Surprising Connections (you probably didn't know these)
- `Graph Report` --semantically_similar_to--> `Mock`  [INFERRED] [semantically similar]
  graphify-out/GRAPH_REPORT.md → testkit/mock.go
- `Documentation Template System Rewrite` --conceptually_related_to--> `Config`  [EXTRACTED]
  docs/superpowers/specs/2026-04-10-doc-template-system-design.md → internal/docgen/config.go
- `Testkit Reference` --conceptually_related_to--> `Mock`  [EXTRACTED]
  skills/using-autoebiten/references/testkit.md → testkit/mock.go
- `TestParsePath()` --calls--> `parsePath()`  [INFERRED]
  state_exporter_test.go → state_exporter.go
- `Capture()` --calls--> `Capture()`  [INFERRED]
  autoebiten_default.go → integrate/integrate.go

## Hyperedges (group relationships)
- **** — patch_method, ebiten_v2_9_9, git_apply, replace_directive, go_mod [INFERRED]
- **** — library_method, update_func, capture_func, setmode_func, register_func, import_autoebiten [INFERRED]
- **** — cli, json_rpc, game_socket, launch_socket, proxy [INFERRED]
- **** — autoui, widgetinfo, ae_tag, xpath_query, proxy_handler, tree_cmd, find_cmd, call_cmd [INFERRED]
- **** — testkit, black_box_mode, white_box_mode, launch_func, mock, state_query [INFERRED]
- **** — launch_cmd, proxy, state_machine, carriage_return_writer, output_manager, err_game_not_connected [INFERRED]
- **** — docgen, docverify, make_docs, template_system, inline_config [INFERRED]
- **** — setmode_func, injection_only_mode, injection_fallback_mode, passthrough_mode, injected_input, real_input [INFERRED]
- **** — radiogroup, tabbook, list, reflection_utils, proxy_handler [INFERRED]
- **** — widgetinfo, state_field, custom_data, snapshot_tree, flat_output, slice_flattening [INFERRED]
- **autoui Architecture Components** — autoui_design, walktree, extractwidgetinfo, extractwidgetstate, extractcustomdata, invoke_method, widgetinfo, widgetnode [INFERRED 0.90]
- **Docgen Template Function Map** — doc_template_system, funcmap, docgen_config, docgen_launch_game, docgen_command, docgen_gocode, docgen_verify, docgen_normalize [INFERRED 0.90]
- **Launch Socket Crash Diagnostics** — launch_socket_crash, launch_socket, state_machine, unified_handler, crash_diagnostics, log_diff, error_accumulation [INFERRED 0.85]

## Communities

### Community 0 - "Input Capture"
Cohesion: 0.03
Nodes (83): Capture(), CursorPosition(), IsKeyJustPressed(), IsKeyJustReleased(), IsKeyPressed(), IsMouseButtonJustPressed(), IsMouseButtonJustReleased(), IsMouseButtonPressed() (+75 more)

### Community 1 - "Custom Commands"
Cohesion: 0.04
Nodes (113): RegisterOptions, TestContextAddOutput(), TestContextGetOutputsEmpty(), TestContextGetOutputsReturnsCopy(), TestContextSetConfig(), TestNewContext(), GetCustomCommand(), ListCustomCommands() (+105 more)

### Community 2 - "Command Execution"
Cohesion: 0.04
Nodes (86): CommandExecutor, Condition, waitLogger, Writer, detectProxyOrDirect(), EnsureTargetPID(), NewCommandExecutor(), captureOutput() (+78 more)

### Community 3 - "autoui Features"
Cohesion: 0.03
Nodes (127): _addr Attribute, ae_tag, any/interface{} Support, autoui.at, autoui Bug Fixes, autoui.call, autoui.call Extended Type Support, autoui Commands (+119 more)

### Community 4 - "autoui Finder"
Cohesion: 0.05
Nodes (84): CallRequest, CallResponse, CoordinateRequest, ExistsResponse, HighlightRequest, mockCommandContext, stringKeyValue, WidgetNode (+76 more)

### Community 5 - "Launch/Mode"
Cohesion: 0.05
Nodes (63): Mode, LaunchCommand, LaunchOptions, delayedMockClient, handleTestConnection(), mustMarshal(), newTestClient(), newTestHandler() (+55 more)

### Community 6 - "autoui Core"
Cohesion: 0.03
Nodes (42): at_cmd, autoebiten, autoui, Game, PlayerCard, call_cmd, Game, Game (+34 more)

### Community 7 - "Registry/Caller"
Cohesion: 0.06
Nodes (81): WidgetInfo, convertArg(), convertArgs(), InvokeMethod(), isWhitelistedSignature(), TestConvertArg_AnyParam(), TestConvertArg_IntToEnum(), TestConvertArg_NonEmptyInterfaceNotImplemented() (+73 more)

### Community 8 - "Custom Handlers"
Cohesion: 0.06
Nodes (41): testHandler, handleCallCommand, decodeParams(), ErrorResponse(), marshalResult(), ProcessRequest(), TestConcurrentProxyRequests(), TestProxyServerCleanup() (+33 more)

### Community 9 - "Testkit Config"
Cohesion: 0.07
Nodes (46): Config, black_box_mode, getStateExporterGameBinary(), TestEnemyStateQuery(), TestHealthModification(), TestPlayerMovement(), TestScreenshotCapture(), DefaultConfig() (+38 more)

### Community 10 - "RadioGroup"
Cohesion: 0.07
Nodes (44): RadioGroupElementInfo, RadioGroupHandler, schemaValidator, runSchemaCommand(), handleRadioGroupActiveIndex(), handleRadioGroupActiveLabel(), handleRadioGroupElements(), handleRadioGroupSetActiveByIndex() (+36 more)

### Community 11 - "Reflection/TabBook"
Cohesion: 0.08
Nodes (38): ProxyHandler, TabInfo, GetProxyHandler(), handleSetTabByIndex(), handleSetTabByLabel(), handleTabIndex(), handleTabLabel(), handleTabs() (+30 more)

### Community 12 - "Doc Templates"
Cohesion: 0.06
Nodes (42): AST Transforms, Command Execution, Crash Output, delay Function, dict Function, Doc Example Automation, Document Template Rewrite Workflow, Documentation Template System Rewrite (+34 more)

### Community 13 - "Script Executor"
Cohesion: 0.13
Nodes (23): TestScriptExecutorEndToEnd(), formatCustomCmd(), formatInputCmd(), formatMouseCmd(), formatScreenshotCmd(), formatWheelCmd(), NewExecutor(), TestExecutorBasic() (+15 more)

### Community 14 - "Script AST"
Cohesion: 0.07
Nodes (18): unmarshalCommand(), unmarshalRepeat(), Entry, CommandSchema, CommandWrapper, CustomCmd, DelayCmd, InputCmd (+10 more)

### Community 15 - "State Export Tests"
Cohesion: 0.12
Nodes (24): gameWithInterface, taggedGameState, taggedPlayer, TestStateExporterInterfaceField(), TestStateExporterInterfaceWithPointer(), TestStateExporterJSONTags(), TestStateExporterNilInterface(), getFieldByName() (+16 more)

### Community 16 - "Widget State"
Cohesion: 0.14
Nodes (21): extractButtonState(), extractCheckboxState(), extractComboButtonState(), extractLabelState(), extractListComboButtonState(), extractListState(), extractProgressBarState(), extractScrollContainerState() (+13 more)

### Community 17 - "Custom Data"
Cohesion: 0.15
Nodes (24): ExtractCustomData(), extractSliceElements(), extractStructFields(), getXMLAttributeName(), TestExtractCustomData_AETag(), TestExtractCustomData_AETagIgnore(), TestExtractCustomData_Bool(), TestExtractCustomData_EmptySlice() (+16 more)

### Community 18 - "Highlight"
Cohesion: 0.25
Nodes (13): highlightManager, AddHighlight(), ClearHighlights(), drawHighlightsCallback(), newHighlightManager(), SetHighlightDuration(), TestAddHighlight(), TestClearHighlights() (+5 more)

### Community 19 - "RPC Messages"
Cohesion: 0.12
Nodes (16): CustomParams, CustomResult, GetMousePositionResult, GetWheelPositionResult, InputParams, InputResult, MouseParams, MouseResult (+8 more)

### Community 20 - "Launch Socket Tests"
Cohesion: 0.21
Nodes (14): launchSocketPath(), TestExecutableNotFound(), TestLaunchExitsAfterCLIQuery(), TestLaunchSocketExistsBeforeGameStart(), TestMultipleCLIQueriesAfterCrash(), TestPreRPCCrashDiagnostics(), launchSockPath(), TestFindRunningGames_DeduplicationPrefersLaunch() (+6 more)

### Community 21 - "AfterDraw"
Cohesion: 0.21
Nodes (12): AfterDraw, Callback Registry Pattern, DrawHighlights, Import Cycle Solution, AfterDraw(), RegisterDrawHighlights(), TestAfterDrawNoCallback(), TestAfterDrawWithCallback() (+4 more)

### Community 22 - "Doc Rewrite"
Cohesion: 0.2
Nodes (12): Alphabetical Ordering, API Reference, Text-based Decision Trees, Documentation Rewrite, E2E Testing, Explicit Notes, Flat Structure, LLM Optimization (+4 more)

### Community 23 - "Script Parser"
Cohesion: 0.31
Nodes (9): Parse(), ParseBytes(), ParseString(), stripComments(), TestParseComments(), TestParseCustom(), TestParseInvalid(), TestParseRepeat() (+1 more)

### Community 24 - "Output/CR"
Cohesion: 0.43
Nodes (7): Carriage Return Interpretation, CarriageReturnWriter, Carriage Return Interpretation, Mutex-Protected Output, OutputManager, Rationale: Carriage Return Interpretation, Rationale: Mutex Protection

### Community 25 - "Config Tests"
Cohesion: 0.67
Nodes (0): 

### Community 26 - "Widget State Extract"
Cohesion: 0.67
Nodes (3): ExtractWidgetState, Rationale: Widget State Extraction, Widget State Extraction

### Community 27 - "Messages Tests"
Cohesion: 1.0
Nodes (0): 

### Community 28 - "Custom"
Cohesion: 1.0
Nodes (0): 

### Community 29 - "Doc"
Cohesion: 1.0
Nodes (0): 

### Community 30 - "Errors"
Cohesion: 1.0
Nodes (0): 

### Community 31 - "Caller Export"
Cohesion: 1.0
Nodes (0): 

### Community 32 - "autoui Doc"
Cohesion: 1.0
Nodes (0): 

### Community 33 - "testkit Doc"
Cohesion: 1.0
Nodes (0): 

### Community 34 - "Doc Verify"
Cohesion: 1.0
Nodes (1): docverify

### Community 35 - "Socket Env"
Cohesion: 1.0
Nodes (1): AUTOEBITEN_SOCKET

### Community 36 - "Go Code Extract"
Cohesion: 1.0
Nodes (1): Go Code Extraction

### Community 37 - "Troubleshooting"
Cohesion: 1.0
Nodes (1): Troubleshooting

## Knowledge Gaps
- **57 isolated node(s):** `testPlayer`, `testInventoryItem`, `testGameState`, `taggedPlayer`, `gameWithInterface` (+52 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Messages Tests`** (2 nodes): `messages_test.go`, `TestErrorCodes()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Custom`** (1 nodes): `custom.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Doc`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Errors`** (1 nodes): `errors.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Caller Export`** (1 nodes): `caller_export_test.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `autoui Doc`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `testkit Doc`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Doc Verify`** (1 nodes): `docverify`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Socket Env`** (1 nodes): `AUTOEBITEN_SOCKET`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Go Code Extract`** (1 nodes): `Go Code Extraction`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Troubleshooting`** (1 nodes): `Troubleshooting`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `autoui Core` to `Input Capture`, `Custom Commands`, `Command Execution`, `autoui Features`, `autoui Finder`, `Registry/Caller`, `Script Executor`, `State Export Tests`, `Script Parser`?**
  _High betweenness centrality (0.101) - this node is a cross-community bridge._
- **Why does `contains()` connect `Custom Commands` to `Command Execution`, `autoui Finder`, `Launch/Mode`, `autoui Core`, `Registry/Caller`, `Custom Handlers`, `Testkit Config`, `RadioGroup`, `Reflection/TabBook`, `Launch Socket Tests`?**
  _High betweenness centrality (0.090) - this node is a cross-community bridge._
- **Why does `NewContext()` connect `Custom Commands` to `Testkit Config`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **Are the 77 inferred relationships involving `contains()` (e.g. with `main()` and `TestOutputManagerDiffAndUpdateSnapshot()`) actually correct?**
  _`contains()` has 77 INFERRED edges - model-reasoned connections that need verification._
- **Are the 45 inferred relationships involving `WalkTree()` (e.g. with `TestHandleExistsCommand_Found()` and `TestHandleExistsCommand_NotFound()`) actually correct?**
  _`WalkTree()` has 45 INFERRED edges - model-reasoned connections that need verification._
- **Are the 38 inferred relationships involving `createTestNineSlice()` (e.g. with `TestHandleExistsCommand_Found()` and `TestHandleExistsCommand_NotFound()`) actually correct?**
  _`createTestNineSlice()` has 38 INFERRED edges - model-reasoned connections that need verification._
- **Are the 39 inferred relationships involving `ExtractWidgetInfo()` (e.g. with `ExtractWidgetState()` and `ExtractCustomData()`) actually correct?**
  _`ExtractWidgetInfo()` has 39 INFERRED edges - model-reasoned connections that need verification._
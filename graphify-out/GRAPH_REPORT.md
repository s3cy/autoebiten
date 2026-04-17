# Graph Report - /Users/s3cy/Desktop/go/autoebiten  (2026-04-17)

## Corpus Check
- 172 files · ~169,720 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1469 nodes · 3555 edges · 59 communities detected
- Extraction: 48% EXTRACTED · 48% INFERRED · 0% AMBIGUOUS · INFERRED: 1700 edges (avg confidence: 0.8)
- Token cost: 98,900 input · 29,500 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Core Input & Server|Core Input & Server]]
- [[_COMMUNITY_CLI Commands & Recording|CLI Commands & Recording]]
- [[_COMMUNITY_autoui Widget Search|autoui Widget Search]]
- [[_COMMUNITY_autoui Command Handlers|autoui Command Handlers]]
- [[_COMMUNITY_Launch & Game Management|Launch & Game Management]]
- [[_COMMUNITY_Black-box Testing|Black-box Testing]]
- [[_COMMUNITY_autoui Demo Game|autoui Demo Game]]
- [[_COMMUNITY_autoui Registry & Proxy|autoui Registry & Proxy]]
- [[_COMMUNITY_E2E Test Infrastructure|E2E Test Infrastructure]]
- [[_COMMUNITY_Script Execution Engine|Script Execution Engine]]
- [[_COMMUNITY_Crash Diagnostics|Crash Diagnostics]]
- [[_COMMUNITY_Documentation & Skills|Documentation & Skills]]
- [[_COMMUNITY_autoui User Commands|autoui User Commands]]
- [[_COMMUNITY_RadioGroup Proxy|RadioGroup Proxy]]
- [[_COMMUNITY_Proxy & Reflection|Proxy & Reflection]]
- [[_COMMUNITY_Doc Generation Templates|Doc Generation Templates]]
- [[_COMMUNITY_State Exporter|State Exporter]]
- [[_COMMUNITY_Widget State Extraction|Widget State Extraction]]
- [[_COMMUNITY_Custom Data Extraction|Custom Data Extraction]]
- [[_COMMUNITY_autoui Highlight System|autoui Highlight System]]
- [[_COMMUNITY_Docgen Context|Docgen Context]]
- [[_COMMUNITY_Custom Commands Demo|Custom Commands Demo]]
- [[_COMMUNITY_Script Parser|Script Parser]]
- [[_COMMUNITY_Demo Commands|Demo Commands]]
- [[_COMMUNITY_Build Tags|Build Tags]]
- [[_COMMUNITY_Docgen Game Session|Docgen Game Session]]
- [[_COMMUNITY_autoui Design Specs|autoui Design Specs]]
- [[_COMMUNITY_Custom Command Types|Custom Command Types]]
- [[_COMMUNITY_Test Games|Test Games]]
- [[_COMMUNITY_Documentation Plans|Documentation Plans]]
- [[_COMMUNITY_Config Testing|Config Testing]]
- [[_COMMUNITY_Carriage Return & Diagnostics|Carriage Return & Diagnostics]]
- [[_COMMUNITY_RPC Messages Test|RPC Messages Test]]
- [[_COMMUNITY_Custom Server Handler|Custom Server Handler]]
- [[_COMMUNITY_testkit Documentation|testkit Documentation]]
- [[_COMMUNITY_testkit Errors|testkit Errors]]
- [[_COMMUNITY_autoui Caller Export Test|autoui Caller Export Test]]
- [[_COMMUNITY_autoui Documentation|autoui Documentation]]
- [[_COMMUNITY_autoui Internal Doc|autoui Internal Doc]]
- [[_COMMUNITY_Main CLI Entry|Main CLI Entry]]
- [[_COMMUNITY_Docgen CLI Entry|Docgen CLI Entry]]
- [[_COMMUNITY_Output Manager Test|Output Manager Test]]
- [[_COMMUNITY_Output Testing|Output Testing]]
- [[_COMMUNITY_Test Error Types|Test Error Types]]
- [[_COMMUNITY_RPC Socket Utilities|RPC Socket Utilities]]
- [[_COMMUNITY_Custom Commands Example|Custom Commands Example]]
- [[_COMMUNITY_Crash Diagnostic Example|Crash Diagnostic Example]]
- [[_COMMUNITY_Simple Game Example|Simple Game Example]]
- [[_COMMUNITY_testkit Error States|testkit Error States]]
- [[_COMMUNITY_autoui Command Docs|autoui Command Docs]]
- [[_COMMUNITY_TabInfo Type|TabInfo Type]]
- [[_COMMUNITY_autoui Internal Docs|autoui Internal Docs]]
- [[_COMMUNITY_CLAUDE.md Project Doc|CLAUDE.md Project Doc]]
- [[_COMMUNITY_autoui At Command|autoui At Command]]
- [[_COMMUNITY_Input Command Type|Input Command Type]]
- [[_COMMUNITY_Mouse Command Type|Mouse Command Type]]
- [[_COMMUNITY_Wheel Command Type|Wheel Command Type]]
- [[_COMMUNITY_XPath 1.0 Support|XPath 1.0 Support]]
- [[_COMMUNITY_Tick System|Tick System]]

## God Nodes (most connected - your core abstractions)
1. `contains()` - 82 edges
2. `WalkTree()` - 49 edges
3. `createTestNineSlice()` - 45 edges
4. `ExtractWidgetInfo()` - 43 edges
5. `NewContext()` - 35 edges
6. `InvokeMethod()` - 35 edges
7. `NewCommandExecutor()` - 31 edges
8. `Game` - 28 edges
9. `CommandExecutor` - 26 edges
10. `Handler` - 25 edges

## Surprising Connections (you probably didn't know these)
- `readme_doc` ----> `custom_command.go - Custom command registration wrapper`  [EXTRACTED]
  README.md → custom_command.go
- `TestParsePath()` --calls--> `parsePath()`  [INFERRED]
  state_exporter_test.go → state_exporter.go
- `Capture()` --calls--> `Capture()`  [INFERRED]
  autoebiten_default.go → integrate/integrate.go
- `TestStateExporterJSONTags()` --calls--> `navigatePath()`  [INFERRED]
  state_exporter_edge_test.go → state_exporter.go
- `TestStateExporterInterfaceField()` --calls--> `navigatePath()`  [INFERRED]
  state_exporter_edge_test.go → state_exporter.go

## Hyperedges (group relationships)
- **cli_architecture** — cli_commandexecutor, cli_writer, cli_launchcommand, cli_launchoptions, cli_condition, cli_waitlogger, cli_schemavalidator [0.8]
- **script_architecture** — script_script, script_executor, script_inputcmd, script_mousecmd, script_wheelcmd, script_screenshotcmd, script_delaycmd, script_customcmd, script_statecmd, script_waitcmd, script_repeatcmd, script_commandwrapper, script_commandschema, script_scriptschema, script_internalwrapper, script_repeatschema [0.8]
- **server_architecture** — server_serverhandler, rpc_handler [0.8]
- **output_architecture** — output_outputmanager, output_carriagereturnwriter, output_filepath [0.8]
- **test_architecture** — cli_schemavalidator, test_testerror [0.7]
- **cli_command_pattern** — cli_commandexecutor, cli_writer, cli_launchcommand [0.8]
- **script_visitor_pattern** — script_executor, script_commandwrapper, script_inputcmd, script_mousecmd, script_wheelcmd, script_screenshotcmd, script_delaycmd, script_customcmd, script_statecmd, script_waitcmd, script_repeatcmd [0.7]
- **script_facade_pattern** — script_script, script_executor, script_parse, script_parsebytes [0.6]
- **output_manager_pattern** — output_outputmanager, output_carriagereturnwriter [0.7]
- **server_handler_pattern** — server_serverhandler, rpc_handler [0.8]
- **cli_factory_pattern** — cli_newcommandexecutor, cli_newwriter, cli_newlaunchcommand, script_newexecutor [0.7]
- **script_strategy_pattern** — script_executor, script_executor_setinputfunc, script_executor_setmousefunc, script_executor_setwheelfunc, script_executor_setscreenshotfunc, script_executor_setcustomfunc, script_executor_setstatefunc, script_executor_setwaitfunc [0.7]
- **cli_proxy_layer** — cli_detectproxyordirect, cli_sendrequestwithproxy, cli_handlerequestwithproxy [0.8]
- **cli_launch_workflow** — cli_launchcommand, cli_createlaunchsocket, cli_creategamecommand, cli_teeoutput, cli_acceptorloop, cli_waitforgamerpc, cli_setupsignalhandling, cli_waitforexit, cli_terminategame, cli_cleanup [0.8]
- **script_execution_workflow** — script_parse, script_parsebytes, script_newexecutor, script_executor_execute, script_executor_executescommands, script_executor_executecommand [0.8]
- **server_request_workflow** — server_update, server_startsocketserver, server_processrequest, server_handler_handleinput, server_handler_handlemouse, server_handler_handlewheel, server_handler_handlescreenshot, server_handler_handlecustom [0.8]
- **output_capture_workflow** — output_derivepaths, output_createlogfile, output_newoutputmanager, output_newcarriagereturnwriter, output_write, output_diffandupdatesnapshot [0.7]
- **cli_wait_condition_workflow** — cli_runwaitforcommand, cli_parsecondition, cli_pollcondition, cli_checkcondition, cli_extractresponsepath, cli_timeouterror [0.8]
- **test_validation_suite** — cli_testlaunch, cli_testschema, cli_testresponse, cli_teststate, cli_testwait, script_testparse, script_testexecutor [0.7]
- **Template Function Map Pattern** — docgen_template_funcmap, docgen_config, docgen_context, docgen_launchgame, docgen_endgame, docgen_executecommand, docgen_delay, docgen_normalize, docgen_verifyoutputs, docgen_extractgocode [0.9]
- **Docgen Game Lifecycle** — docgen_launchgame, docgen_endgame, docgen_executecommand, docgen_gamesession, docgen_context [0.85]
- **JSON-RPC Protocol Types** — rpc_rpcrequest, rpc_rpcresponse, rpc_rpcerror, rpc_request [0.95]
- **Command Handler Interface Pattern** — rpc_handler, rpc_inputparams, rpc_mouseparams, rpc_wheelparams, rpc_screenshotparams, rpc_processrequest [0.9]
- **Socket Discovery and Connection** — rpc_socketpath, rpc_findrunninggames, rpc_autoselectgame, rpc_gameinfo, rpc_serve, rpc_client [0.85]
- **Functional Options Pattern** — testkit_option, testkit_config, testkit_withtimeout, testkit_withargs, testkit_withenv [0.9]
- **Command Registry Pattern** — custom_register, custom_unregister, custom_get, custom_list, custom_context [0.9]
- **testing_architecture** — testkit_game, testkit_mock, whitebox_test, blackbox_test [1.0]
- **autoui_integration_pattern** — autoui_register, example_autoui, ebitenui_ui, ebitenui_widget [1.0]
- **testgame_pattern** — testgame_stateful, testgame_simple, testgame_custom, integration_test [1.0]
- **example_games** — example_simple, example_custom_cmd, example_crash_diag, example_state_exporter, example_autoui [1.0]
- **state_export_pattern** — example_state_exporter, gamestate, player, enemy, autoebiten_register_state_exporter [1.0]
- **core_integration_pattern** — autoebiten_update, autoebiten_capture, autoebiten_register [1.0]
- **functional_options_pattern** — testkit_option, testkit_config, testkit_withtimeout, testkit_withargs, testkit_withenv, testkit_defaultconfig [1.0]
- **Widget Tree Traversal Flow** — autoui_tree, autoui_xml, autoui_xpath, autoui_finder, widgetinfo, widgetnode [1.0]
- **Method Invocation Flow** — autoui_handlers, autoui_caller, autoui_proxy, autoui_proxy_radiogroup [1.0]
- **Command Registration Pattern** — autoui_handlers, autoui_tree, autoui_finder, autoui_xpath, autoui_caller, autoui_highlight, autoui_registry [1.0]
- **RadioGroup Management** — autoui_registry, autoui_proxy_radiogroup, autoui_tree [1.0]
- **Visual Highlight System** — autoui_highlight, highlightmanager, highlight [1.0]
- **XML Marshaling Pipeline** — widgetinfo, widgetnode, autoui_xml [1.0]
- **Widget Finding System** — autoui_finder, widgetinfo [1.0]
- **Reflection-based Method Caller** — autoui_caller, proxyhandler, radiogrouphandler [1.0]
- **autoui_internal_package** — autoui_internal_reflection, autoui_internal_widgetstate, autoui_internal_customdata, autoui_internal_doc, reflection_test, widgetstate_test, customdata_test [INFERRED 1.00]
- **e2e_test_suite** — e2e_tests, e2e_launch_tests, e2e_crash_tests, e2e_output_tests [INFERRED 1.00]
- **documentation_set** — readme_doc, integration_doc, commands_doc, claude_doc [INFERRED 1.00]
- **integration_methods** — patch_method, library_method, build_tags [INFERRED 1.00]
- **automation_features** — input_injection, screenshot_capture, script_executor, custom_command, state_exporter, wait_for_command [INFERRED 0.90]
- **crash_diagnostic_system** — crash_diagnostics, launch_socket, proxy_server, output_manager, carriage_return_writer [INFERRED 0.80]
- **Documentation File Structure** — docs/tutorial.md, docs/autoui.md, docs/SPEC.md, docs/testkit.md [INFERRED 0.80]
- **Implementation Plans** — docs/superpowers/plans/2026-04-09-widget-state-extraction.md, docs/superpowers/plans/2026-04-08-launch-socket-crash-diagnostics.md, docs/superpowers/plans/2026-04-07-carriage-return-interpretation-plan.md, docs/superpowers/plans/2026-04-06-documentation-rewrite.md, docs/superpowers/plans/2026-04-08-autoui.md, docs/superpowers/plans/2026-04-10-doc-template-system.md, docs/superpowers/plans/2026-04-09-doc-example-automation.md, docs/superpowers/plans/2026-04-11-autoui-call-type-support.md, docs/superpowers/plans/2026-04-09-integrate-afterdraw.md, docs/superpowers/plans/2026-04-07-game-output-capture.md, docs/superpowers/plans/2026-04-09-autoui-bugfixes.md [INFERRED 0.90]
- **April 2026 autoui Feature Set** —  [INFERRED 0.80]
- **Documentation Generation System** —  [INFERRED 0.90]
- **Launch Infrastructure Components** —  [INFERRED 0.90]
- **autoui_architecture** — autoui, autoebiten, tree_walker, finder, caller, ebitenui [INFERRED 1.00]
- **integration_methods** — patch_method, library_method, ebiten, autoebiten [INFERRED 1.00]
- **testing_modes** — black_box_testing, white_box_testing, testkit [INFERRED 1.00]
- **autoui_commands** — autoui_tree, autoui_at, autoui_find, autoui_xpath, autoui_exists, autoui_call, autoui_highlight [INFERRED 1.00]
- **crash_diagnostics_flow** — launch_command, crash_diagnostics, log_diff, proxy_error [INFERRED 1.00]
- **e2e_test_pattern** — autoui_tree, autoui_find, autoui_call, screenshot_command, testkit [INFERRED 0.90]
- **custom_data_extraction** — custom_data, ae_tag, widgetnode [INFERRED 1.00]
- **widget_automation_workflow** — autoui_tree, autoui_find, autoui_call [INFERRED 0.90]
- **heal_flow** — autoebiten_command_heal, autoebiten_command_system, autoebiten_player_stats [INFERRED 1.00]
- **damage_flow** — autoebiten_command_damage, autoebiten_command_system, autoebiten_player_stats [INFERRED 1.00]
- **info_flow** — autoebiten_command_getplayerinfo, autoebiten_command_system, autoebiten_player_stats [INFERRED 1.00]

## Communities

### Community 0 - "Core Input & Server"
Cohesion: 0.03
Nodes (87): Capture(), CursorPosition(), IsKeyJustPressed(), IsKeyJustReleased(), IsKeyPressed(), IsMouseButtonJustPressed(), IsMouseButtonJustReleased(), IsMouseButtonPressed() (+79 more)

### Community 1 - "CLI Commands & Recording"
Cohesion: 0.05
Nodes (74): CommandExecutor, Condition, waitLogger, Writer, detectProxyOrDirect(), NewCommandExecutor(), captureOutput(), TestHandleResponseBothDiffAndProxyError() (+66 more)

### Community 2 - "autoui Widget Search"
Cohesion: 0.04
Nodes (97): autoui, stringKeyValue, WidgetNode, xpath.go, Caller Whitelist, ExtractCustomData, filterWidgets(), FindAt() (+89 more)

### Community 3 - "autoui Command Handlers"
Cohesion: 0.05
Nodes (85): CallRequest, CallResponse, CoordinateRequest, ExistsResponse, HighlightRequest, mockCommandContext, RegisterOptions, GetCustomCommand() (+77 more)

### Community 4 - "Launch & Game Management"
Cohesion: 0.04
Nodes (66): Mode, LaunchCommand, LaunchOptions, handleTestConnection(), mustMarshal(), newTestClient(), newTestHandler(), sendRequest() (+58 more)

### Community 5 - "Black-box Testing"
Cohesion: 0.05
Nodes (60): Config, getStateExporterGameBinary(), TestEnemyStateQuery(), TestHealthModification(), TestPlayerMovement(), TestScreenshotCapture(), blackbox_test, DefaultConfig() (+52 more)

### Community 6 - "autoui Demo Game"
Cohesion: 0.04
Nodes (36): autoebiten, Game, PlayerCard, Game, Game, CustomGame, Documentation Template System, docgen (+28 more)

### Community 7 - "autoui Registry & Proxy"
Cohesion: 0.08
Nodes (64): WidgetInfo, convertArg(), convertArgs(), InvokeMethod(), isWhitelistedSignature(), TestConvertArg_AnyParam(), TestConvertArg_IntToEnum(), TestConvertArg_NonEmptyInterfaceNotImplemented() (+56 more)

### Community 8 - "E2E Test Infrastructure"
Cohesion: 0.07
Nodes (37): testHandler, decodeParams(), ErrorResponse(), marshalResult(), ProcessRequest(), delayedMockClient, GameClient, Handler (+29 more)

### Community 9 - "Script Execution Engine"
Cohesion: 0.06
Nodes (41): unmarshalCommand(), unmarshalRepeat(), TestScriptExecutorEndToEnd(), formatCustomCmd(), formatInputCmd(), formatMouseCmd(), formatScreenshotCmd(), formatWheelCmd() (+33 more)

### Community 10 - "Crash Diagnostics"
Cohesion: 0.06
Nodes (45): EnsureTargetPID(), launchSocketPath(), TestExecutableNotFound(), TestLaunchExitsAfterCLIQuery(), TestLaunchSocketExistsBeforeGameStart(), TestMultipleCLIQueriesAfterCrash(), TestPreRPCCrashDiagnostics(), persistentPreRunRootCommand() (+37 more)

### Community 11 - "Documentation & Skills"
Cohesion: 0.05
Nodes (48): autoui_exists, Black-Box Testing, build_tags, carriage_return_writer, commands_doc, crash_diagnostics, custom_command.go - Custom command registration wrapper, e2e_crash_tests (+40 more)

### Community 12 - "autoui User Commands"
Cohesion: 0.05
Nodes (49): ae_tag, autoui_call, caller.go, autoui_find, finder.go, handlers.go, Highlight, autoui_internal_customdata (+41 more)

### Community 13 - "RadioGroup Proxy"
Cohesion: 0.08
Nodes (42): RadioGroupElementInfo, RadioGroupHandler, schemaValidator, runSchemaCommand(), handleRadioGroupActiveIndex(), handleRadioGroupActiveLabel(), handleRadioGroupElements(), handleRadioGroupSetActiveByIndex() (+34 more)

### Community 14 - "Proxy & Reflection"
Cohesion: 0.08
Nodes (38): ProxyHandler, TabInfo, GetProxyHandler(), handleSetTabByIndex(), handleSetTabByLabel(), handleTabIndex(), handleTabLabel(), handleTabs() (+30 more)

### Community 15 - "Doc Generation Templates"
Cohesion: 0.09
Nodes (42): commandFunc(), configFunc(), delayFunc(), dictFunc(), endGameFunc(), FuncMap(), gocodeFunc(), launchGameFunc() (+34 more)

### Community 16 - "State Exporter"
Cohesion: 0.12
Nodes (25): gameWithInterface, taggedGameState, taggedPlayer, TestStateExporterInterfaceField(), TestStateExporterInterfaceWithPointer(), TestStateExporterJSONTags(), TestStateExporterNilInterface(), getFieldByName() (+17 more)

### Community 17 - "Widget State Extraction"
Cohesion: 0.14
Nodes (21): extractButtonState(), extractCheckboxState(), extractComboButtonState(), extractLabelState(), extractListComboButtonState(), extractListState(), extractProgressBarState(), extractScrollContainerState() (+13 more)

### Community 18 - "Custom Data Extraction"
Cohesion: 0.15
Nodes (24): ExtractCustomData(), extractSliceElements(), extractStructFields(), getXMLAttributeName(), TestExtractCustomData_AETag(), TestExtractCustomData_AETagIgnore(), TestExtractCustomData_Bool(), TestExtractCustomData_EmptySlice() (+16 more)

### Community 19 - "autoui Highlight System"
Cohesion: 0.22
Nodes (14): highlightManager, AddHighlight(), ClearHighlights(), DrawHighlights(), drawHighlightsCallback(), newHighlightManager(), SetHighlightDuration(), TestAddHighlight() (+6 more)

### Community 20 - "Docgen Context"
Cohesion: 0.18
Nodes (8): TestContextAddOutput(), TestContextGetOutputsEmpty(), TestContextGetOutputsReturnsCopy(), TestContextSetConfig(), TestNewContext(), Config, Context, NormalizeRule

### Community 21 - "Custom Commands Demo"
Cohesion: 0.13
Nodes (15): command_console, command_list, custom_commands_demo, damage_cmd, dark_blue_theme, deferred_cmd, ebitengine_game, echo_cmd (+7 more)

### Community 22 - "Script Parser"
Cohesion: 0.31
Nodes (9): Parse(), ParseBytes(), ParseString(), stripComments(), TestParseComments(), TestParseCustom(), TestParseInvalid(), TestParseRepeat() (+1 more)

### Community 23 - "Demo Commands"
Cohesion: 0.2
Nodes (10): damage Command, deferred Command, echo Command, getPlayerInfo Command, heal Command, Custom Commands Demo, Command List, Health Display (+2 more)

### Community 24 - "Build Tags"
Cohesion: 0.33
Nodes (9): autoebiten_default.go - non-release build with integration, autoebiten.go - Mode constants and management, autoebiten_release.go - release build stubs, internal/input/input.go - VirtualInput state management, internal/input/input_time.go - InputTime tick/subtick, internal/input/keys.go - Key constants and lookup, internal/input/mouse_buttons.go - Mouse button constants, integrate/integrate.go - Ebiten integration layer (+1 more)

### Community 25 - "Docgen Game Session"
Cohesion: 0.25
Nodes (8): Delay, EndGame, ExecuteCommand, ExtractGoCode, LaunchGame, Normalize, FuncMap, VerifyOutputs

### Community 26 - "autoui Design Specs"
Cohesion: 0.25
Nodes (8): autoui Bug Fixes Design, autoui.call Extended Type Support Design, autoui.exists Command Design, autoui.exists Command Implementation Plan, autoui RadioGroup & TabBook Support Design, autoui RadioGroup & TabBook Support Implementation Plan, integrate.AfterDraw() Design, Widget State Extraction Expansion Design

### Community 27 - "Custom Command Types"
Cohesion: 0.43
Nodes (8): Damage Command, Deferred Command, Echo Command, GetPlayerInfo Command, Heal Command, Command Processing System, Custom Commands Demo, Player Stats Display

### Community 28 - "Test Games"
Cohesion: 0.33
Nodes (6): integration_test, inventory_item, player, testgame_custom, testgame_simple, testgame_stateful

### Community 29 - "Documentation Plans"
Cohesion: 0.5
Nodes (4): Doc Example Automation Design, Document Template Rewrite Workflow Design, Documentation Template System Rewrite Design, Documentation Rewrite Design

### Community 30 - "Config Testing"
Cohesion: 0.67
Nodes (0): 

### Community 31 - "Carriage Return & Diagnostics"
Cohesion: 0.67
Nodes (3): Carriage Return Interpretation + Mutex-Protected Output Design, Game Output Capture Design, Launch Socket for Crash Diagnostics Design

### Community 32 - "RPC Messages Test"
Cohesion: 1.0
Nodes (0): 

### Community 33 - "Custom Server Handler"
Cohesion: 1.0
Nodes (0): 

### Community 34 - "testkit Documentation"
Cohesion: 1.0
Nodes (0): 

### Community 35 - "testkit Errors"
Cohesion: 1.0
Nodes (0): 

### Community 36 - "autoui Caller Export Test"
Cohesion: 1.0
Nodes (0): 

### Community 37 - "autoui Documentation"
Cohesion: 1.0
Nodes (0): 

### Community 38 - "autoui Internal Doc"
Cohesion: 1.0
Nodes (0): 

### Community 39 - "Main CLI Entry"
Cohesion: 1.0
Nodes (1): cmd/autoebiten/main.go - CLI entry point with cobra

### Community 40 - "Docgen CLI Entry"
Cohesion: 1.0
Nodes (1): cmd/docgen/main.go - Documentation generator

### Community 41 - "Output Manager Test"
Cohesion: 1.0
Nodes (1): internal/output/manager_test.go - Output manager tests

### Community 42 - "Output Testing"
Cohesion: 1.0
Nodes (1): internal/output/output_test.go - Output path tests

### Community 43 - "Test Error Types"
Cohesion: 1.0
Nodes (1): test_testerror

### Community 44 - "RPC Socket Utilities"
Cohesion: 1.0
Nodes (1): SocketPath

### Community 45 - "Custom Commands Example"
Cohesion: 1.0
Nodes (1): example_custom_cmd

### Community 46 - "Crash Diagnostic Example"
Cohesion: 1.0
Nodes (1): example_crash_diag

### Community 47 - "Simple Game Example"
Cohesion: 1.0
Nodes (1): example_simple

### Community 48 - "testkit Error States"
Cohesion: 1.0
Nodes (1): testkit_errinvalidstate

### Community 49 - "autoui Command Docs"
Cohesion: 1.0
Nodes (1): doc.go

### Community 50 - "TabInfo Type"
Cohesion: 1.0
Nodes (1): TabInfo

### Community 51 - "autoui Internal Docs"
Cohesion: 1.0
Nodes (1): autoui_internal_doc

### Community 52 - "CLAUDE.md Project Doc"
Cohesion: 1.0
Nodes (1): claude_doc

### Community 53 - "autoui At Command"
Cohesion: 1.0
Nodes (1): autoui_at

### Community 54 - "Input Command Type"
Cohesion: 1.0
Nodes (1): input_command

### Community 55 - "Mouse Command Type"
Cohesion: 1.0
Nodes (1): mouse_command

### Community 56 - "Wheel Command Type"
Cohesion: 1.0
Nodes (1): wheel_command

### Community 57 - "XPath 1.0 Support"
Cohesion: 1.0
Nodes (1): xpath_1_0

### Community 58 - "Tick System"
Cohesion: 1.0
Nodes (1): ticks

## Knowledge Gaps
- **67 isolated node(s):** `testPlayer`, `testInventoryItem`, `testGameState`, `taggedPlayer`, `gameWithInterface` (+62 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `RPC Messages Test`** (2 nodes): `messages_test.go`, `TestErrorCodes()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Custom Server Handler`** (1 nodes): `custom.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `testkit Documentation`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `testkit Errors`** (1 nodes): `errors.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `autoui Caller Export Test`** (1 nodes): `caller_export_test.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `autoui Documentation`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `autoui Internal Doc`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Main CLI Entry`** (1 nodes): `cmd/autoebiten/main.go - CLI entry point with cobra`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Docgen CLI Entry`** (1 nodes): `cmd/docgen/main.go - Documentation generator`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Output Manager Test`** (1 nodes): `internal/output/manager_test.go - Output manager tests`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Output Testing`** (1 nodes): `internal/output/output_test.go - Output path tests`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Test Error Types`** (1 nodes): `test_testerror`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `RPC Socket Utilities`** (1 nodes): `SocketPath`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Custom Commands Example`** (1 nodes): `example_custom_cmd`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Crash Diagnostic Example`** (1 nodes): `example_crash_diag`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Simple Game Example`** (1 nodes): `example_simple`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `testkit Error States`** (1 nodes): `testkit_errinvalidstate`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `autoui Command Docs`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `TabInfo Type`** (1 nodes): `TabInfo`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `autoui Internal Docs`** (1 nodes): `autoui_internal_doc`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `CLAUDE.md Project Doc`** (1 nodes): `claude_doc`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `autoui At Command`** (1 nodes): `autoui_at`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Input Command Type`** (1 nodes): `input_command`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Mouse Command Type`** (1 nodes): `mouse_command`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Wheel Command Type`** (1 nodes): `wheel_command`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `XPath 1.0 Support`** (1 nodes): `xpath_1_0`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Tick System`** (1 nodes): `ticks`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `contains()` connect `autoui Command Handlers` to `CLI Commands & Recording`, `autoui Widget Search`, `Launch & Game Management`, `Black-box Testing`, `autoui Demo Game`, `autoui Registry & Proxy`, `E2E Test Infrastructure`, `Crash Diagnostics`, `RadioGroup Proxy`, `Proxy & Reflection`, `Doc Generation Templates`?**
  _High betweenness centrality (0.055) - this node is a cross-community bridge._
- **Why does `Game` connect `Black-box Testing` to `CLI Commands & Recording`, `Crash Diagnostics`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Why does `main()` connect `autoui Demo Game` to `Core Input & Server`, `CLI Commands & Recording`, `autoui Command Handlers`, `Script Execution Engine`, `Crash Diagnostics`, `Documentation & Skills`, `State Exporter`, `Script Parser`?**
  _High betweenness centrality (0.047) - this node is a cross-community bridge._
- **Are the 77 inferred relationships involving `contains()` (e.g. with `main()` and `TestOutputManagerDiffAndUpdateSnapshot()`) actually correct?**
  _`contains()` has 77 INFERRED edges - model-reasoned connections that need verification._
- **Are the 45 inferred relationships involving `WalkTree()` (e.g. with `TestHandleExistsCommand_Found()` and `TestHandleExistsCommand_NotFound()`) actually correct?**
  _`WalkTree()` has 45 INFERRED edges - model-reasoned connections that need verification._
- **Are the 38 inferred relationships involving `createTestNineSlice()` (e.g. with `TestHandleExistsCommand_Found()` and `TestHandleExistsCommand_NotFound()`) actually correct?**
  _`createTestNineSlice()` has 38 INFERRED edges - model-reasoned connections that need verification._
- **Are the 39 inferred relationships involving `ExtractWidgetInfo()` (e.g. with `ExtractWidgetState()` and `ExtractCustomData()`) actually correct?**
  _`ExtractWidgetInfo()` has 39 INFERRED edges - model-reasoned connections that need verification._
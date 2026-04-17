# Graph Report - .  (2026-04-17)

## Corpus Check
- 12 files · ~173,132 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 1558 nodes · 3714 edges · 65 communities detected
- Extraction: 48% EXTRACTED · 48% INFERRED · 0% AMBIGUOUS · INFERRED: 1769 edges (avg confidence: 0.8)
- Token cost: 12,500 input · 8,500 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]
- [[_COMMUNITY_Community 23|Community 23]]
- [[_COMMUNITY_Community 24|Community 24]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 34|Community 34]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 38|Community 38]]
- [[_COMMUNITY_Community 39|Community 39]]
- [[_COMMUNITY_Community 40|Community 40]]
- [[_COMMUNITY_Community 41|Community 41]]
- [[_COMMUNITY_Community 42|Community 42]]
- [[_COMMUNITY_Community 43|Community 43]]
- [[_COMMUNITY_Community 44|Community 44]]
- [[_COMMUNITY_Community 45|Community 45]]
- [[_COMMUNITY_Community 46|Community 46]]
- [[_COMMUNITY_Community 47|Community 47]]
- [[_COMMUNITY_Community 48|Community 48]]
- [[_COMMUNITY_Community 49|Community 49]]
- [[_COMMUNITY_Community 50|Community 50]]
- [[_COMMUNITY_Community 51|Community 51]]
- [[_COMMUNITY_Community 52|Community 52]]
- [[_COMMUNITY_Community 53|Community 53]]
- [[_COMMUNITY_Community 54|Community 54]]
- [[_COMMUNITY_Community 55|Community 55]]
- [[_COMMUNITY_Community 56|Community 56]]
- [[_COMMUNITY_Community 57|Community 57]]
- [[_COMMUNITY_Community 58|Community 58]]
- [[_COMMUNITY_Community 59|Community 59]]
- [[_COMMUNITY_Community 60|Community 60]]
- [[_COMMUNITY_Community 61|Community 61]]
- [[_COMMUNITY_Community 62|Community 62]]
- [[_COMMUNITY_Community 63|Community 63]]
- [[_COMMUNITY_Community 64|Community 64]]

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
- `Mock` --semantically_similar_to--> `testkit.Launch`  [INFERRED] [semantically similar]
  testkit/mock.go → skills/using-autoebiten/SKILL.md
- `state_exporter.go - State query via reflection` --rationale_for--> `Rationale: JSON Tags Ignored for Path Navigation`  [INFERRED]
  state_exporter.go → skills/using-autoebiten/references/integration.md
- `Mock` --describes--> `using-autoebiten Skill`  [EXTRACTED]
  testkit/mock.go → skills/using-autoebiten/SKILL.md
- `state_exporter.go - State query via reflection` --enables--> `State Queries`  [INFERRED]
  state_exporter.go → skills/using-autoebiten/SKILL.md
- `readme_doc` ----> `custom_command.go - Custom command registration wrapper`  [EXTRACTED]
  README.md → custom_command.go

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

### Community 0 - "Community 0"
Cohesion: 0.04
Nodes (96): CommandExecutor, Condition, waitLogger, Writer, detectProxyOrDirect(), NewCommandExecutor(), captureOutput(), TestHandleResponseBothDiffAndProxyError() (+88 more)

### Community 1 - "Community 1"
Cohesion: 0.04
Nodes (126): CallRequest, CallResponse, CoordinateRequest, ExistsResponse, HighlightRequest, mockCommandContext, RegisterOptions, GetCustomCommand() (+118 more)

### Community 2 - "Community 2"
Cohesion: 0.03
Nodes (87): Capture(), CursorPosition(), IsKeyJustPressed(), IsKeyJustReleased(), IsKeyPressed(), IsMouseButtonJustPressed(), IsMouseButtonJustReleased(), IsMouseButtonPressed() (+79 more)

### Community 3 - "Community 3"
Cohesion: 0.03
Nodes (103): ae Tags, autoui.at Command, autoui.call Command, autoui_exists, autoui.exists Command, autoui_find, autoui.find Command, finder.go (+95 more)

### Community 4 - "Community 4"
Cohesion: 0.05
Nodes (86): stringKeyValue, WidgetNode, filterWidgets(), FindAt(), FindByQuery(), FindByQueryJSON(), getWidgetAttributeValue(), matchesCriteria() (+78 more)

### Community 5 - "Community 5"
Cohesion: 0.05
Nodes (64): Mode, LaunchCommand, LaunchOptions, delayedMockClient, handleTestConnection(), mustMarshal(), newTestClient(), newTestHandler() (+56 more)

### Community 6 - "Community 6"
Cohesion: 0.06
Nodes (82): WidgetInfo, convertArg(), convertArgs(), InvokeMethod(), isWhitelistedSignature(), TestConvertArg_AnyParam(), TestConvertArg_IntToEnum(), TestConvertArg_NonEmptyInterfaceNotImplemented() (+74 more)

### Community 7 - "Community 7"
Cohesion: 0.05
Nodes (59): Config, getStateExporterGameBinary(), TestEnemyStateQuery(), TestHealthModification(), TestPlayerMovement(), TestScreenshotCapture(), blackbox_test, DefaultConfig() (+51 more)

### Community 8 - "Community 8"
Cohesion: 0.06
Nodes (41): testHandler, decodeParams(), ErrorResponse(), marshalResult(), ProcessRequest(), TestConcurrentProxyRequests(), TestProxyServerCleanup(), delayedMockClient (+33 more)

### Community 9 - "Community 9"
Cohesion: 0.04
Nodes (36): autoebiten, Game, PlayerCard, Game, Game, CustomGame, Documentation Template System, docgen (+28 more)

### Community 10 - "Community 10"
Cohesion: 0.04
Nodes (65): ae_tag, autoui, autoui_call, caller.go, handlers.go, Highlight, highlightManager, autoui_internal_customdata (+57 more)

### Community 11 - "Community 11"
Cohesion: 0.06
Nodes (46): EnsureTargetPID(), launchSocketPath(), TestExecutableNotFound(), TestLaunchExitsAfterCLIQuery(), TestLaunchSocketExistsBeforeGameStart(), TestMultipleCLIQueriesAfterCrash(), TestPreRPCCrashDiagnostics(), persistentPreRunRootCommand() (+38 more)

### Community 12 - "Community 12"
Cohesion: 0.07
Nodes (42): autoui_internal_reflection, ProxyHandler, TabInfo, GetProxyHandler(), handleSetTabByIndex(), handleSetTabByLabel(), handleTabIndex(), handleTabLabel() (+34 more)

### Community 13 - "Community 13"
Cohesion: 0.07
Nodes (44): RadioGroupElementInfo, RadioGroupHandler, schemaValidator, runSchemaCommand(), handleRadioGroupActiveIndex(), handleRadioGroupActiveLabel(), handleRadioGroupElements(), handleRadioGroupSetActiveByIndex() (+36 more)

### Community 14 - "Community 14"
Cohesion: 0.09
Nodes (18): unmarshalCommand(), unmarshalRepeat(), Entry, CommandSchema, CommandWrapper, CustomCmd, DelayCmd, InputCmd (+10 more)

### Community 15 - "Community 15"
Cohesion: 0.12
Nodes (25): gameWithInterface, taggedGameState, taggedPlayer, TestStateExporterInterfaceField(), TestStateExporterInterfaceWithPointer(), TestStateExporterJSONTags(), TestStateExporterNilInterface(), getFieldByName() (+17 more)

### Community 16 - "Community 16"
Cohesion: 0.15
Nodes (24): ExtractCustomData(), extractSliceElements(), extractStructFields(), getXMLAttributeName(), TestExtractCustomData_AETag(), TestExtractCustomData_AETagIgnore(), TestExtractCustomData_Bool(), TestExtractCustomData_EmptySlice() (+16 more)

### Community 17 - "Community 17"
Cohesion: 0.09
Nodes (26): autoui.at Command, autoui.call Command, autoui.exists Command, autoui.find Command, autoui.highlight Command, autoui.xpath Command, ae Tags (Custom Attributes), RadioGroup Operations (+18 more)

### Community 18 - "Community 18"
Cohesion: 0.18
Nodes (8): TestContextAddOutput(), TestContextGetOutputsEmpty(), TestContextGetOutputsReturnsCopy(), TestContextSetConfig(), TestNewContext(), Config, Context, NormalizeRule

### Community 19 - "Community 19"
Cohesion: 0.13
Nodes (15): command_console, command_list, custom_commands_demo, damage_cmd, dark_blue_theme, deferred_cmd, ebitengine_game, echo_cmd (+7 more)

### Community 20 - "Community 20"
Cohesion: 0.31
Nodes (9): Parse(), ParseBytes(), ParseString(), stripComments(), TestParseComments(), TestParseCustom(), TestParseInvalid(), TestParseRepeat() (+1 more)

### Community 21 - "Community 21"
Cohesion: 0.2
Nodes (10): damage Command, deferred Command, echo Command, getPlayerInfo Command, heal Command, Custom Commands Demo, Command List, Health Display (+2 more)

### Community 22 - "Community 22"
Cohesion: 0.33
Nodes (9): autoebiten_default.go - non-release build with integration, autoebiten.go - Mode constants and management, autoebiten_release.go - release build stubs, internal/input/input.go - VirtualInput state management, internal/input/input_time.go - InputTime tick/subtick, internal/input/keys.go - Key constants and lookup, internal/input/mouse_buttons.go - Mouse button constants, integrate/integrate.go - Ebiten integration layer (+1 more)

### Community 23 - "Community 23"
Cohesion: 0.25
Nodes (8): Delay, EndGame, ExecuteCommand, ExtractGoCode, LaunchGame, Normalize, FuncMap, VerifyOutputs

### Community 24 - "Community 24"
Cohesion: 0.25
Nodes (8): autoui Bug Fixes Design, autoui.call Extended Type Support Design, autoui.exists Command Design, autoui.exists Command Implementation Plan, autoui RadioGroup & TabBook Support Design, autoui RadioGroup & TabBook Support Implementation Plan, integrate.AfterDraw() Design, Widget State Extraction Expansion Design

### Community 25 - "Community 25"
Cohesion: 0.43
Nodes (8): Damage Command, Deferred Command, Echo Command, GetPlayerInfo Command, Heal Command, Command Processing System, Custom Commands Demo, Player Stats Display

### Community 26 - "Community 26"
Cohesion: 0.33
Nodes (6): integration_test, inventory_item, player, testgame_custom, testgame_simple, testgame_stateful

### Community 27 - "Community 27"
Cohesion: 0.5
Nodes (4): Doc Example Automation Design, Document Template Rewrite Workflow Design, Documentation Template System Rewrite Design, Documentation Rewrite Design

### Community 28 - "Community 28"
Cohesion: 0.67
Nodes (0): 

### Community 29 - "Community 29"
Cohesion: 0.67
Nodes (3): Carriage Return Interpretation + Mutex-Protected Output Design, Game Output Capture Design, Launch Socket for Crash Diagnostics Design

### Community 30 - "Community 30"
Cohesion: 0.67
Nodes (3): autoui.tree Command, XML Format, WidgetInfo

### Community 31 - "Community 31"
Cohesion: 1.0
Nodes (0): 

### Community 32 - "Community 32"
Cohesion: 1.0
Nodes (2): Writing Documentation Rules, autoebiten Overview

### Community 33 - "Community 33"
Cohesion: 1.0
Nodes (2): state Command, State Exporter

### Community 34 - "Community 34"
Cohesion: 1.0
Nodes (2): Thread Safety, EbitenUI Integration

### Community 35 - "Community 35"
Cohesion: 1.0
Nodes (2): launch Command, Crash Diagnostics

### Community 36 - "Community 36"
Cohesion: 1.0
Nodes (0): 

### Community 37 - "Community 37"
Cohesion: 1.0
Nodes (0): 

### Community 38 - "Community 38"
Cohesion: 1.0
Nodes (0): 

### Community 39 - "Community 39"
Cohesion: 1.0
Nodes (0): 

### Community 40 - "Community 40"
Cohesion: 1.0
Nodes (0): 

### Community 41 - "Community 41"
Cohesion: 1.0
Nodes (0): 

### Community 42 - "Community 42"
Cohesion: 1.0
Nodes (1): cmd/autoebiten/main.go - CLI entry point with cobra

### Community 43 - "Community 43"
Cohesion: 1.0
Nodes (1): cmd/docgen/main.go - Documentation generator

### Community 44 - "Community 44"
Cohesion: 1.0
Nodes (1): internal/output/manager_test.go - Output manager tests

### Community 45 - "Community 45"
Cohesion: 1.0
Nodes (1): internal/output/output_test.go - Output path tests

### Community 46 - "Community 46"
Cohesion: 1.0
Nodes (1): test_testerror

### Community 47 - "Community 47"
Cohesion: 1.0
Nodes (1): SocketPath

### Community 48 - "Community 48"
Cohesion: 1.0
Nodes (1): example_custom_cmd

### Community 49 - "Community 49"
Cohesion: 1.0
Nodes (1): example_crash_diag

### Community 50 - "Community 50"
Cohesion: 1.0
Nodes (1): example_simple

### Community 51 - "Community 51"
Cohesion: 1.0
Nodes (1): testkit_errinvalidstate

### Community 52 - "Community 52"
Cohesion: 1.0
Nodes (1): doc.go

### Community 53 - "Community 53"
Cohesion: 1.0
Nodes (1): TabInfo

### Community 54 - "Community 54"
Cohesion: 1.0
Nodes (1): autoui_internal_doc

### Community 55 - "Community 55"
Cohesion: 1.0
Nodes (1): claude_doc

### Community 56 - "Community 56"
Cohesion: 1.0
Nodes (1): autoui_at

### Community 57 - "Community 57"
Cohesion: 1.0
Nodes (1): input_command

### Community 58 - "Community 58"
Cohesion: 1.0
Nodes (1): mouse_command

### Community 59 - "Community 59"
Cohesion: 1.0
Nodes (1): wheel_command

### Community 60 - "Community 60"
Cohesion: 1.0
Nodes (1): xpath_1_0

### Community 61 - "Community 61"
Cohesion: 1.0
Nodes (1): ticks

### Community 62 - "Community 62"
Cohesion: 1.0
Nodes (1): Installation Instructions

### Community 63 - "Community 63"
Cohesion: 1.0
Nodes (1): Graphify Knowledge Graph

### Community 64 - "Community 64"
Cohesion: 1.0
Nodes (1): ping Command

## Knowledge Gaps
- **93 isolated node(s):** `testPlayer`, `testInventoryItem`, `testGameState`, `taggedPlayer`, `gameWithInterface` (+88 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Community 31`** (2 nodes): `messages_test.go`, `TestErrorCodes()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 32`** (2 nodes): `Writing Documentation Rules`, `autoebiten Overview`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 33`** (2 nodes): `state Command`, `State Exporter`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 34`** (2 nodes): `Thread Safety`, `EbitenUI Integration`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 35`** (2 nodes): `launch Command`, `Crash Diagnostics`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 36`** (1 nodes): `custom.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 37`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 38`** (1 nodes): `errors.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 39`** (1 nodes): `caller_export_test.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 40`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 41`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 42`** (1 nodes): `cmd/autoebiten/main.go - CLI entry point with cobra`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 43`** (1 nodes): `cmd/docgen/main.go - Documentation generator`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 44`** (1 nodes): `internal/output/manager_test.go - Output manager tests`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 45`** (1 nodes): `internal/output/output_test.go - Output path tests`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 46`** (1 nodes): `test_testerror`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 47`** (1 nodes): `SocketPath`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 48`** (1 nodes): `example_custom_cmd`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 49`** (1 nodes): `example_crash_diag`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 50`** (1 nodes): `example_simple`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 51`** (1 nodes): `testkit_errinvalidstate`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 52`** (1 nodes): `doc.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 53`** (1 nodes): `TabInfo`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 54`** (1 nodes): `autoui_internal_doc`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 55`** (1 nodes): `claude_doc`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 56`** (1 nodes): `autoui_at`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 57`** (1 nodes): `input_command`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 58`** (1 nodes): `mouse_command`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 59`** (1 nodes): `wheel_command`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 60`** (1 nodes): `xpath_1_0`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 61`** (1 nodes): `ticks`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 62`** (1 nodes): `Installation Instructions`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 63`** (1 nodes): `Graphify Knowledge Graph`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 64`** (1 nodes): `ping Command`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `contains()` connect `Community 1` to `Community 0`, `Community 4`, `Community 5`, `Community 6`, `Community 7`, `Community 8`, `Community 9`, `Community 11`, `Community 12`, `Community 13`?**
  _High betweenness centrality (0.089) - this node is a cross-community bridge._
- **Why does `main()` connect `Community 9` to `Community 0`, `Community 1`, `Community 2`, `Community 3`, `Community 11`, `Community 15`, `Community 20`?**
  _High betweenness centrality (0.070) - this node is a cross-community bridge._
- **Why does `Executor` connect `Community 0` to `Community 3`, `Community 14`?**
  _High betweenness centrality (0.070) - this node is a cross-community bridge._
- **Are the 77 inferred relationships involving `contains()` (e.g. with `main()` and `TestOutputManagerDiffAndUpdateSnapshot()`) actually correct?**
  _`contains()` has 77 INFERRED edges - model-reasoned connections that need verification._
- **Are the 45 inferred relationships involving `WalkTree()` (e.g. with `TestHandleExistsCommand_Found()` and `TestHandleExistsCommand_NotFound()`) actually correct?**
  _`WalkTree()` has 45 INFERRED edges - model-reasoned connections that need verification._
- **Are the 38 inferred relationships involving `createTestNineSlice()` (e.g. with `TestHandleExistsCommand_Found()` and `TestHandleExistsCommand_NotFound()`) actually correct?**
  _`createTestNineSlice()` has 38 INFERRED edges - model-reasoned connections that need verification._
- **Are the 39 inferred relationships involving `ExtractWidgetInfo()` (e.g. with `ExtractWidgetState()` and `ExtractCustomData()`) actually correct?**
  _`ExtractWidgetInfo()` has 39 INFERRED edges - model-reasoned connections that need verification._
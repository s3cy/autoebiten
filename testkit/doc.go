// Package testkit provides a testing framework for autoebiten games.
//
// testkit offers two testing modes:
//
//  1. Black-Box Testing (Game): Launches the game in a separate process and
//     controls it via JSON-RPC over Unix sockets. This mode tests the game
//     as a black box, identical to how real users interact with it.
//
//  2. White-Box Testing (Mock): Tests game logic in the same process with
//     mocked inputs. This mode provides fine-grained control over game state
//     and is ideal for unit testing specific behaviors.
//
// Both modes use a similar API for input injection (key presses, mouse movements)
// and state observation, allowing tests to be written in a mode-agnostic way.
//
// Example black-box test:
//
//	func TestPlayerMovement(t *testing.T) {
//	    game := testkit.Launch(t, "./mygame", testkit.WithTimeout(30*time.Second))
//	    defer game.Shutdown()
//
//	    // Wait for game to be ready
//	    game.WaitFor(func() bool {
//	        err := game.Ping()
//	        return err == nil
//	    }, 5*time.Second)
//
//	    // Press movement key for 10 ticks
//	    game.HoldKey(ebiten.KeyArrowRight, 10)
//
//	    // Verify player moved
//	    x, err := game.StateQuery("gamestate", "Player.X")
//	    require.NoError(t, err)
//	    assert.Greater(t, x.(float64), 0.0)
//	}
//
// Example white-box test:
//
//	func TestPlayerMovement(t *testing.T) {
//	    myGame := NewMyGame()
//	    mock := testkit.NewMock(t, myGame)
//
//	    // Inject key press
//	    mock.InjectKeyPress(ebiten.KeyRight)
//
//	    // Advance 10 ticks
//	    mock.Ticks(10)
//
//	    // Verify player moved
//	    assert.Greater(t, myGame.Player.X, 0)
//	}
//
// State Export:
//
// Games can export internal state using autoebiten.RegisterStateExporter, which
// enables reflection-based state queries via dot-notation paths like
// "Player.X" or "Inventory.0.Name":
//
//	func init() {
//	    autoebiten.RegisterStateExporter("gamestate", &gameInstance)
//	}
package testkit

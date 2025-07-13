package plugin

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rsmacapinlac/pomodux/internal/logger"
)

// TestMain initializes the logger for all tests in this package
func TestMain(m *testing.M) {
	// Initialize logger for tests
	logConfig := &logger.Config{
		Level:      logger.LogLevelDebug,
		Format:     "text",
		Output:     "console",
		ShowCaller: false,
	}
	if err := logger.Init(logConfig); err != nil {
		panic("Failed to initialize logger for tests: " + err.Error())
	}

	// Run tests
	os.Exit(m.Run())
}

func TestNewPluginManager(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)

	if pm == nil {
		t.Fatal("NewPluginManager returned nil")
	}

	if pm.pluginsDir != pluginsDir {
		t.Errorf("Expected plugins directory %s, got %s", pluginsDir, pm.pluginsDir)
	}

	// Cleanup
	pm.Shutdown()
}

func TestLoadPlugin(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Test plugin code
	pluginCode := `
pomodux.register_plugin({
    name = "test_plugin",
    version = "1.0.0",
    description = "Test plugin",
    author = "Test Author"
})

pomodux.register_hook("timer_started", function(event)
    pomodux.log("Timer started hook called")
end)
`

	err := pm.LoadPlugin("test_plugin", pluginCode)
	if err != nil {
		t.Fatalf("Failed to load plugin: %v", err)
	}

	// Check if plugin was loaded
	plugin, exists := pm.GetPlugin("test_plugin")
	if !exists {
		t.Fatal("Plugin not found after loading")
	}

	if plugin.Name != "test_plugin" {
		t.Errorf("Expected plugin name 'test_plugin', got '%s'", plugin.Name)
	}

	if plugin.Version != "1.0.0" {
		t.Errorf("Expected plugin version '1.0.0', got '%s'", plugin.Version)
	}

	if plugin.Description != "Test plugin" {
		t.Errorf("Expected plugin description 'Test plugin', got '%s'", plugin.Description)
	}

	if plugin.Author != "Test Author" {
		t.Errorf("Expected plugin author 'Test Author', got '%s'", plugin.Author)
	}
}

func TestLoadPluginFromFile(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Create a test plugin file
	pluginFile := filepath.Join(pluginsDir, "test_plugin.lua")
	pluginCode := `
pomodux.register_plugin({
    name = "test_plugin",
    version = "1.0.0",
    description = "Test plugin from file",
    author = "Test Author"
})

pomodux.register_hook("timer_completed", function(event)
    pomodux.log("Timer completed hook called")
end)
`

	err := os.WriteFile(pluginFile, []byte(pluginCode), 0644)
	if err != nil {
		t.Fatalf("Failed to write plugin file: %v", err)
	}

	// Load the plugin from file
	err = pm.LoadPluginFromFile(pluginFile)
	if err != nil {
		t.Fatalf("Failed to load plugin from file: %v", err)
	}

	// Check if plugin was loaded
	plugin, exists := pm.GetPlugin("test_plugin")
	if !exists {
		t.Fatal("Plugin not found after loading from file")
	}

	if plugin.Name != "test_plugin" {
		t.Errorf("Expected plugin name 'test_plugin', got '%s'", plugin.Name)
	}
}

func TestLoadPlugins(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Create multiple test plugin files
	plugins := []struct {
		name string
		code string
	}{
		{
			name: "plugin1.lua",
			code: `
pomodux.register_plugin({
    name = "plugin1",
    version = "1.0.0",
    description = "First test plugin",
    author = "Test Author"
})
`,
		},
		{
			name: "plugin2.lua",
			code: `
pomodux.register_plugin({
    name = "plugin2",
    version = "1.0.0",
    description = "Second test plugin",
    author = "Test Author"
})
`,
		},
		{
			name: "not_a_plugin.txt", // Should be ignored
			code: "This is not a plugin",
		},
	}

	// Write plugin files
	for _, p := range plugins {
		pluginFile := filepath.Join(pluginsDir, p.name)
		err := os.WriteFile(pluginFile, []byte(p.code), 0644)
		if err != nil {
			t.Fatalf("Failed to write plugin file %s: %v", p.name, err)
		}
	}

	// Load all plugins
	err := pm.LoadPlugins()
	if err != nil {
		t.Fatalf("Failed to load plugins: %v", err)
	}

	// Check if plugins were loaded
	loadedPlugins := pm.ListPlugins()
	if len(loadedPlugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(loadedPlugins))
	}

	// Check specific plugins
	plugin1, exists := pm.GetPlugin("plugin1")
	if !exists {
		t.Fatal("Plugin1 not found")
	}

	plugin2, exists := pm.GetPlugin("plugin2")
	if !exists {
		t.Fatal("Plugin2 not found")
	}

	if plugin1.Description != "First test plugin" {
		t.Errorf("Expected plugin1 description 'First test plugin', got '%s'", plugin1.Description)
	}

	if plugin2.Description != "Second test plugin" {
		t.Errorf("Expected plugin2 description 'Second test plugin', got '%s'", plugin2.Description)
	}
}

func TestEmitEvent(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Create a test plugin that logs events
	pluginCode := `
local event_count = 0

pomodux.register_plugin({
    name = "event_test",
    version = "1.0.0",
    description = "Event test plugin",
    author = "Test Author"
})

pomodux.register_hook("timer_started", function(event)
    event_count = event_count + 1
    pomodux.log("Event received: " .. event.type)
end)
`

	err := pm.LoadPlugin("event_test", pluginCode)
	if err != nil {
		t.Fatalf("Failed to load plugin: %v", err)
	}

	// Emit an event
	event := Event{
		Type:      EventTimerStarted,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"session_type": "work",
			"duration":     1500,
		},
	}

	pm.EmitEvent(event)

	// Give some time for the event to be processed
	time.Sleep(100 * time.Millisecond)
}

func TestEnableDisablePlugin(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Load a test plugin
	pluginCode := `
pomodux.register_plugin({
    name = "test_plugin",
    version = "1.0.0",
    description = "Test plugin",
    author = "Test Author"
})
`

	err := pm.LoadPlugin("test_plugin", pluginCode)
	if err != nil {
		t.Fatalf("Failed to load plugin: %v", err)
	}

	// Check initial state
	plugin, exists := pm.GetPlugin("test_plugin")
	if !exists {
		t.Fatal("Plugin not found")
	}

	if !plugin.Enabled {
		t.Error("Plugin should be enabled by default")
	}

	// Disable plugin
	err = pm.EnablePlugin("test_plugin", false)
	if err != nil {
		t.Fatalf("Failed to disable plugin: %v", err)
	}

	plugin, _ = pm.GetPlugin("test_plugin")
	if plugin.Enabled {
		t.Error("Plugin should be disabled")
	}

	// Enable plugin
	err = pm.EnablePlugin("test_plugin", true)
	if err != nil {
		t.Fatalf("Failed to enable plugin: %v", err)
	}

	plugin, _ = pm.GetPlugin("test_plugin")
	if !plugin.Enabled {
		t.Error("Plugin should be enabled")
	}
}

func TestUnloadPlugin(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Load a test plugin
	pluginCode := `
pomodux.register_plugin({
    name = "test_plugin",
    version = "1.0.0",
    description = "Test plugin",
    author = "Test Author"
})
`

	err := pm.LoadPlugin("test_plugin", pluginCode)
	if err != nil {
		t.Fatalf("Failed to load plugin: %v", err)
	}

	// Check plugin exists
	_, exists := pm.GetPlugin("test_plugin")
	if !exists {
		t.Fatal("Plugin not found after loading")
	}

	// Unload plugin
	err = pm.UnloadPlugin("test_plugin")
	if err != nil {
		t.Fatalf("Failed to unload plugin: %v", err)
	}

	// Check plugin no longer exists
	_, exists = pm.GetPlugin("test_plugin")
	if exists {
		t.Error("Plugin still exists after unloading")
	}
}

func TestInvalidPlugin(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Test plugin without registration
	invalidCode := `
-- This plugin doesn't register itself
print("Hello world")
`

	err := pm.LoadPlugin("invalid_plugin", invalidCode)
	if err == nil {
		t.Error("Expected error when loading invalid plugin")
	}
}

func TestDuplicatePlugin(t *testing.T) {
	pluginsDir := t.TempDir()
	pm := NewPluginManager(pluginsDir)
	defer pm.Shutdown()

	// Load the same plugin twice
	pluginCode := `
pomodux.register_plugin({
    name = "duplicate_plugin",
    version = "1.0.0",
    description = "Test plugin",
    author = "Test Author"
})
`

	err := pm.LoadPlugin("duplicate_plugin", pluginCode)
	if err != nil {
		t.Fatalf("Failed to load plugin first time: %v", err)
	}

	err = pm.LoadPlugin("duplicate_plugin", pluginCode)
	if err == nil {
		t.Error("Expected error when loading duplicate plugin")
	}
}
